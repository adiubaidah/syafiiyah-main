package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
)

type SantriPresenceHandler interface {
	CreateSantriPresenceHandler(c *gin.Context)
}

type santriPresenceHandler struct {
	logger  *logrus.Logger
	usecase usecase.SantriPresenceUseCase
}

func NewSantriPresenceHandler(logger *logrus.Logger, usecase usecase.SantriPresenceUseCase) SantriPresenceHandler {
	return &santriPresenceHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (h *santriPresenceHandler) CreateSantriPresenceHandler(c *gin.Context) {
	var santriPresenceRequest model.CreateSantriPresenceRequest
	if err := c.ShouldBindJSON(&santriPresenceRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	santriPresenceRequest.CreatedBy = db.PresenceCreatedByTypeAdmin

	result, err := h.usecase.CreateSantriPresence(context.Background(), &santriPresenceRequest)
	if err != nil {
		h.logger.Error(err)

		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		} else {
			c.JSON(http.StatusInternalServerError, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
			return
		}

	}
	c.JSON(http.StatusCreated, model.ResponseData[model.SantriPresenceResponse]{Code: http.StatusCreated, Status: "Created", Data: *result})
}
