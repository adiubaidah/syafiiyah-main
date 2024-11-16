// arduino_handler_test.go

package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/handler"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase/mocks"

	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
)

func TestCreateArduinoHandler(t *testing.T) {
	mockUsecase := new(mocks.ArduinoUseCase)
	mockUsecase.On("CreateArduino", mock.Anything, mock.AnythingOfType("*model.CreateArduinoRequest")).
		Return(&model.ArduinoResponse{
			ID:   3,
			Name: "Arduino 3",
		}, nil)

	h := handler.NewArduinoHandler(mockUsecase)

	router := gin.Default()
	router.POST("/arduino", h.CreateArduinoHandler)

	requestBody := model.CreateArduinoRequest{
		Name:  "Arduino 3",
		Modes: []db.ArduinoModeType{"excuse"},
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/arduino", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response model.ResponseData[model.ArduinoResponse]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	// additional assertions
}
