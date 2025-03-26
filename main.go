package main

import (
	"fmt"
	"log/slog"
	"mock-api/internal/config"
	httpDelivery "mock-api/internal/delivery/http"
	"mock-api/internal/repository"
	"mock-api/internal/usecase"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file", err)
	}

	cfg := config.NewConfig()

	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
	})

	fileRepository := repository.NewFileRepository(cfg)
	mockUseCase := usecase.NewMockUseCase(cfg, fileRepository)
	mockHandler := httpDelivery.NewMockHandler(mockUseCase)

	app.All("*", mockHandler.Handle)

	fmt.Println("Listening on port " + cfg.Port)
	logger.Info(fmt.Sprintf("Mock API server running on http://localhost:%s 🚀", cfg.Port))
	logger.Info(fmt.Sprintf("%v", app.Listen(":"+cfg.Port)))
}
