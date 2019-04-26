package SQLiteDatabase

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Josiah-B/Cyclone/Interfaces"
	// this is needed for the sqlite3 database to function properly
	_ "github.com/mattn/go-sqlite3"
)

// createDBquery is an SQL query for creating the initial database
// the 'Settings' Table is for version control
const createDBQuery string = `CREATE TABLE 'Settings'
(
	Version INTEGER
);

CREATE TABLE 'Station'
(
	StationID INTEGER PRIMARY KEY AUTOINCREMENT,
	Name TEXT,
	Description TEXT,
	Latitude TEXT,
	Longitude TEXT
);

CREATE TABLE 'UnitType'
(
	UnitTypeID INTEGER PRIMARY KEY AUTOINCREMENT,
	Name TEXT,
	UnitOfMeasure TEXT,
	Description TEXT
);

CREATE TABLE 'Sensor'
(
	SensorID INTEGER PRIMARY KEY AUTOINCREMENT,
	Name TEXT,
	Description TEXT
);

CREATE TABLE 'ObserverdProperty'
(
	PropertyID INTEGER PRIMARY KEY AUTOINCREMENT,
	Name TEXT,
	Description TEXT
);

CREATE TABLE 'DataStream'
(
	StreamID INTEGER PRIMARY KEY AUTOINCREMENT,
	StationID INTEGER REFERENCES Station(StationID),
	ObservedPropertyID INTEGER REFERENCES ObservedProperty(ObservedPropertyID),
	SensorID INTEGER REFERENCES Sensor(SensorID),
	UnitTypeID INTEGER REFERENCES UnitType(UnitTypeID)
);

CREATE TABLE 'Observation'
(
	ObservationID INTEGER PRIMARY KEY AUTOINCREMENT,
	DataStreamID INTEGER REFERENCES DataStream(DataStreamID),
	TimeStamp DATETIME NOT NULL,
	Value TEXT NOT NULL
);

INSERT INTO Settings (Version) VALUES (1);`

//

// SqliteSettings holds connection related settings for the database
type SqliteSettings struct {
	//Path Defines the filepath to the database
	Path string
}

// DataBase is the type class for the SQLite3 database implementation
type DataBase struct {
	// Configuration contains the various connection settings for the database
	Configuration SqliteSettings
	// backingDB is the actual database interface
	BackingDB *sql.DB
	// databaseError is where we can put any errors we run into during operation. This makes it easier to check for errors at any point in the programs excution
	databaseError error
	// Version is the current version of the database
	version int
}

func (db *DataBase) Initilize() {
	// initilize array that will hold the weather stations current conditions
}

//Open connects to the database, creating it if nessasary
func (db *DataBase) Open() {
	db.BackingDB, db.databaseError = sql.Open("sqlite3", db.Configuration.Path)
	if db.databaseError != nil {
		fmt.Println(db.databaseError)
	} else {
		db.createAndUpgrade()
	}

}

func (db *DataBase) createAndUpgrade() {
	db.version = db.getSchemaVersion() // get the current database schema version

	if db.version < 1 { // if the database schema version is < 1 we know the tables have not been created yet
		db.create()                        // create the database tables
		db.version = db.getSchemaVersion() // refresh the database version with the current; this also allows os to verify that the create operation succeded
	}
	fmt.Printf("DataBase Version is: %v\n", db.version)
}

// Create creates the database
func (db *DataBase) create() {
	fmt.Println("Creating Database...")
	_, db.databaseError = db.BackingDB.Exec(createDBQuery)
	if db.databaseError != nil {
		fmt.Println(db.databaseError)
	}
}

func (db *DataBase) getSchemaVersion() int {
	var version = -1

	var row *sql.Row
	row = db.BackingDB.QueryRow("SELECT Version FROM Settings LIMIT 1")

	db.databaseError = row.Scan(&version)
	// if the 'Select version' query throws an error then the database has not been created
	if db.databaseError != nil {
		db.databaseError = nil // clear the database error
		version = -1           // set the database version to -1 meaning we do not know
		fmt.Println(db.databaseError)
	}
	return version
}

// Close closes the database flushing any changes to disk
func (db *DataBase) Close() {
	db.BackingDB.Close()
}

// Upgrade checks the database schema version and (if needed) upgrades it to the latest version
func (db *DataBase) Upgrade() {

}

// AddOrUpdateStation adds a new station to the database. If a station with the stn.StationID is already in the database this function will update it with the new values.
func (db *DataBase) AddOrUpdateStation(stn *Interfaces.Station) {
	tx, err := db.BackingDB.Begin()
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO Station (Name, Description, Latitude, Longitude) VALUES ((SELECT StationID FROM Station WHERE StationID = ?),?,?,?,?)")
	stmt.Exec(stn.StationID, stn.Name, stn.Description, stn.Latitude, stn.Longitude)

	for i := 0; i < len(stn.Streams); i++ {
		stmt, err = tx.Prepare("INSERT OR REPLACE INTO DataStream (StationID, StreamID, ObservedPropertyID, SensorID, UnitTypeID) VALUES (?,(SELECT StreamID FROM DataStreams WHERE StreamID = ?),?,?,?)")
		stmt.Exec(&stn.StationID, &stn.Streams[i].StreamID, &stn.Streams[i].ObservedPropertyID, &stn.Streams[i].SensorID, &stn.Streams[i].UnitTypeID)
	}

	// rollback the operation if there was an error
	if err != nil {
		tx.Rollback()
	} else { // commit the changes if there were no errors
		tx.Commit()
	}

}

func (db *DataBase) addStation(stn *Interfaces.StationUploadTemplate) error {
	tx, err := db.BackingDB.Begin()
	stmt, err := tx.Prepare("INSERT INTO Station (Name, Description, Latitude, Longitude) VALUES (?,?,?,?)")
	// rollback the operation if there was an error
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}

	res, err2 := stmt.Exec(stn.StationName, "none", "n/a", "n/a")

	if err2 != nil {
		tx.Rollback()
		return err2
	}

	// commit the changes if there were no errors
	tx.Commit()
	fmt.Println(res)
	return nil

}

//LogConditions pushes the current station to the database thus logging the conditions
func (db *DataBase) LogConditions(currentConditions *Interfaces.StationUploadTemplate) error {

	currentStation, err := db.GetStationByName(currentConditions.StationName)
	fmt.Println(err)
	fmt.Print("Station ")

	if currentStation == nil {
		fmt.Print(currentConditions.StationName)
		fmt.Println("does not exist. Attempting to create it...")
		err := db.addStation(currentConditions)
		if err != nil {
			return errors.New("Cannot Find Station; failed to add it to the database")
		}
		currentStation, err = db.GetStationByName(currentConditions.StationName)

		fmt.Println(currentStation, err)
		fmt.Println("Station created!")
	}

	fmt.Println(currentStation.StationID)

	for sensorName, value := range currentConditions.SensorReadings {
		var currentDataStream = db.GetDataStreamBySensorName(sensorName, int64(currentStation.StationID))
		var dataStreamID = -1
		fmt.Print("Current Data Stream: ")
		fmt.Println(currentDataStream)
		if currentDataStream != nil {
			dataStreamID = currentDataStream.StreamID
		} else {
			// create a new datastream if one does not exist for this sensor
			dataStreamID = int(db.createDataStream(currentStation.StationID, sensorName, ""))
		}

		observ := Interfaces.Observation{
			DataStreamID: dataStreamID,
			TimeStamp:    currentConditions.TimeStamp,
			Value:        value}

		// push the current observation to the DB
		db.AddObservation(&observ)
	}
	return nil
}

// creates a new dataStream and returns the StreamID
func (db *DataBase) createDataStream(stationID int, sensorName string, sensorDescription string) int64 {
	// Create the sensor
	tempSensor := Interfaces.Sensor{Name: sensorName, Description: sensorDescription}
	sensorID := db.AddOrUpdateSensor(&tempSensor)

	// Create the observed Property... which can be created later

	// create the DataStream
	dataStream := Interfaces.DataStream{
		StationID: stationID,
		SensorID:  int(sensorID)}
	fmt.Println("DataStream:")
	fmt.Println(dataStream)
	return db.AddDataStream(&dataStream)
}

func (db *DataBase) AddOrUpdateUnitType(unit *Interfaces.UnitType) {
	var updateQuery = "INSERT OR REPLACE INTO UnitType (UnitTypeID,Name,UnitOfMeasure,Description) VALUES ((SELECT UnitTypeID FROM UnitType WHERE UnitTypeID = ?),?,?,?)"
	db.BackingDB.Exec(updateQuery, &unit.UnitTypeID, &unit.Name, &unit.UnitOfMeasure, &unit.Description)
}

// AddOrUpdateSensor updates the current sensor or creates a new one if it does not exist. Returns the ID number fot the newly created sensor.
func (db *DataBase) AddOrUpdateSensor(sensor *Interfaces.Sensor) int64 {
	var updateQuery = "INSERT OR REPLACE INTO Sensor (SensorID, Name, Description) VALUES ((SELECT SensorID FROM Sensor WHERE SensorID = ?),?,?)"
	res, _ := db.BackingDB.Exec(updateQuery, &sensor.SensorID, &sensor.Name, &sensor.Description)
	id, _ := res.LastInsertId()
	return id

}

func (db *DataBase) AddOrUpdateObservedProperty(observedProperty *Interfaces.ObservedProperty) {
	var updateQuery = "INSERT OR REPLACE INTO ObservedProperty (PropertyID, Name, Description) VALUES ((SELECT PropertyID FROM ObservedProprty WHERE PropertyID = ?),?,?)"
	db.BackingDB.Exec(updateQuery, &observedProperty.PropertyID, &observedProperty.Name, &observedProperty.Description)
}

func (db *DataBase) AddDataStream(dataStream *Interfaces.DataStream) int64 {
	var updateQuery = "INSERT INTO DataStream (StationID, ObservedPropertyID, SensorID, UnitTypeID) VALUES (?,?,?,?)"
	res, _ := db.BackingDB.Exec(updateQuery, &dataStream.StationID, &dataStream.ObservedPropertyID, &dataStream.SensorID, &dataStream.UnitTypeID)
	id, _ := res.LastInsertId()
	return id
}

func (db *DataBase) UpdateDataStream(dataStream *Interfaces.DataStream) {
	var updateQuery = "REPLACE INTO DataStream (StationID, StreamID, ObservedPropertyID, SensorID, UnitTypeID) VALUES (?,(SELECT StreamID FROM DataStream WHERE StreamID = ?),?,?,?,?)"
	db.BackingDB.Exec(updateQuery, &dataStream.StationID, &dataStream.StreamID, &dataStream.ObservedPropertyID, &dataStream.SensorID, &dataStream.UnitTypeID)
}

func (db *DataBase) AddObservation(observation *Interfaces.Observation) {
	var updateQuery = "INSERT INTO Observation (DataStreamID, TimeStamp, Value) VALUES (?,?,?)"
	db.BackingDB.Exec(updateQuery, &observation.DataStreamID, &observation.TimeStamp, &observation.Value)
}

// GetStation retrieves from the database a Station with ID number 'stationID'
func (db *DataBase) GetStation(stationID string) *Interfaces.Station {
	var stn Interfaces.Station
	row := db.BackingDB.QueryRow("SELECT StationID, Name, Description, Latitude, Longitude FROM Station WHERE StationID = ?", stationID)
	err := row.Scan(&stn.StationID, &stn.Name, &stn.Description, &stn.Latitude, &stn.Longitude)

	if err != nil {
		return nil
	}

	return &stn
}

// GetStation retrieves from the database a Station with name 'stationName'
func (db *DataBase) GetStationByName(stationName string) (*Interfaces.Station, error) {
	fmt.Println("Attempting to get the following station from the database: ", stationName)
	var stn Interfaces.Station
	row := db.BackingDB.QueryRow("SELECT StationID, Name, Description, Latitude, Longitude FROM Station WHERE Name = ?", stationName)
	err := row.Scan(&stn.StationID, &stn.Name, &stn.Description, &stn.Latitude, &stn.Longitude)
	fmt.Println("Results of the database search: ", stn, err)
	if err != nil {
		return nil, err
	}

	return &stn, nil
}

// GetStations retrieves all of the stations from the database
func (db *DataBase) GetStations() []Interfaces.Station {
	var stns []Interfaces.Station

	rows, err := db.BackingDB.Query("SELECT StationID, Name, Description, Latitude, Longitude FROM Station")

	for rows.Next() {
		var stn = *new(Interfaces.Station)

		err = rows.Scan(&stn.StationID, &stn.Name, &stn.Description, &stn.Latitude, &stn.Longitude)

		stns = append(stns, stn)
	}

	if err != nil {
		return nil
	}

	return stns
}

// GetDataStream returns a datastream from the database based on the 'streamID'
func (db *DataBase) GetDataStream(streamID string) *Interfaces.DataStream {
	var dataStream Interfaces.DataStream
	row := db.BackingDB.QueryRow("SELECT StationID, StreamID, ObservedPropertyID, SensorID, UnitTypeID FROM DataStream WHERE StreamID = ?", streamID)
	err := row.Scan(&dataStream.StationID, &dataStream.StreamID, &dataStream.ObservedPropertyID, &dataStream.SensorID, &dataStream.UnitTypeID)

	if err != nil {
		return nil
	}

	return &dataStream
}

// GetDataStream returns a datastream from the database based on the 'sensorName'
func (db *DataBase) GetDataStreamBySensorName(sensorName string, stationIDNum int64) *Interfaces.DataStream {
	var dataStream Interfaces.DataStream
	row := db.BackingDB.QueryRow("SELECT StationID, StreamID, ObservedPropertyID, DataStream.SensorID, UnitTypeID FROM DataStream INNER JOIN Sensor ON DataStream.SensorID = Sensor.SensorID WHERE Sensor.Name = ? AND DataStream.StationID = ?", sensorName, stationIDNum)
	err := row.Scan(&dataStream.StationID, &dataStream.StreamID, &dataStream.ObservedPropertyID, &dataStream.SensorID, &dataStream.UnitTypeID)

	if err != nil {
		fmt.Println("DataBase: GetDataStreamBySensorName: ", err)
		return nil
	}

	return &dataStream
}

func (db *DataBase) GetDataStreams() *[]Interfaces.DataStream {
	var streams []Interfaces.DataStream

	rows, err := db.BackingDB.Query("SELECT StationID, StreamID, ObservedPropertyID, SensorID, UnitTypeID FROM DataStream")

	if err != nil {
		fmt.Println(err)
		return nil
	}

	for rows.Next() {
		var dataStream = *new(Interfaces.DataStream)
		err = rows.Scan(&dataStream.StationID, &dataStream.StreamID, &dataStream.ObservedPropertyID, &dataStream.SensorID, &dataStream.UnitTypeID)
		streams = append(streams, dataStream)
	}

	return &streams
}

func (db *DataBase) GetUnitType(unitTypeID string) *Interfaces.UnitType {
	var unittype = Interfaces.UnitType{}

	row := db.BackingDB.QueryRow("SELECT UnitTypeID, Name, UnitOfMeasure, Description FROM UnitType WHERE UnitTypeID = ?", unitTypeID)
	err := row.Scan(&unittype.UnitTypeID, &unittype.Name, &unittype.UnitOfMeasure, &unittype.Description)
	if err != nil {
		return nil
	}

	return &unittype

}

func (db *DataBase) GetUnitTypes() *[]Interfaces.UnitType {
	var stns []Interfaces.UnitType

	rows, err := db.BackingDB.Query("SELECT UnitTypeID, Name, UnitOfMeasure, Description FROM UnitType")

	if err != nil {
		return nil
	}

	for rows.Next() {
		var unittype = *new(Interfaces.UnitType)

		err = rows.Scan(&unittype.UnitTypeID, &unittype.Name, &unittype.UnitOfMeasure, &unittype.Description)

		stns = append(stns, unittype)
	}

	return &stns
}

func (db *DataBase) GetSensor(sensorID string) *Interfaces.Sensor {
	var sensor = Interfaces.Sensor{}

	row := db.BackingDB.QueryRow("SELECT SensorID, Name, Description FROM Sensor WHERE SensorID = ?", sensorID)
	err := row.Scan(&sensor.SensorID, &sensor.Name, &sensor.Description)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &sensor
}

func (db *DataBase) GetSensors() *[]Interfaces.Sensor {
	var sensors []Interfaces.Sensor

	rows, err := db.BackingDB.Query("SELECT SensorID, Name, Description FROM Sensor")

	if err != nil {
		fmt.Println(err)
		return nil
	}

	for rows.Next() {
		var sensor = *new(Interfaces.Sensor)
		err = rows.Scan(&sensor.SensorID, &sensor.Name, &sensor.Description)
		sensors = append(sensors, sensor)
	}

	return &sensors
}

// GetObservations returns the obseration with 'ObservationID'
func (db *DataBase) GetObservation(observationID string) *Interfaces.Observation {
	var selectQuery = "SELECT ObservationID, DataStreamID, TimeStamp, Value FROM Observation WHERE ObservationID = ?"

	row := db.BackingDB.QueryRow(selectQuery, observationID)

	var obser = Interfaces.Observation{}

	err := row.Scan(&obser.ObservationID, &obser.DataStreamID, &obser.TimeStamp, &obser.Value)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &obser
}

func (db *DataBase) GetObservations(parameters Interfaces.ObservationParameters) *[]Interfaces.Observation {
	var observations []Interfaces.Observation

	rows, err := db.BackingDB.Query("SELECT ObservationID, DataStreamID, TimeStamp, Value FROM Observation INNER JOIN DataStream ON DataStream.StreamID = Observation.DataStreamID WHERE DataStream.SensorID = ?", parameters.SensorID)

	if err != nil {
		return nil
	}

	for rows.Next() {
		var observation = *new(Interfaces.Observation)
		err = rows.Scan(&observation.ObservationID, &observation.DataStreamID, &observation.TimeStamp, &observation.Value)
		observations = append(observations, observation)
	}

	return &observations
}

func (db *DataBase) GetObservedProperty(propertyID string) *Interfaces.ObservedProperty {
	var property = Interfaces.ObservedProperty{}

	row := db.BackingDB.QueryRow("SELECT PropertyID, Name, Description FROM ObserverdProperty WHERE PropertyID = ?", propertyID)
	err := row.Scan(&property.PropertyID, &property.Name, &property.Description)
	if err != nil {
		return nil
	}
	return &property
}

func (db *DataBase) GetObservedProperties() *[]Interfaces.ObservedProperty {
	var properties []Interfaces.ObservedProperty

	rows, err := db.BackingDB.Query("SELECT PropertyID, Name, Description FROM ObserverdProperty")

	if err != nil {
		return nil
	}

	for rows.Next() {
		var property = *new(Interfaces.ObservedProperty)
		err = rows.Scan(&property.PropertyID, &property.Name, &property.Description)
		properties = append(properties, property)
	}

	return &properties
}
