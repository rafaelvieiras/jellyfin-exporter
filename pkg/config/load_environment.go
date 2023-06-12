package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// LoadEnvironment loads environment variables from a file and sets them.
func LoadEnvironment(filename string) error {
	// Loading from the .env file
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				if os.Getenv(key) == "" { // Only set if not previously set
					os.Setenv(key, value)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}

	// Check if the necessary environment variables are set
	if os.Getenv("JELLYFIN_API_URL") == "" || os.Getenv("JELLYFIN_TOKEN") == "" {
		return errors.New("JELLYFIN_API_URL and JELLYFIN_TOKEN environment variables are required")
	}

	// Set default server port if not set
	if os.Getenv("SERVER_PORT") == "" {
		os.Setenv("SERVER_PORT", "8080")
	}

	// Set default metrics path if not set
	if os.Getenv("METRICS_PATH") == "" {
		os.Setenv("METRICS_PATH", "/metrics")
	}

	// Set default server address if not set
	if os.Getenv("SERVER_ADDR") == "" {
		os.Setenv("SERVER_ADDR", "")
	}

	return nil
}
