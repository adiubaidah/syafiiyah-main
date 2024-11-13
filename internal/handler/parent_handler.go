package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/config"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
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
	usecase     usecase.ParentUseCase
	userUseCase usecase.UserUseCase
}

func NewParentHandler(config *config.Config, logger *logrus.Logger, usecase usecase.ParentUseCase, userUseCase usecase.UserUseCase) ParentHandler {
	return &parentHandler{
		config:      config,
		logger:      logger,
		usecase:     usecase,
		userUseCase: userUseCase,
	}
}

func (h *parentHandler) CreateParentHandler(c *gin.Context) {
	var parentRequest model.CreateParentRequest
	if err := c.ShouldBind(&parentRequest); err != nil {
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
		if err := c.SaveUploadedFile(photo, fmt.Sprintf("%s/%s", config.PathPhoto, fileName)); err != nil {
			h.logger.Error(err)
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
			return
		}
		parentRequest.Photo = fileName
	}

	if parentRequest.UserID != 0 {
		_, err := h.userUseCase.GetUser(c, parentRequest.UserID, "")
		if err != nil {
			h.logger.Error(err)
			c.JSON(404, model.ResponseMessage{Code: 404, Status: "error", Message: err.Error()})
			return
		}
	}

	result, err := h.usecase.CreateParent(c, &parentRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(201, model.ResponseData[model.ParentResponse]{
		Code:   201,
		Status: "success",
		Data:   result,
	})
}

func (h *parentHandler) ListParentHandler(c *gin.Context) {
	var listParentRequest model.ListParentRequest
	if err := c.ShouldBind(&listParentRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}
	if listParentRequest.Limit == 0 {
		listParentRequest.Limit = 10
	}
	if listParentRequest.Page == 0 {
		listParentRequest.Page = 1
	}
	result, err := h.usecase.ListParents(c, &listParentRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	//format all image from filename to url
	for i, parent := range result {
		if parent.Photo == "" {
			continue
		}
		result[i].Photo = fmt.Sprintf("%s/photo/%s", h.config.ServerPublicUrl, parent.Photo)
	}

	count, err := h.usecase.CountParents(c, &listParentRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	pagination := model.Pagination{
		CurrentPage:  listParentRequest.Page,
		TotalPages:   (count + listParentRequest.Limit - 1) / listParentRequest.Limit,
		TotalItems:   count,
		ItemsPerPage: listParentRequest.Limit,
	}

	c.JSON(200, model.ResponseData[model.ListParentResponse]{
		Code:   200,
		Status: "success",
		Data: model.ListParentResponse{
			Items:      result,
			Pagination: pagination,
		},
	})
}

func (h *parentHandler) UpdateParentHandler(c *gin.Context) {
	idParam := c.Param("id")
	parentId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	oldParent, err := h.usecase.GetParent(c, int32(parentId))
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, model.ResponseMessage{Code: 404, Status: "error", Message: err.Error()})
		return
	}

	var parentRequest model.UpdateParentRequest

	if err := c.ShouldBind(&parentRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	if parentRequest.UserID != 0 {
		user, err := h.userUseCase.GetUser(c, parentRequest.UserID, "")
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
		if err := c.SaveUploadedFile(photo, fmt.Sprintf("%s/%s", config.PathPhoto, fileName)); err != nil {
			h.logger.Error(err)
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
			return
		}
		parentRequest.Photo = fileName

		if oldParent.Photo != "" {
			util.DeleteFile(filepath.Join(config.PathPhoto, oldParent.Photo))
		}
	}

	result, err := h.usecase.UpdateParent(c, &parentRequest, int32(parentId))
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}

		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(200, model.ResponseData[model.ParentResponse]{Code: 200, Status: "success", Data: result})

}

func (h *parentHandler) GetParentHandler(c *gin.Context) {
	idParam := c.Param("id")
	parentId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	result, err := h.usecase.GetParent(c, int32(parentId))
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}

		c.JSON(500, model.ResponseMessage{Code: 404, Status: "error", Message: err.Error()})
		return
	}

	if result.Photo != "" {
		result.Photo = fmt.Sprintf("%s/photo/%s", h.config.ServerPublicUrl, result.Photo)
	}

	c.JSON(200, model.ResponseData[model.ParentResponse]{Code: 200, Status: "success", Data: result})
}

func (h *parentHandler) DeleteParentHandler(c *gin.Context) {
	idParam := c.Param("id")
	parentId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	deletedParent, err := h.usecase.DeleteParent(c, int32(parentId))
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}

		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	if deletedParent.Photo != "" {
		util.DeleteFile(filepath.Join(config.PathPhoto, deletedParent.Photo))
	}

	c.JSON(200, model.ResponseData[model.ParentResponse]{Code: 200, Status: "success", Data: deletedParent})
}
