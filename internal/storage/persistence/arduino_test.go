package persistence

import (
	"context"
	"fmt"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/random"
	"github.com/stretchr/testify/require"
)

func clearArduinoTable(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), "DELETE FROM arduino_mode")
	require.NoError(t, err)
	_, err = sqlStore.db.Exec(context.Background(), "DELETE FROM arduino")
	require.NoError(t, err)
}

func createRandomArduinoModesParams(name string) []CreateArduinoModesParams {
	return []CreateArduinoModesParams{
		{
			Mode:                 ArduinoModeTypeRecord,
			InputTopic:           fmt.Sprintf("%s/input/%s", name, ArduinoModeTypeRecord),
			AcknowledgementTopic: fmt.Sprintf("%s/acknowledgment/%s", name, ArduinoModeTypeRecord),
		},
		{
			Mode:                 ArduinoModeTypeExcuse,
			InputTopic:           fmt.Sprintf("%s/input/%s", name, ArduinoModeTypeExcuse),
			AcknowledgementTopic: fmt.Sprintf("%s/acknowledgment/%s", name, ArduinoModeTypeExcuse),
		},
		{
			Mode:                 ArduinoModeTypePresence,
			InputTopic:           fmt.Sprintf("%s/input/%s", name, ArduinoModeTypePresence),
			AcknowledgementTopic: fmt.Sprintf("%s/acknowledgment/%s", name, ArduinoModeTypePresence),
		},
	}
}

func TestCreateArduinoWithModes(t *testing.T) {
	nameRandom := random.RandomString(10)
	modeParams := createRandomArduinoModesParams(nameRandom)

	arduino, err := sqlStore.CreateArduinoWithModes(context.Background(), nameRandom, modeParams)
	require.NoError(t, err)
	require.NotEmpty(t, arduino)
	require.NotZero(t, arduino.ID)
	require.Equal(t, nameRandom, arduino.Name)

}
