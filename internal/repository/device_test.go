package persistence

import (
	"context"
	"fmt"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/random"
	"github.com/stretchr/testify/require"
)

func clearDeviceTable(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), "DELETE FROM device_mode")
	require.NoError(t, err)
	_, err = sqlStore.db.Exec(context.Background(), "DELETE FROM device")
	require.NoError(t, err)
}

func createRandomArduinoModesParams(name string) []CreateDeviceModesParams {
	return []CreateDeviceModesParams{
		{
			Mode:                 DeviceModeTypePermission,
			InputTopic:           fmt.Sprintf("%s/input/%s", name, DeviceModeTypePermission),
			AcknowledgementTopic: fmt.Sprintf("%s/acknowledgment/%s", name, DeviceModeTypePermission),
		},
		{
			Mode:                 DeviceModeTypePing,
			InputTopic:           fmt.Sprintf("%s/input/%s", name, DeviceModeTypePing),
			AcknowledgementTopic: fmt.Sprintf("%s/acknowledgment/%s", name, DeviceModeTypePing),
		},
		{
			Mode:                 DeviceModeTypePresence,
			InputTopic:           fmt.Sprintf("%s/input/%s", name, DeviceModeTypePresence),
			AcknowledgementTopic: fmt.Sprintf("%s/acknowledgment/%s", name, DeviceModeTypePresence),
		},
	}
}

func TestCreateArduinoWithModes(t *testing.T) {
	nameRandom := random.RandomString(10)
	modeParams := createRandomArduinoModesParams(nameRandom)

	arduino, err := sqlStore.CreateDeviceWithModes(context.Background(), nameRandom, modeParams)
	require.NoError(t, err)
	require.NotEmpty(t, arduino)
	require.NotZero(t, arduino.ID)
	require.Equal(t, nameRandom, arduino.Name)

}
