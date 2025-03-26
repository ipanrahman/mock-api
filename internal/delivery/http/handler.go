package http

import (
	"mock-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type MockHandler struct {
	mockUseCase usecase.MockUseCase
}

func NewMockHandler(mockUseCase usecase.MockUseCase) *MockHandler {
	return &MockHandler{mockUseCase: mockUseCase}
}

func (h *MockHandler) Handle(ctx *fiber.Ctx) error {
	mockResponse, err := h.mockUseCase.Process(ctx)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"details": err.Error(),
			},
		})
	}

	for key, value := range mockResponse.Headers {
		ctx.Response().Header.Set(key, value)
	}

	return ctx.Status(mockResponse.Status).JSON(mockResponse.Body)
}
