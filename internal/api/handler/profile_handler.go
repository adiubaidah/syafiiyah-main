package handler

import (
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProfileHandler interface {
	GetProfile(c *gin.Context)
}

type profileHandler struct {
	logger          *logrus.Logger
	employeeUseCase usecase.EmployeeUseCase
	parentUseCase   usecase.ParentUseCase
}

func NewProfileHandler(logger *logrus.Logger, employeeUseCase usecase.EmployeeUseCase, parentUseCase usecase.ParentUseCase) ProfileHandler {
	return &profileHandler{
		logger:          logger,
		employeeUseCase: employeeUseCase,
		parentUseCase:   parentUseCase,
	}
}

func (h *profileHandler) GetProfile(c *gin.Context) {
	userValue, _ := c.Get("user")
	user, _ := userValue.(*model.User)

	if user.Role == db.RoleTypeParent {
		parent, err := h.parentUseCase.GetParentByUserID(c, user.ID)
		if err != nil {
			h.logger.Error(err)
			if appErr, ok := err.(*exception.AppError); ok {
				c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
				return
			}
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		}
		c.JSON(200, model.ResponseData[*model.ParentResponse]{Code: 200, Status: "success", Data: parent})
	} else {
		employee, err := h.employeeUseCase.GetEmployeeByUserID(c, user.ID)
		if err != nil {
			h.logger.Error(err)
			if appErr, ok := err.(*exception.AppError); ok {
				c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
				return
			}
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		}
		c.JSON(200, model.ResponseData[*model.Employee]{Code: 200, Status: "success", Data: employee})
	}
}
