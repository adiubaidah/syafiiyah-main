package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
)

type EmployeeOccupationHandler interface {
	CreateEmployeeOccupationHandler(c *gin.Context)
	ListEmployeeOccupationHandler(c *gin.Context)
	UpdateEmployeeOccupationHandler(c *gin.Context)
	DeleteEmployeeOccupationHandler(c *gin.Context)
}

type employeeOccupationHandler struct {
	logger  *logrus.Logger
	usecase usecase.EmployeeOccuapationUsecase
}

func NewEmployeeOccupationHandler(logger *logrus.Logger, usecase usecase.EmployeeOccuapationUsecase) EmployeeOccupationHandler {
	return &employeeOccupationHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (h *employeeOccupationHandler) CreateEmployeeOccupationHandler(c *gin.Context) {
	var employeeOccupationRequest model.CreateEmployeeOccupationRequest
	if err := c.ShouldBindJSON(&employeeOccupationRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	result, err := h.usecase.CreateEmployeeOccupation(context.Background(), &employeeOccupationRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, model.ResponseData[*model.EmployeeOccupationResponse]{Code: http.StatusCreated, Status: "Created", Data: result})
}

func (h *employeeOccupationHandler) ListEmployeeOccupationHandler(c *gin.Context) {
	result, err := h.usecase.ListEmployeeOccupations(context.Background())
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ResponseData[*[]model.EmployeeOccupationWithCountResponse]{Code: http.StatusOK, Status: "OK", Data: result})
}

func (h *employeeOccupationHandler) UpdateEmployeeOccupationHandler(c *gin.Context) {

	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	var employeeOccupationRequest model.UpdateEmployeeOccupationRequest
	if err := c.ShouldBindJSON(&employeeOccupationRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	id := int32(idParam)

	result, err := h.usecase.UpdateEmployeeOccupation(context.Background(), &employeeOccupationRequest, id)
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}

		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ResponseData[*model.EmployeeOccupationResponse]{Code: http.StatusOK, Status: "OK", Data: result})
}

func (h *employeeOccupationHandler) DeleteEmployeeOccupationHandler(c *gin.Context) {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
		c.JSON(500, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	id := int32(idParam)
	result, err := h.usecase.DeleteEmployeeOccupation(context.Background(), id)
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}

		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ResponseData[*model.EmployeeOccupationResponse]{Code: http.StatusOK, Status: "OK", Data: result})
}
