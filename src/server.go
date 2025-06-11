package main

import (
	"errors"
	"html/template"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/storage/sqlite3"

	"github.com/meldiron/wp-admin-api/src/config"
	"github.com/meldiron/wp-admin-api/src/resources"
)

func main() {
	// Sessions
	storage := sqlite3.New()
	sessionStore := session.NewStore(session.Config{
		Storage: storage,
	})

	// Default handler for all requests (HTMX)
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if c.Locals("error") == nil {
				c.Locals("error", err.Error())
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
			html := template.Must(template.ParseFiles("src/index.html"))

			data := make(map[any]interface{})

			// Add whitelisted locals data
			keys := []string{"error", "component", "debugStatuses"} // TODO: Make dynamic
			for _, key := range keys {
				value := c.Locals(key)
				if value != nil {
					data[key] = value
				}
			}

			// Add all session data
			sess, err := sessionStore.Get(c)
			if err != nil {
				return errors.New("Session store not working.")
			}
			defer sess.Release()

			for _, key := range sess.Keys() {
				value := sess.Get(key)
				if value != nil {
					data[key] = value.(string)
				}
			}

			if data["component"] != nil {
				html.ExecuteTemplate(c.Response().BodyWriter(), data["component"].(string), data)
			} else {
				html.Execute(c.Response().BodyWriter(), data)
			}

			return nil
		},
	})

	// Expose public folder
	app.Get("/*", static.New("./public"))

	// Security
	// TODO: Add csrf and cors

	// Allow to use panic safely
	app.Use(recover.New())

	// Expose authorized data
	app.Use(func(c fiber.Ctx) error {
		sess, err := sessionStore.Get(c)
		if err != nil {
			return errors.New("Session store not working.")
		}
		defer sess.Release()

		if sess.Get("username") != nil {
			c.Locals("authorized", true)

			statuses, err := resources.GetServersStatuses()
			if err != nil {
				return err
			}

			c.Locals("debugStatuses", statuses)
		}

		return c.Next()
	})

	// Homepage (endpoint)
	app.Get("/", func(c fiber.Ctx) error {
		return errors.New("") // Render HTMX
	})

	// V1 API controller
	v1 := app.Group("/v1")

	// Create session (endpoint)
	v1.Post("/sessions", func(c fiber.Ctx) error {
		c.Locals("component", "auth")

		type CreateSessionsBody struct {
			Password string `form:"password"`
			Username string `form:"username"`
		}

		body := new(CreateSessionsBody)
		if err := c.Bind().Body(body); err != nil {
			return err
		}

		if !config.ValidateCredentials(body.Username, body.Password) {
			return errors.New("Invalid credentials.")
		}

		sess, err := sessionStore.Get(c)
		if err != nil {
			return errors.New("Session store not working.")
		}
		defer sess.Release()

		sess.Set("username", body.Username)
		err = sess.Save()
		if err != nil {
			return errors.New("Session could not be saved.")
		}

		c.Response().Header.Add("HX-Redirect", "/")
		c.Response().SetStatusCode(fiber.StatusNoContent)
		return nil
	})

	// Delete session (endpoint)
	v1.Delete("/sessions", func(c fiber.Ctx) error {
		c.Locals("component", "card-user")

		sess, err := sessionStore.Get(c)
		if err != nil {
			return errors.New("Session store not working.")
		}
		defer sess.Release()

		if sess.Get("username") != nil {
			sess.Delete("username")

			err = sess.Save()
			if err != nil {
				return errors.New("Session could not be saved.")
			}
		}

		c.Response().Header.Add("HX-Redirect", "/")
		c.Response().SetStatusCode(fiber.StatusNoContent)
		return nil
	})

	// Create debug action (endpoint)
	v1.Post("/actions/debug", func(c fiber.Ctx) error {
		c.Locals("component", "card-debug")

		if c.Locals("authorized") == nil {
			return errors.New("Please sign in first.")
		}

		type CreateDebugAction struct {
			Server string `form:"server"`
		}

		body := new(CreateDebugAction)
		if err := c.Bind().Body(body); err != nil {
			return err
		}

		path := resources.GetServerPath(body.Server)

		if path == "" {
			return errors.New("Server not found.")
		}

		status, err := resources.IsDebugEnabled(path)
		if err != nil {
			return err
		}

		if status == true {
			resources.ToggleDebugMode(path, false)
		} else {
			resources.ToggleDebugMode(path, true)
		}

		statuses, err := resources.GetServersStatuses()
		if err != nil {
			return err
		}

		c.Locals("debugStatuses", statuses)

		return errors.New("") // Render HTMX
	})

	app.Listen(":3000")
}
