package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/config"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	"github.com/adiubaidah/rfid-syafiiyah/platform/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ParentHandler interface {
	CreateParentHandler(c *gin.Context)
	ListParentHandler(c *gin.Context)
	GetParentHandler(c *gin.Context)
	UpdateParentHandler(c *gin.Context)
	DeleteParentHandler(c *gin.Context)
}

type parentHandler struct {
	config      *config.Config
	logger      *logrus.Logger
	storage     *storage.StorageManager
	usecase     usecase.ParentUseCase
	userUseCase usecase.UserUseCase
}

func NewParentHandler(logger *logrus.Logger, config *config.Config, storage *storage.StorageManager, usecase usecase.ParentUseCase, userUseCase usecase.UserUseCase) ParentHandler {
	return &parentHandler{
		config:      config,
		logger:      logger,
		usecase:     usecase,
		storage:     storage,
		userUseCase: userUseCase,
	}
}

func (h *parentHandler) CreateParentHandler(c *gin.Context) {
	var request model.CreateParentRequest
	if err := c.ShouldBind(&request); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	photo, err := c.FormFile("photo")
	if err != nil {
		if err != http.ErrMissingFile {
			h.logger.Error(err)
			c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
			return
		}
	} else {
		if err := util.ValidatePhoto(photo); err != nil {
			c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
			return
		}
		fileName := fmt.Sprintf("%s%s", uuid.New().String(), util.GetFileExtension(photo))
		if request.Photo, err = h.storage.UploadFile(c, photo, fileName); err != nil {
			h.logger.Error(err)
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Failed to save photo"})
			return
		}
	}

	if request.UserID != 0 {
		_, err := h.userUseCase.GetUser(c, request.UserID, "")
		if err != nil {
			h.logger.Error(err)
			c.JSON(404, model.ResponseMessage{Code: 404, Status: "error", Message: err.Error()})
			return
		}
	}

	result, err := h.usecase.CreateParent(c, &request)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(201, model.ResponseData[model.ParentResponse]{
		Code:   201,
		Status: "success",
		Data:   *result,
	})
}

func (h *parentHandler) ListParentHandler(c *gin.Context) {
	var request model.ListParentRequest
	if err := c.ShouldBind(&request); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}
	if request.Limit == 0 {
		request.Limit = 10
	}
	if request.Page == 0 {
		request.Page = 1
	}
	result, err := h.usecase.ListParents(c, &request)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	count, err := h.usecase.CountParents(c, &request)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	pagination := model.Pagination{
		CurrentPage:  request.Page,
		TotalPages:   int32((count + int64(request.Limit) - 1) / int64(request.Limit)),
		TotalItems:   count,
		ItemsPerPage: request.Limit,
	}

	c.JSON(200, model.ResponseData[model.ListParentResponse]{
		Code:   200,
		Status: "success",
		Data: model.ListParentResponse{
			Items:      *result,
			Pagination: pagination,
		},
	})
}

func (h *parentHandler) UpdateParentHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	oldData, err := h.usecase.GetParent(c, int32(id))
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, model.ResponseMessage{Code: 404, Status: "error", Message: err.Error()})
		return
	}

	var request model.UpdateParentRequest

	if err := c.ShouldBind(&request); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	if request.UserID != 0 {
		user, err := h.userUseCase.GetUser(c, request.UserID, "")
		if err != nil {
			h.logger.Error(err)
			c.JSON(404, model.ResponseMessage{Code: 404, Status: "error", Message: err.Error()})
			return
		}

		if user.Role != "parent" {
			c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "User is not a parent"})
			return
		}

	}

	photo, err := c.FormFile("photo")
	if err != nil {
		if err != http.ErrMissingFile {
			h.logger.Error(err)
			c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
			return
		}
	} else {
		if err := util.ValidatePhoto(photo); err != nil {
			c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
			return
		}
		fileName := fmt.Sprintf("%s%s", uuid.New().String(), util.GetFileExtension(photo))
		if request.Photo, err = h.storage.UploadFile(c, photo, fileName); err != nil {
			h.logger.Error(err)
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Failed to save photo"})
			return
		}

		if oldData.Photo != "" {
			h.storage.DeleteFile(context.Background(), oldData.Photo)
		}
	}

	result, err := h.usecase.UpdateParent(c, &request, int32(id))
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}

		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(200, model.ResponseData[model.ParentResponse]{Code: 200, Status: "success", Data: *result})

}

func (h *parentHandler) GetParentHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	result, err := h.usecase.GetParent(c, int32(id))
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}

		c.JSON(500, model.ResponseMessage{Code: 404, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(200, model.ResponseData[model.ParentResponse]{Code: 200, Status: "success", Data: *result})
}

func (h *parentHandler) DeleteParentHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	result, err := h.usecase.DeleteParent(c, int32(id))
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}

		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	if result.Photo != "" {
		h.storage.DeleteFile(context.Background(), result.Photo)
	}

	c.JSON(200, model.ResponseData[model.ParentResponse]{Code: 200, Status: "success", Data: *result})
}
