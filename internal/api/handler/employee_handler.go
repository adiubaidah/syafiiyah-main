package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	storage "github.com/adiubaidah/rfid-syafiiyah/platform/storage"
)

type EmployeeHandler struct {
	Logger  *logrus.Logger
	Storage *storage.StorageManager
	UseCase *usecase.EmployeeUseCase
}

func NewEmployeeHandler(args *EmployeeHandler) *EmployeeHandler {
	return args
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var request model.CreateEmployeeRequest
	if err := c.ShouldBind(&request); err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	photo, err := c.FormFile("photo")
	if err != nil {
		if err != http.ErrMissingFile {
			h.Logger.Error(err)
			c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
			return
		}
	} else {
		if err := util.ValidatePhoto(photo); err != nil {
			c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
			return
		}
		fileName := fmt.Sprintf("%s%s", uuid.New().String(), util.GetFileExtension(photo))
		if request.Photo, err = h.Storage.UploadFile(c, photo, fileName); err != nil {
			h.Logger.Error(err)
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Failed to save photo"})
			return
		}

	}
	result, err := h.UseCase.Create(c, &request)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(201, model.ResponseData[*model.Employee]{Code: 201, Status: "Created", Data: result})
}

func (h *EmployeeHandler) List(c *gin.Context) {

	var request model.ListEmployeeRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	if request.Limit == 0 {
		request.Limit = 10
	}
	if request.Page == 0 {
		request.Page = 1
	}
	result, err := h.UseCase.List(c, &request)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	count, err := h.UseCase.Count(c, &request)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	pagination := model.Pagination{
		CurrentPage:  request.Page,
		TotalPages:   (int32(count) + request.Limit - 1) / request.Limit,
		TotalItems:   count,
		ItemsPerPage: request.Limit,
	}

	c.JSON(200, model.ResponseData[model.ListEmployeeResponse]{
		Code:   200,
		Status: "OK",
		Data: model.ListEmployeeResponse{
			Items:      result,
			Pagination: pagination,
		},
	})
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	id := c.Param("id")
	employeeID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	var request model.UpdateEmployeeRequest
	if err := c.ShouldBind(&request); err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	photo, err := c.FormFile("photo")
	if err != nil {
		if err != http.ErrMissingFile {
			h.Logger.Error(err)
			c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
			return
		}
	} else {
		if err := util.ValidatePhoto(photo); err != nil {
			c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
			return
		}
		fileName := fmt.Sprintf("%s%s", uuid.New().String(), util.GetFileExtension(photo))
		if request.Photo, err = h.Storage.UploadFile(c, photo, fileName); err != nil {
			h.Logger.Error(err)
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Failed to save photo"})
			return
		}

	}
	result, err := h.UseCase.Update(c, &request, int32(employeeID))
	if err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(200, model.ResponseData[*model.Employee]{Code: 200, Status: "OK", Data: result})
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	employeeID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}
	result, err := h.UseCase.Delete(c, int32(employeeID))
	if err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	h.Storage.DeleteFile(c, result.Photo)

	c.JSON(200, model.ResponseData[*model.Employee]{Code: 200, Status: "OK", Data: result})
}
