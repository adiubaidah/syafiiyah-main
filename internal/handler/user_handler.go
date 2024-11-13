package handler

import (
	"strconv"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler interface {
	CreateUserHandler(c *gin.Context)
	ListUserHandler(c *gin.Context)
	GetUserHandler(c *gin.Context)
	UpdateUserHandler(c *gin.Context)
	DeleteUserHandler(c *gin.Context)
}

type userHandler struct {
	logger  *logrus.Logger
	usecase usecase.UserUseCase
}

func NewUserHandler(logger *logrus.Logger, usecase usecase.UserUseCase) UserHandler {
	return &userHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (h *userHandler) CreateUserHandler(c *gin.Context) {
	var userRequest model.CreateUserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	userResponse, err := h.usecase.CreateUser(c, &userRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(201, userResponse)
}

func (h *userHandler) ListUserHandler(c *gin.Context) {
	var listUserRequest model.ListUserRequest
	if err := c.ShouldBindQuery(&listUserRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	if listUserRequest.Limit == 0 {
		listUserRequest.Limit = 10
	}
	if listUserRequest.Page == 0 {
		listUserRequest.Page = 1
	}

	result, err := h.usecase.ListUsers(c, &listUserRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	count, err := h.usecase.CountUsers(c, &listUserRequest)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	pagination := model.Pagination{
		CurrentPage:  listUserRequest.Page,
		TotalPages:   (count + listUserRequest.Limit - 1) / listUserRequest.Limit,
		TotalItems:   count,
		ItemsPerPage: listUserRequest.Limit,
	}

	c.JSON(200, model.ResponseData[model.ListUserResponse]{
		Code:   200,
		Status: "success",
		Data: model.ListUserResponse{
			Items:      result,
			Pagination: pagination,
		},
	})
}

func (h *userHandler) GetUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	result, err := h.usecase.GetUser(c, int32(userId), "")
	if err != nil {
		h.logger.Error(err)

		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}

		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(200, model.ResponseData[model.UserResponse]{Code: 200, Status: "success", Data: model.UserResponse{
		ID:       result.ID,
		Username: result.Username,
		Role:     result.Role,
	}})
}

func (h *userHandler) UpdateUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	var userRequest model.UpdateUserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	result, err := h.usecase.UpdateUser(c, &userRequest, int32(userId))
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(200, model.ResponseData[model.UserResponse]{Code: 200, Status: "success", Data: result})
}

func (h *userHandler) DeleteUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	deletedUser, err := h.usecase.DeleteUser(c, int32(userId))
	if err != nil {
		h.logger.Error(err)

		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(200, model.ResponseData[model.UserResponse]{Code: 200, Status: "success", Data: deletedUser})
}
