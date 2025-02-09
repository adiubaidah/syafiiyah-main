package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	repo "github.com/adiubaidah/rfid-syafiiyah/internal/repository"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
)

type DeviceUseCase struct {
	store repo.Store
}

func NewDeviceUseCase(store repo.Store) *DeviceUseCase {
	return &DeviceUseCase{store: store}
}

func (c *DeviceUseCase) CreateDevice(ctx context.Context, request *model.CreateDeviceRequest) (*model.DeviceResponse, error) {
	sqlStore := c.store.(*repo.SQLStore)
	modeParams := make([]repo.CreateDeviceModesParams, 0)

	for _, mode := range request.Modes {
		modeParams = append(modeParams, repo.CreateDeviceModesParams{
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

func (c *DeviceUseCase) ListDevices(ctx context.Context) (*[]model.DeviceWithModesResponse, error) {
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

func (c *DeviceUseCase) UpdateDevice(ctx context.Context, request *model.CreateDeviceRequest, deviceId int32) (*model.DeviceResponse, error) {
	sqlStore := c.store.(*repo.SQLStore)
	modeParams := make([]repo.CreateDeviceModesParams, 0)

	for _, mode := range request.Modes {
		modeParams = append(modeParams, repo.CreateDeviceModesParams{
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

func (c *DeviceUseCase) DeleteDevice(ctx context.Context, deviceId int32) (*model.DeviceResponse, error) {
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
