package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
)

type DeviceUseCase interface {
	CreateDevice(ctx context.Context, request *model.CreateDeviceRequest) (*model.DeviceResponse, error)
	ListDevices(ctx context.Context) (*[]model.DeviceWithModesResponse, error)
	UpdateDevice(ctx context.Context, request *model.CreateDeviceRequest, deviceId int32) (*model.DeviceResponse, error)
	DeleteDevice(ctx context.Context, deviceId int32) (*model.DeviceResponse, error)
}

type deviceService struct {
	store db.Store
}

func NewDeviceUseCase(store db.Store) DeviceUseCase {
	return &deviceService{store: store}
}

func (c *deviceService) CreateDevice(ctx context.Context, request *model.CreateDeviceRequest) (*model.DeviceResponse, error) {
	sqlStore := c.store.(*db.SQLStore)
	modeParams := make([]db.CreateDeviceModesParams, 0)

	for _, mode := range request.Modes {
		modeParams = append(modeParams, db.CreateDeviceModesParams{
			Mode:                 mode,
			InputTopic:           fmt.Sprintf("%s/input/%s", util.ToSnakeCase(request.Name), mode),
			AcknowledgementTopic: fmt.Sprintf("%s/acknowledgment/%s", util.ToSnakeCase(request.Name), mode),
		})
	}

	device, err := sqlStore.CreateDeviceWithModes(ctx, request.Name, modeParams)
	if err != nil {
		return nil, err
	}

	return &model.DeviceResponse{
		ID:   device.ID,
		Name: device.Name,
	}, nil
}

func (c *deviceService) ListDevices(ctx context.Context) (*[]model.DeviceWithModesResponse, error) {
	devices, err := c.store.ListDevices(ctx)
	if err != nil {
		return nil, err
	}

	deviceMap := make(map[int32]*model.DeviceWithModesResponse)

	for _, device := range devices {
		if _, exists := deviceMap[device.ID]; !exists {
			deviceMap[device.ID] = &model.DeviceWithModesResponse{
				ID:    device.ID,
				Name:  device.Name,
				Modes: []model.DeviceMode{},
			}
		}

		if device.DeviceModeID.Valid {
			mode := model.DeviceMode{
				Mode:                device.DeviceModeMode.DeviceModeType,
				InputTopic:          device.DeviceModeInputTopic.String,
				AcknowledgmentTopic: device.DeviceModeAcknowledgementTopic.String,
			}
			deviceMap[device.ID].Modes = append(deviceMap[device.ID].Modes, mode)
		}
	}
	var responses []model.DeviceWithModesResponse
	for _, device := range deviceMap {
		responses = append(responses, *device)
	}

	return &responses, nil
}

func (c *deviceService) UpdateDevice(ctx context.Context, request *model.CreateDeviceRequest, deviceId int32) (*model.DeviceResponse, error) {
	sqlStore := c.store.(*db.SQLStore)
	modeParams := make([]db.CreateDeviceModesParams, 0)

	for _, mode := range request.Modes {
		modeParams = append(modeParams, db.CreateDeviceModesParams{
			Mode:                 mode,
			InputTopic:           fmt.Sprintf("%s/input/%s", util.ToSnakeCase(request.Name), mode),
			AcknowledgementTopic: fmt.Sprintf("%s/acknowledgment/%s", util.ToSnakeCase(request.Name), mode),
		})
	}

	device, err := sqlStore.UpdateDeviceWithModes(ctx, deviceId, request.Name, modeParams)
	if err != nil {
		return nil, err
	}

	return &model.DeviceResponse{
		ID:   device.ID,
		Name: device.Name,
	}, nil
}

func (c *deviceService) DeleteDevice(ctx context.Context, deviceId int32) (*model.DeviceResponse, error) {
	device, err := c.store.DeleteDevice(ctx, deviceId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Device not found")
		}
		return nil, err
	}

	return &model.DeviceResponse{
		ID:   device.ID,
		Name: device.Name,
	}, nil
}
