package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	pb "github.com/adiubaidah/syafiiyah-main/internal/protobuf"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SantriScheduleHandler interface {
	CreateSantriScheduleHandler(c *gin.Context)
	ListSantriScheduleHandler(c *gin.Context)
	UpdateSantriScheduleHandler(c *gin.Context)
	DeleteSantriScheduleHandler(c *gin.Context)
}

type santriScheduleHandler struct {
	logger  *logrus.Logger
	service pb.SantriScheduleServiceClient
}

func NewSantriScheduleHandler(logger *logrus.Logger, service pb.SantriScheduleServiceClient) SantriScheduleHandler {
	return &santriScheduleHandler{
		logger:  logger,
		service: service,
	}
}

func (h *santriScheduleHandler) CreateSantriScheduleHandler(c *gin.Context) {
	var santriScheduleRequest model.CreateSantriScheduleRequest
	if err := c.ShouldBindJSON(&santriScheduleRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	resp, err := h.service.CreateSantriSchedule(context.Background(), &pb.CreateSantriScheduleRequest{
		Name:          santriScheduleRequest.Name,
		Description:   santriScheduleRequest.Description,
		StartPresence: santriScheduleRequest.StartPresence,
		StartTime:     santriScheduleRequest.StartTime,
		FinishTime:    santriScheduleRequest.FinishTime,
	})
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, model.ResponseData[*model.SantriScheduleResponse]{Code: http.StatusCreated, Status: "Created", Data: &model.SantriScheduleResponse{
		ID:            int32(resp.Id),
		Name:          resp.Name,
		Description:   resp.Description,
		StartPresence: resp.StartPresence,
		StartTime:     resp.StartTime,
		FinishTime:    resp.FinishTime,
	}})
}

func (h *santriScheduleHandler) ListSantriScheduleHandler(c *gin.Context) {
	resp, err := h.service.ListSantriSchedule(context.Background(), &pb.ListSantriScheduleRequest{})
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}
	var santriScheduleResponses []model.SantriScheduleResponse
	for _, santriSchedule := range resp.Schedules {
		santriScheduleResponses = append(santriScheduleResponses, model.SantriScheduleResponse{
			ID:            int32(santriSchedule.Id),
			Name:          santriSchedule.Name,
			Description:   santriSchedule.Description,
			StartPresence: santriSchedule.StartPresence,
			StartTime:     santriSchedule.StartTime,
			FinishTime:    santriSchedule.FinishTime,
		})
	}
	c.JSON(http.StatusOK, model.ResponseData[*[]model.SantriScheduleResponse]{Code: http.StatusOK, Status: "OK", Data: &santriScheduleResponses})
}

func (h *santriScheduleHandler) UpdateSantriScheduleHandler(c *gin.Context) {
	var santriScheduleRequest model.UpdateSantriScheduleRequest
	if err := c.ShouldBindJSON(&santriScheduleRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	idParam := c.Param("id")
	santriScheduleId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	resp, err := h.service.UpdateSantriSchedule(context.Background(), &pb.UpdateSantriScheduleRequest{
		Schedule: &pb.SantriSchedule{
			Id:            int32(santriScheduleId),
			Name:          santriScheduleRequest.Name,
			Description:   santriScheduleRequest.Description,
			StartPresence: santriScheduleRequest.StartPresence,
			StartTime:     santriScheduleRequest.StartTime,
			FinishTime:    santriScheduleRequest.FinishTime,
		},
	})

	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ResponseData[*model.SantriScheduleResponse]{Code: http.StatusOK, Status: "OK", Data: &model.SantriScheduleResponse{
		ID:            int32(resp.Id),
		Name:          resp.Name,
		Description:   resp.Description,
		StartPresence: resp.StartPresence,
		StartTime:     resp.StartTime,
		FinishTime:    resp.FinishTime,
	}})

}

func (h *santriScheduleHandler) DeleteSantriScheduleHandler(c *gin.Context) {
	idParam := c.Param("id")
	santriScheduleId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	resp, err := h.service.DeleteSantriSchedule(context.Background(), &pb.DeleteSantriScheduleRequest{
		Id: int32(santriScheduleId),
	})

	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ResponseData[*model.SantriScheduleResponse]{Code: http.StatusOK, Status: "OK", Data: &model.SantriScheduleResponse{
		ID:            int32(resp.Id),
		Name:          resp.Name,
		Description:   resp.Description,
		StartPresence: resp.StartPresence,
		StartTime:     resp.StartTime,
		FinishTime:    resp.FinishTime,
	}})
}
