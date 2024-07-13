package serve_test

import (
	serve2 "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	"encoding/json"
	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestTraineeControllerEndpoints(t *testing.T) {
	app := fiber.New()

	controller := &serve2.TraineeController{}
	api := app.Group("/trainee")
	controller.RegisterRoutes(api)

	os.Setenv("JWT_SECRET_KEY", "your_secret_key")

	t.Run("EditProfile", func(t *testing.T) {
		reqBody := dto.TraineeEdit{
			UserName: "testuser",
		}
		token := createTestToken("1")
		e := httpexpect.WithConfig(httpexpect.Config{
			BaseURL:  "http://localhost",
			Reporter: httpexpect.NewAssertReporter(t),
		})
		e.PUT("/trainee/profile/1").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(reqBody).
			Expect().
			Status(http.StatusOK).JSON().Object().Value("UserName").String().NotEmpty()
	})

	t.Run("GetTraineeProfile", func(t *testing.T) {
		token := createTestToken("1")
		e := httpexpect.WithConfig(httpexpect.Config{
			BaseURL:  "http://localhost", // Set your base URL here
			Reporter: httpexpect.NewAssertReporter(t),
		})
		e.GET("/trainee/profile/1").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).JSON().Object().Value("UserName").String().NotEmpty()
	})
}

func createTestToken(userID string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	return tokenString
}

func TestCreateProgramRequest(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1haGRpZWhtb2doaXNlaDgxdHJhaW5lZTJAZ21haWwuY29tIiwiZXhwIjoxNzE2NzU1Mjk2LCJ1c2VyX2lkIjoxfQ.IizpOZ5At3WcwYdR6jK4tGE8eDeRSKvJlT8FqBkcPnE"

	requestBody := dto.ProgramRequest{
		Description: "Sample program description",
		ActiveDays:  []bool{true, false, true, false, false, true, true},
	}

	requestBodyJSON, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "http://127.0.0.1:1234/trainee/request", bytes.NewBuffer(requestBodyJSON))
	req.Header.Set("Authorization", token)
	resp := httptest.NewRecorder()

	assert.Equal(t, http.StatusOK, resp.Code)
}
