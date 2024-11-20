package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase/mocks"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSantriOccupationHandler(t *testing.T) {
	// Mock usecase
	mockUsecase := new(mocks.SantriOccuapationUsecase)

	// Mock logger
	logger := logrus.New()

	// Create handler
	h := NewSantriOccupationHandler(logger, mockUsecase)

	// Initialize Gin router
	router := gin.Default()
	router.POST("/santri-occupation", h.CreateSantriOccupationHandler)

	t.Run("Success - Valid Input", func(t *testing.T) {
		// Mock response from usecase
		mockUsecase.On("CreateSantriOccupation", mock.Anything, mock.AnythingOfType("*model.CreateSantriOccupationRequest")).
			Return(model.SantriOccupationResponse{
				ID:          1,
				Name:        "Pendidik",
				Description: "Mengajar santri",
			}, nil)

		// Define input
		requestBody := model.CreateSantriOccupationRequest{
			Name:        "Pendidik",
			Description: "Mengajar santri",
		}

		// Convert to JSON
		bodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP request
		req, _ := http.NewRequest("POST", "/santri-occupation", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Mock HTTP response
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)

		var response model.ResponseData[model.SantriOccupationResponse]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Created", response.Status)
		assert.Equal(t, int32(1), response.Data.ID)
		assert.Equal(t, "Pendidik", response.Data.Name)
		assert.Equal(t, "Mengajar santri", response.Data.Description)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("Fail - Invalid Input", func(t *testing.T) {
		// Define invalid input (missing required name field)
		requestBody := map[string]string{
			"description": "Mengajar santri",
		}

		// Convert to JSON
		bodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP request
		req, _ := http.NewRequest("POST", "/santri-occupation", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Mock HTTP response
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "error", response["status"])
		assert.Contains(t, response["message"], "Key: 'CreateSantriOccupationRequest.Name'")
	})
}

func TestListSantriOccupationHandler(t *testing.T) {
	// Mock usecase
	mockUsecase := new(mocks.SantriOccuapationUsecase)

	// Mock logger
	logger := logrus.New()

	// Create handler
	h := NewSantriOccupationHandler(logger, mockUsecase)

	// Initialize Gin router
	router := gin.Default()
	router.GET("/santri-occupation", h.ListSantriOccupationHandler)

	t.Run("Success", func(t *testing.T) {
		// Mock response from usecase
		mockUsecase.On("ListSantriOccupations", mock.Anything).
			Return([]model.SantriOccupationWithCountResponse{
				{
					SantriOccupationResponse: model.SantriOccupationResponse{
						ID:          1,
						Name:        "Pendidik",
						Description: "Mengajar santri",
					},
					Count: 5,
				},
				{
					SantriOccupationResponse: model.SantriOccupationResponse{
						ID:          2,
						Name:        "Pengasuh",
						Description: "Mengasuh santri",
					},
					Count: 3,
				},
			}, nil)

		// Create HTTP request
		req, _ := http.NewRequest("GET", "/santri-occupation", nil)

		// Mock HTTP response
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response model.ResponseData[[]model.SantriOccupationWithCountResponse]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "OK", response.Status)
		assert.Len(t, response.Data, 2)
		assert.Equal(t, int32(1), response.Data[0].ID)
		assert.Equal(t, "Pendidik", response.Data[0].Name)
		assert.Equal(t, "Mengajar santri", response.Data[0].Description)
		assert.Equal(t, int32(5), response.Data[0].Count)
		assert.Equal(t, int32(2), response.Data[1].ID)
		assert.Equal(t, "Pengasuh", response.Data[1].Name)
		assert.Equal(t, "Mengasuh santri", response.Data[1].Description)
		assert.Equal(t, int32(3), response.Data[1].Count)
	})
}
