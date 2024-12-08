package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
)

type SantriPresenceHandler interface {
	CreateSantriPresenceHandler(c *gin.Context)
	ListSantriPresencesHandler(c *gin.Context)
	UpdateSantriPresenceHandler(c *gin.Context)
	DeleteSantriPresenceHandler(c *gin.Context)
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

func (h *santriPresenceHandler) ListSantriPresencesHandler(c *gin.Context) {
	var listSantriPresenceRequest model.ListSantriPresenceRequest
	if err := c.ShouldBindQuery(&listSantriPresenceRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	if listSantriPresenceRequest.Limit == 0 {
		listSantriPresenceRequest.Limit = 10
	}

	if listSantriPresenceRequest.Page == 0 {
		listSantriPresenceRequest.Page = 1
	}

	result, err := h.usecase.ListSantriPresences(context.Background(), &listSantriPresenceRequest)
	h.logger.Println("result", result)
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
	count, err := h.usecase.CountSantriPresences(context.Background(), &listSantriPresenceRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusInternalServerError, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	pagination := model.Pagination{
		CurrentPage:  listSantriPresenceRequest.Page,
		TotalPages:   int32((count + int64(listSantriPresenceRequest.Limit) - 1) / int64(listSantriPresenceRequest.Limit)),
		TotalItems:   count,
		ItemsPerPage: listSantriPresenceRequest.Limit,
	}

	c.JSON(http.StatusOK, model.ResponseData[model.ListSantriPresenceResponse]{Code: http.StatusOK, Status: "success", Data: model.ListSantriPresenceResponse{Data: *result, Pagination: pagination}})

}

func (h *santriPresenceHandler) UpdateSantriPresenceHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	var updateRequest model.UpdateSantriPresenceRequest

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
	}

	result, err := h.usecase.UpdateSantriPresence(context.Background(), &updateRequest, int32(id))
	if err != nil {
		h.logger.Error(err)

		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return

	}

	c.JSON(http.StatusOK, model.ResponseData[model.SantriPresenceResponse]{Code: http.StatusOK, Status: "success", Data: *result})
}

func (h *santriPresenceHandler) DeleteSantriPresenceHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	deleted, err := h.usecase.DeleteSantriPresence(context.Background(), int32(id))

	if err != nil {
		h.logger.Error(err)

		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ResponseData[model.SantriPresenceResponse]{Code: http.StatusOK, Status: "success", Data: *deleted})
}
