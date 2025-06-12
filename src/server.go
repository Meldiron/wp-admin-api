package main

import (
	"errors"
	"html/template"
	"time"

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
	storage := sqlite3.New(sqlite3.Config{
		Database: "./sqlite/fiber.sqlite3",
	})
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
			keys := []string{"error", "component", "debugStatuses", "success"} // TODO: Make dynamic
			for _, key := range keys {
				value := c.Locals(key)
				if value != nil {
					data[key] = value
				}
			}

			// Add all session data
			sess, err := sessionStore.Get(c)
			if err != nil {
				return errors.New("session store not working")
			}
			defer sess.Release()

			for _, key := range sess.Keys() {
				value := sess.Get(key)
				if value != nil {
					data[key] = value.(string)
				}
			}

			if data["component"] != nil {
				err := html.ExecuteTemplate(c.Response().BodyWriter(), data["component"].(string), data)
				if err != nil {
					return err
				}
			} else {
				err := html.Execute(c.Response().BodyWriter(), data)
				if err != nil {
					return err
				}
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
			return errors.New("session store not working")
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

	// Healthcheck endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"time":    time.Now().UTC().Format(time.RFC3339),
			"service": "wp-admin-api",
		})
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
			return errors.New("invalid credentials")
		}

		sess, err := sessionStore.Get(c)
		if err != nil {
			return errors.New("session store not working")
		}
		defer sess.Release()

		sess.Set("username", body.Username)
		err = sess.Save()
		if err != nil {
			return errors.New("session could not be saved")
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
			return errors.New("session store not working")
		}
		defer sess.Release()

		if sess.Get("username") != nil {
			sess.Delete("username")

			err = sess.Save()
			if err != nil {
				return errors.New("session could not be saved")
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
			return errors.New("please sign in first")
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
			return errors.New("server not found")
		}

		status, err := resources.IsDebugEnabled(path)
		if err != nil {
			return err
		}

		if status {
			err := resources.ToggleDebugMode(path, false)
			if err != nil {
				return err
			}
		} else {
			err := resources.ToggleDebugMode(path, true)
			if err != nil {
				return err
			}
		}

		statuses, err := resources.GetServersStatuses()
		if err != nil {
			return err
		}

		c.Locals("debugStatuses", statuses)

		return errors.New("") // Render HTMX
	})

	// Create restart action (endpoint)
	v1.Post("/actions/restart", func(c fiber.Ctx) error {
		c.Locals("component", "card-restart")

		if c.Locals("authorized") == nil {
			return errors.New("please sign in first")
		}

		type CreateRestartAction struct {
			Server string `form:"server"`
		}

		body := new(CreateRestartAction)
		if err := c.Bind().Body(body); err != nil {
			return err
		}

		path := resources.GetServerPath(body.Server)

		if path == "" {
			return errors.New("server not found")
		}

		err := resources.RestartServer(path)
		if err != nil {
			return err
		}

		c.Locals("success", "server restarted successfully")

		return errors.New("") // Render HTMX
	})

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
