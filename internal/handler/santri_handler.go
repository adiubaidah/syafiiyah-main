package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/config"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SantriHandler interface {
	CreateSantriHandler(c *gin.Context)
	ListSantriHandler(c *gin.Context)
	GetSantriHandler(c *gin.Context)
	UpdateSantriHandler(c *gin.Context)
	DeleteSantriHandler(c *gin.Context)
}

type santriHandler struct {
	config  *config.Config
	logger  *logrus.Logger
	usecase usecase.SantriUseCase
}

func NewSantriHandler(config *config.Config, logger *logrus.Logger, usecase usecase.SantriUseCase) SantriHandler {
	return &santriHandler{
		config:  config,
		logger:  logger,
		usecase: usecase,
	}
}

func (h *santriHandler) CreateSantriHandler(c *gin.Context) {
	var santriRequest model.CreateSantriRequest
	if err := c.ShouldBind(&santriRequest); err != nil {
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
		photoPath := filepath.Join(config.PathPhoto, fileName)
		if err := c.SaveUploadedFile(photo, photoPath); err != nil {
			h.logger.Error(err)
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Failed to save photo"})
			return
		}

		santriRequest.Photo = fileName
	}
	result, err := h.usecase.CreateSantri(c, &santriRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(201, model.ResponseData[model.SantriResponse]{Code: 201, Status: "Created", Data: result})
}

func (h *santriHandler) ListSantriHandler(c *gin.Context) {
	limitNumber, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limitNumber = 10
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	generation, err := strconv.Atoi(c.Query("generation"))
	if err != nil {
		generation = 0
	}

	isActive, err := strconv.Atoi(c.Query("is_active"))
	if err != nil {
		isActive = -1
	}

	occupationId, err := strconv.Atoi(c.Query("occupation_id"))
	if err != nil {
		occupationId = 0
	}

	search := c.Query("q")

	listSantriRequest := model.ListSantriRequest{
		Q:            search,
		Order:        c.Query("order"),
		Limit:        int32(limitNumber),
		Page:         int32(page),
		Generation:   int32(generation),
		IsActive:     isActive,
		OccupationID: int32(occupationId),
	}

	if err := c.ShouldBind(&listSantriRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	result, err := h.usecase.ListSantri(c, &listSantriRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	//format all image from filename to url
	for i, santri := range result {
		if santri.Photo == "" {
			continue
		}
		result[i].Photo = fmt.Sprintf("%s/photo/%s", h.config.ServerPublicUrl, santri.Photo)
	}

	count, err := h.usecase.CountSantri(c, &listSantriRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	pagination := model.Pagination{
		CurrentPage:  page,
		TotalPages:   count/limitNumber + 1,
		TotalItems:   count,
		ItemsPerPage: limitNumber,
	}

	c.JSON(200, model.ResponseData[model.ListSantriResponse]{
		Code:   200,
		Status: "OK",
		Data: model.ListSantriResponse{
			Items:      result,
			Pagination: pagination,
		},
	})
}

func (h *santriHandler) GetSantriHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	santriId := int32(id)
	result, err := h.usecase.GetSantri(c, &santriId)
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	result.Photo = fmt.Sprintf("%s/photo/%s", h.config.ServerPublicUrl, result.Photo)
	c.JSON(200, model.ResponseData[model.SantriCompleteResponse]{Code: 200, Status: "OK", Data: result})
}

func (h *santriHandler) UpdateSantriHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}
	santriId := int32(id)
	oldSantri, err := h.usecase.GetSantri(c, &santriId)
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	var santriRequest model.UpdateSantriRequest
	if err := c.ShouldBind(&santriRequest); err != nil {
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
		photoPath := filepath.Join(config.PathPhoto, fileName)
		if err := c.SaveUploadedFile(photo, photoPath); err != nil {
			h.logger.Error(err)
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Failed to save photo"})
			return
		}
		santriRequest.Photo = fileName
		//delete old photo
		util.DeleteFile(filepath.Join(config.PathPhoto, oldSantri.Photo))
	}
	result, err := h.usecase.UpdateSantri(c, &santriRequest, &santriId)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(201, model.ResponseData[model.SantriResponse]{Code: 201, Status: "Created", Data: result})
}

func (h *santriHandler) DeleteSantriHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Error(err)
		ctx.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}
	santriId := int32(id)

	deletedSantri, err := h.usecase.DeleteSantri(ctx, &santriId)
	if err != nil {
		h.logger.Error(err)
		ctx.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	//delete photo also
	util.DeleteFile(filepath.Join(config.PathPhoto, deletedSantri.Photo))

	ctx.JSON(200, model.ResponseData[model.SantriResponse]{Code: 200, Status: "OK", Data: deletedSantri})
}
