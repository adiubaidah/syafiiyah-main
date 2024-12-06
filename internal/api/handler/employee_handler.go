package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/config"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
)

type EmployeeHandler interface {
	CreateEmployeeHandler(c *gin.Context)
	// ListEmployeeHandler(c *gin.Context)
	// UpdateEmployeeHandler(c *gin.Context)
	// DeleteEmployeeHandler(c *gin.Context)
}

type employeeHandler struct {
	logger  *logrus.Logger
	config  *config.Config
	usecase usecase.EmployeeUseCase
}

func NewEmployeeHandler(logger *logrus.Logger, config *config.Config, usecase usecase.EmployeeUseCase) EmployeeHandler {
	return &employeeHandler{
		logger:  logger,
		config:  config,
		usecase: usecase,
	}
}

func (h *employeeHandler) CreateEmployeeHandler(c *gin.Context) {
	var request model.CreateEmployeeRequest
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
		photoPath := filepath.Join(config.PathPhoto, fileName)
		if err := c.SaveUploadedFile(photo, photoPath); err != nil {
			h.logger.Error(err)
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Failed to save photo"})
			return
		}

		request.Photo = fileName
	}
	result, err := h.usecase.CreateEmployee(c, &request)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(201, model.ResponseData[*model.Employee]{Code: 201, Status: "Created", Data: result})
}

// func (h *employeeHandler) ListEmployeeHandler(c *gin.Context) {

// 	var request model.ListEmployeeRequest

// 	if err := c.ShouldBindQuery(&request); err != nil {
// 		h.logger.Error(err)
// 		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
// 		return
// 	}

// 	if request.Limit == 0 {
// 		request.Limit = 10
// 	}
// 	if request.Page == 0 {
// 		request.Page = 1
// 	}
// 	result, err := h.usecase.ListEmployees(c, &request)
// 	if err != nil {
// 		h.logger.Error(err)
// 		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
// 		return
// 	}

// 	//format all image from filename to url
// 	for i := range *result {
// 		if (*result)[i].Photo == "" {
// 			continue
// 		}
// 		(*result)[i].Photo = fmt.Sprintf("%s/photo/%s", h.config.ServerPublicUrl, (*result)[i].Photo)
// 	}

// 	count, err := h.usecase.CountEmployees(c, &request)
// 	if err != nil {
// 		h.logger.Error(err)
// 		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
// 		return
// 	}

// 	pagination := model.Pagination{
// 		CurrentPage:  request.Page,
// 		TotalPages:   (int32(count) + request.Limit - 1) / request.Limit,
// 		TotalItems:   count,
// 		ItemsPerPage: request.Limit,
// 	}

// 	c.JSON(200, model.ResponseData[model.ListEmployeeResponse]{
// 		Code:   200,
// 		Status: "OK",
// 		Data: model.ListEmployeeResponse{
// 			Items:      result,
// 			Pagination: pagination,
// 		},
// 	})
// }
