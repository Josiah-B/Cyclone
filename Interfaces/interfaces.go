package Interfaces

import (
	"net/http"
	"time"
)

// Structures used inside the database. Each one represents a database table.
type UnitType struct {
	UnitTypeID    int
	Name          string
	UnitOfMeasure string
	Description   string
}

type Sensor struct {
	SensorID    int
	Name        string
	Description string
}

type ObservedProperty struct {
	PropertyID  int
	Name        string
	Description string
}

type DataStream struct {
	StationID          int
	StreamID           int
	ObservedPropertyID int
	SensorID           int
	UnitTypeID         int
	PropertyObserved   ObservedProperty
	SensorUnit         Sensor
}

type Observation struct {
	ObservationID int
	DataStreamID  int
	TimeStamp     time.Time
	Value         string
}

type Station struct {
	StationID   int
	Name        string
	Description string
	Latitude    string
	Longitude   string
	Streams     []DataStream
}

// StationUploadTemplate is a template for uploading weather station sensor readings to the api. It is also the format the api spits back current conditions in.
type StationUploadTemplate struct {
	StationName    string
	TimeStamp      time.Time
	SensorReadings map[string]string
}

// Template for configuration settings
type Configuration struct {
	Settings map[string]interface{}
}

// APIRoute holds info that is used to setup the API routing
type APIRoute struct {
	Route         string
	HandlerMethod func(http.ResponseWriter, *http.Request) `json:"-"` // we do not need any json exports to know what the handling method is
	HTTPMethod    string
	Description   string
}

//ObservationParameters holds the paramaters for getting historical graph data from the database
type ObservationParameters struct {
	SensorID  int
	StartTime time.Time
	EndTime   time.Time
}
type Storage interface {
	Initilize()
	// AddOrUpdateStation adds a new station to the database. If a station with the stn.StationID is already in the database this function will update it with the new values.
	AddOrUpdateStation(stn *Station)

	//LogConditions pushes the current station to the database thus logging the conditions
	LogConditions(currentConditions *StationUploadTemplate) error

	AddOrUpdateUnitType(unit *UnitType)

	// AddOrUpdateSensor updates the current sensor or creates a new one if it does not exist. Returns the ID number fot the newly created sensor.
	AddOrUpdateSensor(sensor *Sensor) int64
	AddOrUpdateObservedProperty(observedProperty *ObservedProperty)
	AddDataStream(dataStream *DataStream) int64
	UpdateDataStream(dataStream *DataStream)
	AddObservation(observation *Observation)

	// GetStation retrieves from the database a Station with ID number 'stationID'
	GetStation(stationID string) *Station

	// GetStation retrieves from the database a Station with name 'stationName'
	GetStationByName(stationName string) *Station

	// GetStations retrieves all of the stations from the database
	GetStations() []Station

	// GetDataStream returns a datastream from the database based on the 'streamID'
	GetDataStream(streamID string) *DataStream

	// GetDataStream returns a datastream from the database based on the 'sensorName'
	GetDataStreamBySensorName(sensorName string, stationIDNum int64) *DataStream
	GetDataStreams() *[]DataStream
	GetUnitType(unitTypeID string) *UnitType
	GetUnitTypes() *[]UnitType
	GetSensor(sensorID string) *Sensor
	GetSensors() *[]Sensor

	// GetObservations returns the obseration with 'ObservationID'
	GetObservation(observationID string) *Observation
	GetObservations(parameters ObservationParameters) *[]Observation
	GetObservedProperty(propertyID string) *ObservedProperty
	GetObservedProperties() *[]ObservedProperty

	GetCurrentSensorReadings(StationName string) StationUploadTemplate
	SetCurrentSensorReadings(currentSensorReadings StationUploadTemplate)
}

// new layer will have the following
/* works between:
api - logger - performs operations on the data and current weather conditions
database wrapper - caches in memory some frequently changing (or frequently requested) items such as the current sensor readings
database - holds the raw data and provides easier access by wrapping SQL commands as methods
*/
