package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"mock-api/internal/config"
	"mock-api/internal/domain"
	"mock-api/internal/repository"
	"mock-api/internal/utils"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type MockUseCase interface {
	Process(ctx *fiber.Ctx) (*domain.Response, error)
}
type mockUseCase struct {
	config         *config.Config
	fileRepository repository.FileRepository
}

func NewMockUseCase(config *config.Config, fileRepository repository.FileRepository) MockUseCase {
	return &mockUseCase{
		config:         config,
		fileRepository: fileRepository,
	}
}

func (uc *mockUseCase) Process(ctx *fiber.Ctx) (*domain.Response, error) {
	mockFilePath := uc.fileRepository.FindFilePath(
		uc.config.MockDir,
		ctx.Params("*"),
		ctx.Method(),
		ctx.Queries(),
	)
	if mockFilePath == "" {
		return nil, domain.ErrMockFileDoesNotExist
	}

	var mockResponse domain.MockResponse
	if err := uc.fileRepository.Find(mockFilePath, &mockResponse); err != nil {
		return nil, err
	}

	if len(mockResponse.Conditions) > 0 {
		if newMock, matched := uc.evaluateConditions(mockResponse.Conditions, ctx.Body(), &mockResponse); matched {
			mockResponse.Response = *newMock
		}
	}
	return &mockResponse.Response, nil
}

func (uc *mockUseCase) isValidAuth(token string, validTokens []string) bool {
	for _, validToken := range validTokens {
		if token == validToken {
			return true
		}
	}
	return false
}

func (uc *mockUseCase) evaluateConditions(conditions []domain.Condition, requestBody []byte, defaultResponse *domain.MockResponse) (*domain.Response, bool) {
	var reqData map[string]interface{}
	if err := json.Unmarshal(requestBody, &reqData); err != nil {
		return nil, false
	}

	for _, cond := range conditions {
		value, exists := reqData[cond.Field]
		if !exists {
			continue
		}

		switch cond.Operator {
		case "equals":
			fmt.Println("Equals: ", cond.Field, cond.Value, cond.Operator)
			if value == cond.Value {
				log.Println("Condition matched: equals")
				return &cond.Response, true
			}
		case "contains":
			if str, ok := value.(string); ok && strings.Contains(str, cond.Value.(string)) {
				log.Println("Condition matched: contains")
				return &cond.Response, true
			}
		case "pattern":
			if str, ok := value.(string); ok {
				if matched, _ := regexp.MatchString(cond.Value.(string), str); matched {
					log.Println("Condition matched: pattern")
					return &cond.Response, true
				}
			}
		case "greater_than":
			valueNum, valueOk := utils.ConvertToFloat64(value)
			condNum, condOk := utils.ConvertToFloat64(cond.Value)

			log.Printf("Checking condition: Field=%s, Operator=greater_than, Value Number=%v, Cond Num=%v", cond.Field, valueNum, condNum)

			if valueOk && condOk {
				if !(valueNum > condNum) {
					return &cond.Response, true
				}
			}
		}
	}

	return &defaultResponse.Response, true
}
