package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SantriScheduleHandler interface {
	CreateSantriScheduleHandler(c *gin.Context)
	ListSantriScheduleHandler(c *gin.Context)
	UpdateSantriScheduleHandler(c *gin.Context)
	DeleteSantriScheduleHandler(c *gin.Context)
}

type santriScheduleHandler struct {
	logger  *logrus.Logger
	usecase usecase.SantriScheduleUseCase
}

func NewSantriScheduleHandler(logger *logrus.Logger, usecase usecase.SantriScheduleUseCase) SantriScheduleHandler {
	return &santriScheduleHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (h *santriScheduleHandler) CreateSantriScheduleHandler(c *gin.Context) {
	var santriScheduleRequest model.CreateSantriScheduleRequest
	if err := c.ShouldBindJSON(&santriScheduleRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	result, err := h.usecase.CreateSantriSchedule(context.Background(), &santriScheduleRequest)
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
	c.JSON(http.StatusCreated, model.ResponseData[model.SantriScheduleResponse]{Code: http.StatusCreated, Status: "Created", Data: result})
}

func (h *santriScheduleHandler) ListSantriScheduleHandler(c *gin.Context) {
	result, err := h.usecase.ListSantriSchedule(context.Background())
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ResponseData[[]model.SantriScheduleResponse]{Code: http.StatusOK, Status: "OK", Data: result})
}

func (h *santriScheduleHandler) UpdateSantriScheduleHandler(c *gin.Context) {
	var santriScheduleRequest model.UpdateSantriScheduleRequest
	if err := c.ShouldBindJSON(&santriScheduleRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	idParam := c.Param("id")
	santriScheduleId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	result, err := h.usecase.UpdateSantriSchedule(context.Background(), &santriScheduleRequest, int32(santriScheduleId))
	if err != nil {
		h.logger.Error(err)

		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ResponseData[model.SantriScheduleResponse]{Code: http.StatusOK, Status: "OK", Data: result})
}

func (h *santriScheduleHandler) DeleteSantriScheduleHandler(c *gin.Context) {
	idParam := c.Param("id")
	santriScheduleId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	result, err := h.usecase.DeleteSantriSchedule(context.Background(), int32(santriScheduleId))
	if err != nil {
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		} else {
			c.JSON(http.StatusInternalServerError, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, model.ResponseData[model.SantriScheduleResponse]{Code: http.StatusOK, Status: "OK", Data: result})
}
