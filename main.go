package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port        string
	MockDataDir string
}

// MockResponse represents the structure of a mock API response
type MockResponse struct {
	Status      int               `json:"status"`
	Headers     map[string]string `json:"headers"`
	Body        interface{}       `json:"body"`
	Delay       int               `json:"delay"`       // Delay in milliseconds
	ValidTokens []string          `json:"validTokens"` // Allowed tokens for authentication
}

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Initialize configuration
	cfg := Config{
		Port:        getEnv("PORT", "8080"),
		MockDataDir: getEnv("MOCK_DATA_PATH", "./data"),
	}

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
	})

	// Handle all routes dynamically
	app.All("/*", func(c *fiber.Ctx) error {
		mockFilePath := getMockFilePath(cfg.MockDataDir, c.Params("*"), c.Method(), c.Queries())
		if mockFilePath == "" {
			return c.Status(404).JSON(fiber.Map{"error": "Mock response not found"})
		}

		// Read JSON mock response
		mock := MockResponse{}
		if err := readJSON(mockFilePath, &mock); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Invalid JSON format"})
		}

		// Authorization check (if required)
		if len(mock.ValidTokens) > 0 && !isValidAuth(c.Get("Authorization"), mock.ValidTokens) {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Apply delay if configured
		if mock.Delay > 0 {
			time.Sleep(time.Duration(mock.Delay) * time.Millisecond)
		}

		// Set custom headers
		for key, value := range mock.Headers {
			c.Set(key, value)
		}

		return c.Status(mock.Status).JSON(mock.Body)
	})

	log.Printf("Mock API server running on http://localhost:%s 🚀", cfg.Port)
	log.Fatalf("%v", app.Listen(":"+cfg.Port))
}

// getEnv returns an environment variable value or fallback default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getMockFilePath constructs possible file paths based on method and query parameters
func getMockFilePath(root, path, method string, queries map[string]string) string {
	method = strings.ToLower(method)
	queryString := getSortedQueryString(queries)

	// Normalize path
	path = strings.Trim(path, "/")
	if path == "" {
		path = "index"
	}

	// Possible file paths
	files := []string{
		filepath.Join(root, path+"_"+method+".json"),                       // e.g., data/pets_get.json
		filepath.Join(root, path+"_"+method+"_"+queryString+".json"),       // e.g., data/pets_get_id=10.json
		filepath.Join(root, path, "index_"+method+".json"),                 // e.g., data/pets/index_get.json
		filepath.Join(root, path, "index_"+method+"_"+queryString+".json"), // e.g., data/pets/index_get_id=10.json
	}

	// Check existence of each possible file
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			return file
		}
	}
	return ""
}

// getSortedQueryString formats query parameters into a consistent sorted string
func getSortedQueryString(queries map[string]string) string {
	if len(queries) == 0 {
		return ""
	}
	query := url.Values{}
	for key, value := range queries {
		query.Add(key, value)
	}
	return strings.ReplaceAll(query.Encode(), "&", "_")
}

// readJSON reads and unmarshals a JSON file into a given struct
func readJSON(filePath string, v interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// isValidAuth checks if the Authorization header matches any valid tokens
func isValidAuth(token string, validTokens []string) bool {
	for _, validToken := range validTokens {
		if token == validToken {
			return true
		}
	}
	return false
}
