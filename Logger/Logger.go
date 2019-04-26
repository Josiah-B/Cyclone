package Logger

import (
	"fmt"
	"math"
	"time"

	"github.com/Josiah-B/Cyclone/Interfaces"
)

type Logger struct {
	//Logging frequency (in minutes)
	Interval int

	//the time at which we need to log conditions again (gets set to the (current time + logging interval) every time conditions are logged)
	nextLogTime time.Time

	//hold a listing of weather stations and the last time each one was logged
	loggingInfo LoggingInfo

	//need the ability to add or remove stations from the list of stations, based on the stations that are currently uploading.
	data Interfaces.Storage
}

//LogginMetaInfo contains the logging settings for the stations
type LoggingMetaInfo struct {
	stationName    string
	stationID      int
	lastLogTime    time.Time //the time at which we last logged the sensor data
	loggingEnabled bool
}

type LoggingInfo struct {
	info map[string]LoggingMetaInfo
}

//InitilizeLogger runs the setup and starts the needed timers
func (logger *Logger) Initilize(storage Interfaces.Storage) {
	fmt.Println("Logger initilize: storage: ", storage)
	logger.data = storage

	logger.loggingInfo.info = make(map[string]LoggingMetaInfo)

	go logger.startTimers()
}

func (logger *Logger) startTimers() {
	// check periodically to see if any new weather stations have been added to the program
	ticker := time.NewTicker(15 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				logger.addUploadingStationsToListForLogging()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	//START TIMER THAT WILL HANDLE LOGGING
	//get the amount of time we need to wait till the starting interval
	secsTilInterval := logger.getSecondsUntillInterval()
	fmt.Println("Logger: Seconds till next interval: ", secsTilInterval)

	hourTimerInterval := time.Second * time.Duration(secsTilInterval)
	fmt.Println("Logger: time to next interval ", hourTimerInterval)
	// wait until the hour based interval
	timer := time.NewTimer(hourTimerInterval)
	<-timer.C
	logger.LogStations()

	tempInterval := (time.Minute * time.Duration(logger.Interval))
	fmt.Println("Logger: Logging interval in millisecs: ", tempInterval)

	loggingTicker := time.NewTicker(tempInterval)
	loggingQuit := make(chan struct{})
	go func() {
		for {
			select {
			case <-loggingTicker.C:
				logger.LogStations()
			case <-loggingQuit:
				ticker.Stop()
				return
			}
		}
	}()
}

//getSecondsUntilInterval returns the number of seconds until the next interval; intervals are based on the hour mark
func (logger *Logger) getSecondsUntillInterval() int {
	//interval in seconds = 60 * interval
	intervalInSecs := float64(60 * logger.Interval)

	//get the current time
	now := time.Now()

	//get the number of seconds that have elapsed in this hour
	nSecs := float64((now.Minute() * 60) + now.Second())
	//number of seconds to the next quarter-hour mark
	delta := math.Ceil(nSecs/intervalInSecs)*intervalInSecs - nSecs
	return int(delta)
}

//GetLoggingSettings returns the logging settings for all the weather stations
func (logger *Logger) GetLoggingSettings() LoggingInfo {
	return logger.loggingInfo
}

//EnableLogging enables logging for the specified station
func (logger *Logger) SetLoggingEnabledStatus(stnName string, enabled bool) {
	oldStn := logger.loggingInfo.info[stnName]

	logger.loggingInfo.info[stnName] = LoggingMetaInfo{
		stationName:    oldStn.stationName,
		loggingEnabled: enabled,
		stationID:      oldStn.stationID,
		lastLogTime:    oldStn.lastLogTime}
}

//adds stations to the list that we keep track of the logging status
func (logger *Logger) addUploadingStationsToListForLogging() {
	fmt.Println("Logger: checking for new weather stations ", time.Now())

	var stns = logger.data.GetStations()
	fmt.Println(stns)
	fmt.Println(logger.loggingInfo.info)
	for index := range stns {
		//see if the station is already in our current list

		if _, ok := logger.loggingInfo.info[stns[index].Name]; ok {

		} else {
			fmt.Println(ok)
			station := LoggingMetaInfo{
				stationName:    stns[index].Name,
				loggingEnabled: false,
				stationID:      stns[index].StationID,
				lastLogTime:    time.Now()}

			fmt.Println("Logger: created new station", station)
			logger.loggingInfo.info[station.stationName] = station

		}
	}
	fmt.Println("Logger: New station meta info ", logger.loggingInfo.info)

	/*
		for stnIndex := range logger.data.GetStations() {

		}
	*/
}

func (logger *Logger) LogStations() {
	fmt.Println("Logger: Logging station data", time.Now())
	for key := range logger.loggingInfo.info {
		fmt.Println("Logger: Attempting to log station: ", key)
		// if the station came from the database then we do not want to log it
		//if val.stationID == -1 {
		fmt.Println("Logger: getting current sensor readings")
		stn := logger.data.GetCurrentSensorReadings(key)
		fmt.Println("Logger: Logging conditions")
		logger.data.LogConditions(&stn)
		//}

	}
}
