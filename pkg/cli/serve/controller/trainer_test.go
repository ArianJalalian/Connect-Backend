package serve_test

import (
	serve2 "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestApp() (*fiber.App, *serve2.TrainerController) {
	app := fiber.New()
	trainerCtrl := &serve2.TraineeController{}
	trainerGroup := app.Group("/trainer") // Mimic the group used in the controller
	trainerCtrl.RegisterRoutes(trainerGroup)
	return app, nil
}

func TestEditProfile(t *testing.T) {
	app, _ := setupTestApp()

	editProfileData := dto.TrainerEdit{}
	body, _ := json.Marshal(editProfileData)

	req := httptest.NewRequest(http.MethodPut, "/trainer/profile/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "mockJWTToken") // Set a mock JWT token for authorization

	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.TrainerResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

}

func TestGetTrainerProfile(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/trainer/profile/", nil)
	req.Header.Set("Authorization", "mockJWTToken") // Set a mock JWT token for authorization

	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.TrainerResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
}

func TestTrainerController_GetAllRequests(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1haGRpZWhtb2doaXNlaDgxdHJhaW5lcjJAZ21haWwuY29tIiwiZXhwIjoxNzE2NzU1MzI3LCJ1c2VyX2lkIjoyfQ.EhOIyG5fUDM3aP7g11aUu98QDiEsrVs9yCWKi8EklD8"

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:1234/trainer/requests", nil)
	req.Header.Set("Authorization", token)
	resp := httptest.NewRecorder()

	assert.Equal(t, http.StatusOK, resp.Code)
}
