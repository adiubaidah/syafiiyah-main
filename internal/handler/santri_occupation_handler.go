package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
)

type SantriOccupationHandler interface {
	CreateSantriOccupationHandler(c *gin.Context)
	// ListSantriOccupationHandler(c *gin.Context)
	// UpdateSantriOccupationHandler(c *gin.Context)
	// DeleteSantriOccupationHandler(c *gin.Context)
}

type santriOccupationHandler struct {
	logger  *logrus.Logger
	usecase usecase.SantriOccuapationUsecase
}

func NewSantriOccupationHandler(logger *logrus.Logger, usecase usecase.SantriOccuapationUsecase) SantriOccupationHandler {
	return &santriOccupationHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (h *santriOccupationHandler) CreateSantriOccupationHandler(c *gin.Context) {
	var santriOccupationRequest model.CreateSantriOccupationRequest
	h.logger.Info("Test")
	if err := c.ShouldBindJSON(&santriOccupationRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := h.usecase.CreateSantriOccupation(context.Background(), santriOccupationRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, model.Response[model.SantriOccupationResponse]{Code: http.StatusCreated, Status: "Created", Data: result})
}
