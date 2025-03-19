package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/exception"
	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	"github.com/adiubaidah/syafiiyah-main/internal/usecase"
)

type SantriOccupationHandler interface {
	CreateSantriOccupationHandler(c *gin.Context)
	ListSantriOccupationHandler(c *gin.Context)
	UpdateSantriOccupationHandler(c *gin.Context)
	DeleteSantriOccupationHandler(c *gin.Context)
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
		c.JSON(400, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	result, err := h.usecase.CreateSantriOccupation(context.Background(), &santriOccupationRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, model.ResponseData[model.SantriOccupationResponse]{Code: http.StatusCreated, Status: "Created", Data: *result})
}

func (h *santriOccupationHandler) ListSantriOccupationHandler(c *gin.Context) {
	result, err := h.usecase.ListSantriOccupations(context.Background())
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ResponseData[[]model.SantriOccupationWithCountResponse]{Code: http.StatusOK, Status: "OK", Data: *result})
}

func (h *santriOccupationHandler) UpdateSantriOccupationHandler(c *gin.Context) {

	idParam := c.Param("id")
	santriOccupationId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	var santriOccupationRequest model.UpdateSantriOccupationRequest
	if err := c.ShouldBindJSON(&santriOccupationRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	id := int32(santriOccupationId)

	result, err := h.usecase.UpdateSantriOccupation(context.Background(), &santriOccupationRequest, id)
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}

		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ResponseData[model.SantriOccupationResponse]{Code: http.StatusOK, Status: "OK", Data: *result})
}

func (h *santriOccupationHandler) DeleteSantriOccupationHandler(c *gin.Context) {
	idParam := c.Param("id")
	santriOccupationId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
		c.JSON(500, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	id := int32(santriOccupationId)
	result, err := h.usecase.DeleteSantriOccupation(context.Background(), id)
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}

		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ResponseData[model.SantriOccupationResponse]{Code: http.StatusOK, Status: "OK", Data: *result})
}
