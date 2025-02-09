package handler

import (
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	repo "github.com/adiubaidah/rfid-syafiiyah/internal/repository"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProfileHandler struct {
	Logger          *logrus.Logger
	EmployeeUseCase *usecase.EmployeeUseCase
	ParentUseCase   *usecase.ParentUseCase
}

func NewProfileHandler(args *ProfileHandler) *ProfileHandler {
	return args
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userValue, _ := c.Get("user")
	user, _ := userValue.(*model.User)

	if user.Role == repo.RoleTypeParent {
		parent, err := h.ParentUseCase.GetByUserID(c, user.ID)
		if err != nil {
			h.Logger.Error(err)
			if appErr, ok := err.(*exception.AppError); ok {
				c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
				return
			}
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		}
		c.JSON(200, model.ResponseData[*model.ParentResponse]{Code: 200, Status: "success", Data: parent})
	} else {
		employee, err := h.EmployeeUseCase.GetByUserID(c, user.ID)
		if err != nil {
			h.Logger.Error(err)
			if appErr, ok := err.(*exception.AppError); ok {
				c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
				return
			}
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		}
		c.JSON(200, model.ResponseData[*model.Employee]{Code: 200, Status: "success", Data: employee})
	}
}
