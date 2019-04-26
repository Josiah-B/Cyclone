/*
	MemoryCache wraps the SQLite Database implementation and provides additional memory caching for the current sensor readings
*/
package MemoryCache

import (
	"fmt"

	"github.com/Josiah-B/Cyclone/CurrentDeviceData"
	"github.com/Josiah-B/Cyclone/Interfaces"
	"github.com/Josiah-B/Cyclone/SQLiteDatabase"
)

type Cache struct {
	//The database is used for long term storage of the sensor readings
	database SQLiteDatabase.DataBase
	//This is our memory cache for the current sensor readings
	currentData  CurrentDeviceData.Data
	DataBasePath string
}

func (cache *Cache) Initilize() {
	//cache.database = SQLiteDatabase.DataBase{}

	cache.database.Configuration = SQLiteDatabase.SqliteSettings{Path: cache.DataBasePath}
	cache.database.Initilize()
	cache.database.Open()

	//setup memory cache for the current weather conditions
	cache.currentData = CurrentDeviceData.Data{}
}

//GetCurrentSensorReadings returns the current sensor readings for the specified weather station
func (cache *Cache) GetCurrentSensorReadings(StationName string) Interfaces.StationUploadTemplate {
	return cache.currentData.GetCurrentSensorReadings(StationName)
}

//SetCurrentSensorReadings stores the current sensor readings in memory
func (cache *Cache) SetCurrentSensorReadings(currentConditions Interfaces.StationUploadTemplate) {
	cache.currentData.SetCurrentSensorReadings(&currentConditions)
}

// AddOrUpdateStation adds a new station to the database. If a station with the stn.StationID is already in the database this function will update it with the new values.
func (cache *Cache) AddOrUpdateStation(stn *Interfaces.Station) {
	cache.database.AddOrUpdateStation(stn)
}

//LogConditions pushes the current station to the database thus logging the conditions
func (cache *Cache) LogConditions(currentConditions *Interfaces.StationUploadTemplate) error {
	fmt.Println("MemCache: Logging Conditions", currentConditions)
	return cache.database.LogConditions(currentConditions)
}

func (cache *Cache) AddOrUpdateUnitType(unit *Interfaces.UnitType) {
	cache.database.AddOrUpdateUnitType(unit)
}

// AddOrUpdateSensor updates the current sensor or creates a new one if it does not exist. Returns the ID number fot the newly created sensor.
func (cache *Cache) AddOrUpdateSensor(sensor *Interfaces.Sensor) int64 {
	return cache.database.AddOrUpdateSensor(sensor)

}

func (cache *Cache) AddOrUpdateObservedProperty(observedProperty *Interfaces.ObservedProperty) {
	cache.database.AddOrUpdateObservedProperty(observedProperty)
}

func (cache *Cache) AddDataStream(dataStream *Interfaces.DataStream) int64 {
	return cache.database.AddDataStream(dataStream)
}

func (cache *Cache) UpdateDataStream(dataStream *Interfaces.DataStream) {
	cache.database.UpdateDataStream(dataStream)
}

func (cache *Cache) AddObservation(observation *Interfaces.Observation) {
	cache.database.AddObservation(observation)

}

// GetStation retrieves from the database a Station with ID number 'stationID'
func (cache *Cache) GetStation(stationID string) *Interfaces.Station {
	return cache.database.GetStation(stationID)
}

// GetStation retrieves from the database a Station with name 'stationName'
func (cache *Cache) GetStationByName(stationName string) *Interfaces.Station {
	temp, _ := cache.database.GetStationByName(stationName)
	return temp
}

// GetStations retrieves all of the stations from the database
func (cache *Cache) GetStations() []Interfaces.Station {
	var stns = cache.database.GetStations()

	for _, stnVal := range cache.currentData.SensorReadings {

		//see if the station is already in our current list
		stnInList := false

		for index := range stns {
			if stns[index].Name == stnVal.StationName {
				stnInList = true
				break
			}
		}

		// add the station to our list if it's not in there
		if stnInList == false {
			var tempStn = Interfaces.Station{
				StationID:   -1, //A Station id value of -1 indicates the station is not in the database at this point
				Name:        stnVal.StationName,
				Description: "Not set",
				Latitude:    "",
				Longitude:   ""}
			stns = append(stns, tempStn)
		}
	}
	return stns
}

// GetDataStream returns a datastream from the database based on the 'streamID'
func (cache *Cache) GetDataStream(streamID string) *Interfaces.DataStream {
	return cache.database.GetDataStream(streamID)
}

// GetDataStream returns a datastream from the database based on the 'sensorName'
func (cache *Cache) GetDataStreamBySensorName(sensorName string, stationIDNum int64) *Interfaces.DataStream {
	return cache.database.GetDataStreamBySensorName(sensorName, stationIDNum)
}

func (cache *Cache) GetDataStreams() *[]Interfaces.DataStream {
	return cache.database.GetDataStreams()
}

func (cache *Cache) GetUnitType(unitTypeID string) *Interfaces.UnitType {
	return cache.database.GetUnitType(unitTypeID)

}

func (cache *Cache) GetUnitTypes() *[]Interfaces.UnitType {
	return cache.database.GetUnitTypes()
}

func (cache *Cache) GetSensor(sensorID string) *Interfaces.Sensor {
	return cache.database.GetSensor(sensorID)
}

func (cache *Cache) GetSensors() *[]Interfaces.Sensor {
	return cache.database.GetSensors()
}

// GetObservations returns the obseration with 'ObservationID'
func (cache *Cache) GetObservation(observationID string) *Interfaces.Observation {
	return cache.database.GetObservation(observationID)
}

func (cache *Cache) GetObservations(parameters Interfaces.ObservationParameters) *[]Interfaces.Observation {
	return cache.database.GetObservations(parameters)
}

func (cache *Cache) GetObservedProperty(propertyID string) *Interfaces.ObservedProperty {
	return cache.database.GetObservedProperty(propertyID)
}

func (cache *Cache) GetObservedProperties() *[]Interfaces.ObservedProperty {
	return cache.database.GetObservedProperties()
}
