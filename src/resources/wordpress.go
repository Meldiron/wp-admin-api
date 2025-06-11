package resources

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strings"
)

var servers []map[string]string

func init() {
	usersEnv := os.Getenv("SERVERS")
	for _, user := range strings.Split(usersEnv, ",") {
		details := strings.Split(user, ":")
		if len(details) >= 2 {
			server := map[string]string{
				"name": details[0],
				"path": details[1],
			}
			servers = append(servers, server)
		}
	}
}

func GetServerPath(name string) string {
	for _, server := range servers {
		if server["name"] == name {
			return server["path"]
		}
	}
	return ""
}

func GetServersStatuses() ([]map[string]string, error) {
	statuses := []map[string]string{}

	for _, server := range servers {
		status, err := IsDebugEnabled(server["path"])

		if err != nil {
			return nil, err
		}

		statusVerbose := ""
		if status {
			statusVerbose = "enabled"
		} else {
			statusVerbose = "disabled"
		}

		statusMap := map[string]string{
			"name":   server["name"],
			"status": statusVerbose,
		}
		statuses = append(statuses, statusMap)
	}

	return statuses, nil
}

func IsDebugEnabled(path string) (bool, error) {
	contentBytes, err := os.ReadFile(path + "/docker-compose.yml")
	if err != nil {
		return false, err
	}

	content := string(contentBytes)
	if strings.Contains(content, "WORDPRESS_DEBUG=true") {
		return true, nil
	}

	return false, nil
}

func ToggleDebugMode(path string, newStatusBool bool) error {
	var newStatus string
	var oldStatus string

	if newStatusBool {
		newStatus = "true"
		oldStatus = "false"
	} else {
		newStatus = "false"
		oldStatus = "true"
	}

	contentBytes, err := os.ReadFile(path + "/docker-compose.yml")
	if err != nil {
		return err
	}

	content := string(contentBytes)

	content = strings.Replace(content, "WORDPRESS_DEBUG="+oldStatus, "WORDPRESS_DEBUG="+newStatus, 1)
	content = strings.Replace(content, "WORDPRESS_DEBUG_LOG="+oldStatus, "WORDPRESS_DEBUG_LOG="+newStatus, 1)

	err = os.WriteFile(path+"/docker-compose.yml", []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
