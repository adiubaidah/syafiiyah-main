package handler

import (
	"strconv"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	Logger  *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserHandler(args *UserHandler) *UserHandler {
	return args
}

func (h *UserHandler) Create(c *gin.Context) {
	var userRequest model.CreateUserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	userResponse, err := h.UseCase.Create(c, &userRequest)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(201, userResponse)
}

func (h *UserHandler) List(c *gin.Context) {
	var listUserRequest model.ListUserRequest
	if err := c.ShouldBindQuery(&listUserRequest); err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	if listUserRequest.Limit == 0 {
		listUserRequest.Limit = 10
	}
	if listUserRequest.Page == 0 {
		listUserRequest.Page = 1
	}

	result, err := h.UseCase.List(c, &listUserRequest)
	if err != nil {
		h.Logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	count, err := h.UseCase.Count(c, &listUserRequest)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	pagination := model.Pagination{
		CurrentPage:  listUserRequest.Page,
		TotalPages:   int32((count + int64(listUserRequest.Limit) - 1) / int64(listUserRequest.Limit)),
		TotalItems:   count,
		ItemsPerPage: listUserRequest.Limit,
	}

	c.JSON(200, model.ResponseData[model.ListUserResponse]{
		Code:   200,
		Status: "success",
		Data: model.ListUserResponse{
			Items:      *result,
			Pagination: pagination,
		},
	})
}

func (h *UserHandler) GetUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	result, err := h.UseCase.GetByID(c, int32(userId))
	if err != nil {
		h.Logger.Error(err)

		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}

		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(200, model.ResponseData[model.User]{Code: 200, Status: "success", Data: model.User{
		ID:       result.ID,
		Username: result.Username,
		Role:     result.Role,
	}})
}

func (h *UserHandler) UpdateUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	var userRequest model.UpdateUserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	result, err := h.UseCase.Update(c, &userRequest, int32(userId))
	if err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(200, model.ResponseData[*model.User]{Code: 200, Status: "success", Data: result})
}

func (h *UserHandler) DeleteUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid ID"})
		return
	}

	deletedUser, err := h.UseCase.Delete(c, int32(userId))
	if err != nil {
		h.Logger.Error(err)

		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(200, model.ResponseData[*model.User]{Code: 200, Status: "success", Data: deletedUser})
}
