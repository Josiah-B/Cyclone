package CurrentDeviceData

import "github.com/Josiah-B/Cyclone/Interfaces"

type Data struct {
	SensorReadings map[string]Interfaces.StationUploadTemplate
}

func init() {
}

func (data *Data) GetCurrentSensorReadings(StationName string) Interfaces.StationUploadTemplate {
	return data.SensorReadings[StationName]
}

func (data *Data) SetCurrentSensorReadings(currentConditions *Interfaces.StationUploadTemplate) {
	//Create the data structure if it has not been initilized yet
	if data.SensorReadings == nil {
		data.SensorReadings = make(map[string]Interfaces.StationUploadTemplate)
	}

	data.SensorReadings[currentConditions.StationName] = *currentConditions
}
