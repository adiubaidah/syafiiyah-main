package persistence

import "context"

func (store *SQLStore) CreateArduinoWithModes(ctx context.Context, arduinoName string, modeParams []CreateArduinoModesParams) (Arduino, error) {
	var createdArduino Arduino

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		arduino, err := q.CreateArduino(ctx, arduinoName)
		if err != nil {
			return err
		}
		createdArduino = arduino

		for i := range modeParams {
			modeParams[i].ArduinoID = arduino.ID
		}
		_, err = q.CreateArduinoModes(ctx, modeParams)
		return err

	})
	return createdArduino, err

}
