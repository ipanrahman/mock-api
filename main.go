package main

import (
	"log"
	"mock-api/internal/config"
	httpDelivery "mock-api/internal/delivery/http"
	"mock-api/internal/repository"
	"mock-api/internal/usecase"
	"mock-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := &config.Config{
		Port:    utils.Env("PORT", "8080"),
		MockDir: utils.Env("MOCK_DIR", "./data"),
	}

	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
	})

	fileRepository := repository.NewFileRepository(cfg)
	mockUseCase := usecase.NewMockUseCase(cfg, fileRepository)
	mockHandler := httpDelivery.NewMockHandler(mockUseCase)

	app.All("*", mockHandler.Handle)

	log.Printf("Mock API server running on http://localhost:%s 🚀", cfg.Port)
	log.Fatalf("%v", app.Listen(":"+cfg.Port))
}
