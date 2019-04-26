package main

import (
	"fmt"
	"net/http"
	"strconv"

	"encoding/json"
	"io/ioutil"

	"github.com/Josiah-B/Cyclone/Interfaces"
	"github.com/gorilla/mux"
)

type HTTPMux struct {
	router *mux.Router
	db     Interfaces.Storage
}

var (
	apiRoutes []Interfaces.APIRoute
)

// Create the url mappings for the REST operations
func (httpMux *HTTPMux) Create(storage Interfaces.Storage) {
	fmt.Println("Setting up storage object...")
	httpMux.db = storage
	fmt.Println("Storage object setup")
	httpMux.router = mux.NewRouter()

	//setup the api mapping
	apiRoutes = []Interfaces.APIRoute{
		Interfaces.APIRoute{
			Route:         "/unitTypes/{unitTypeID}",
			HandlerMethod: httpMux.getUnitType,
			HTTPMethod:    "GET",
			Description:   "Returns a specific UnitType"},
		Interfaces.APIRoute{
			Route:         "/unitTypes",
			HandlerMethod: httpMux.getUnitTypes,
			HTTPMethod:    "GET",
			Description:   "Returns all UnitTypes"},
		Interfaces.APIRoute{
			Route:         "/unitTypes",
			HandlerMethod: httpMux.addUnitType,
			HTTPMethod:    "POST",
			Description:   "Adds new UnitType"},
		Interfaces.APIRoute{
			Route:         "/unitTypes",
			HandlerMethod: httpMux.modifyUnitType,
			HTTPMethod:    "PUT",
			Description:   "Modifies a current UnitType"},
		Interfaces.APIRoute{
			Route:         "/unitTypes",
			HandlerMethod: httpMux.deleteUnitType,
			HTTPMethod:    "DELETE",
			Description:   "Deletes a UnitType"},

		Interfaces.APIRoute{
			Route:         "/sensors/{sensorID}",
			HandlerMethod: httpMux.getSensor,
			HTTPMethod:    "GET",
			Description:   "Gets a specific Sensor"},
		Interfaces.APIRoute{
			Route:         "/sensors",
			HandlerMethod: httpMux.getSensors,
			HTTPMethod:    "GET",
			Description:   "Gets all the sensors"},
		Interfaces.APIRoute{
			Route:         "/sensors",
			HandlerMethod: httpMux.addSensor,
			HTTPMethod:    "POST",
			Description:   "Adds a new sensor"},
		Interfaces.APIRoute{
			Route:         "/sensors",
			HandlerMethod: httpMux.modifySensor,
			HTTPMethod:    "PUT",
			Description:   "Modifies a current sensor"},
		Interfaces.APIRoute{
			Route:         "/sensors",
			HandlerMethod: httpMux.deleteSensor,
			HTTPMethod:    "DELETE",
			Description:   "Deletes a sensor"},

		Interfaces.APIRoute{
			Route:         "/observedProperties/{propertyID}",
			HandlerMethod: httpMux.getObservedProperty,
			HTTPMethod:    "GET",
			Description:   "Gets a specific Observed Property"},
		Interfaces.APIRoute{
			Route:         "/getObservedProperties",
			HandlerMethod: httpMux.getObservedProperties,
			HTTPMethod:    "GET",
			Description:   "Gets all the observed properties"},
		Interfaces.APIRoute{
			Route:         "/addObservedProperty",
			HandlerMethod: httpMux.addObservedProperty,
			HTTPMethod:    "POST",
			Description:   "Adds a new Observed Property"},
		Interfaces.APIRoute{
			Route:         "/observedProperties",
			HandlerMethod: httpMux.modifyObservedProperty,
			HTTPMethod:    "PUT",
			Description:   "Modifies a current Observed Property"},
		Interfaces.APIRoute{
			Route:         "/observedProperties",
			HandlerMethod: httpMux.deleteObservedProperty,
			HTTPMethod:    "DELETE",
			Description:   "Deletes an Observed Property"},

		Interfaces.APIRoute{
			Route:         "/dataStreams/{streamID}",
			HandlerMethod: httpMux.getDataStream,
			HTTPMethod:    "GET",
			Description:   "Gets a specific Data Stream"},
		Interfaces.APIRoute{
			Route:         "/dataStreams",
			HandlerMethod: httpMux.getDataStreams,
			HTTPMethod:    "GET",
			Description:   "Gets all Data Streams"},
		Interfaces.APIRoute{
			Route:         "/dataStreams",
			HandlerMethod: httpMux.addDataStream,
			HTTPMethod:    "POST",
			Description:   "Adds a new data stream"},
		Interfaces.APIRoute{
			Route:         "/dataStreams",
			HandlerMethod: httpMux.modifyDataStream,
			HTTPMethod:    "PUT",
			Description:   "Modifies a current data stream"},
		Interfaces.APIRoute{
			Route:         "/dataStreams",
			HandlerMethod: httpMux.deleteDataStream,
			HTTPMethod:    "DELETE",
			Description:   "Deletes a current data stream"},

		Interfaces.APIRoute{
			Route:         "/observation/{observationID}",
			HandlerMethod: httpMux.getObservation,
			HTTPMethod:    "GET",
			Description:   "Gets a specific Observation"},
		Interfaces.APIRoute{
			Route:         "/observations/{sensorID}",
			HandlerMethod: httpMux.getObservations,
			HTTPMethod:    "GET",
			Description:   "Gets all of the logged observations for the specified sensor"},
		Interfaces.APIRoute{
			Route:         "/observation",
			HandlerMethod: httpMux.addObservation,
			HTTPMethod:    "POST",
			Description:   "Adds a new Observation"},
		Interfaces.APIRoute{
			Route:         "/observation",
			HandlerMethod: httpMux.modifyObservation,
			HTTPMethod:    "PUT",
			Description:   "Modifies a current Observation"},
		Interfaces.APIRoute{
			Route:         "/observation/{observationID}",
			HandlerMethod: httpMux.deleteObservation,
			HTTPMethod:    "DELETE",
			Description:   "Deletes a current Observation"},

		Interfaces.APIRoute{
			Route:         "/stations/{stationID}",
			HandlerMethod: httpMux.getStation,
			HTTPMethod:    "GET",
			Description:   "Gets a specific station; by station ID number"},
		Interfaces.APIRoute{
			Route:         "/stations",
			HandlerMethod: httpMux.getStations,
			HTTPMethod:    "GET",
			Description:   "Lists all stations in the database"},
		Interfaces.APIRoute{
			Route:         "/stations",
			HandlerMethod: httpMux.addStation,
			HTTPMethod:    "POST",
			Description:   "Adds a new station"},
		Interfaces.APIRoute{
			Route:         "/stations/logConditions",
			HandlerMethod: httpMux.logConditions,
			HTTPMethod:    "POST",
			Description:   "Logs conditions to the Database"},
		Interfaces.APIRoute{
			Route:         "/stations/setCurrentConditions",
			HandlerMethod: httpMux.setCurrentConditions,
			HTTPMethod:    "POST",
			Description:   "Sets the current conditions for a station"},
		Interfaces.APIRoute{
			Route:         "/stations",
			HandlerMethod: httpMux.modifyStation,
			HTTPMethod:    "PUT",
			Description:   "Modifies a current station"},
		Interfaces.APIRoute{
			Route:         "/stations",
			HandlerMethod: httpMux.deleteStation,
			HTTPMethod:    "DELETE",
			Description:   "Deletes a station"},

		Interfaces.APIRoute{
			Route:         "/current/{stationName}",
			HandlerMethod: httpMux.getCurrentConditions,
			HTTPMethod:    "GET",
			Description:   "Gets the current conditions for a specific station"},
	}

	//fmt.Println(apiRoutes)

	httpMux.router.HandleFunc("/", httpMux.showHomePage)

	//Setup the webserver to handle the apiRoutes
	for _, v := range apiRoutes {
		httpMux.router.HandleFunc(v.Route, v.HandlerMethod).Methods(v.HTTPMethod)
	}
	/* old routing configuration setup; DELETE this once the new API declaration has been tested.

	httpMux.router.HandleFunc("/unitTypes/{unitTypeID}", httpMux.getUnitType).Methods("GET") // gets a specific unitType
	httpMux.router.HandleFunc("/unitTypes", httpMux.getUnitTypes).Methods("GET")             // returns all unitTypes
	httpMux.router.HandleFunc("/unitTypes", httpMux.addUnitType).Methods("POST")
	httpMux.router.HandleFunc("/unitTypes", httpMux.modifyUnitType).Methods("PUT")
	httpMux.router.HandleFunc("/unitTypes", httpMux.deleteUnitType).Methods("DELETE")

	httpMux.router.HandleFunc("/sensors/{sensorID}", httpMux.getSensor).Methods("GET")
	httpMux.router.HandleFunc("/sensors", httpMux.getSensors).Methods("GET")
	httpMux.router.HandleFunc("/sensors", httpMux.addSensor).Methods("POST")
	httpMux.router.HandleFunc("/sensors", httpMux.modifySensor).Methods("PUT")
	httpMux.router.HandleFunc("/sensors", httpMux.deleteSensor).Methods("DELETE")

	httpMux.router.HandleFunc("/observedProperties/{propertyID}", httpMux.getObservedProperty).Methods("GET")
	httpMux.router.HandleFunc("/observedProperties", httpMux.getObservedProperties).Methods("GET")
	httpMux.router.HandleFunc("/observedProperties", httpMux.addObservedProperty).Methods("POST")
	httpMux.router.HandleFunc("/observedProperties", httpMux.modifyObservedProperty).Methods("PUT")
	httpMux.router.HandleFunc("/observedProperties", httpMux.deleteObservedProperty).Methods("DELETE")

	httpMux.router.HandleFunc("/dataStreams/{streamID}", httpMux.getDataStream).Methods("GET")
	httpMux.router.HandleFunc("/dataStreams", httpMux.getDataStreams).Methods("GET")
	httpMux.router.HandleFunc("/dataStreams", httpMux.addDataStream).Methods("POST")
	httpMux.router.HandleFunc("/dataStreams", httpMux.modifyDataStream).Methods("PUT")
	httpMux.router.HandleFunc("/dataStreams", httpMux.deleteDataStream).Methods("DELETE")

	httpMux.router.HandleFunc("/observations/{observationID}", httpMux.getObservation).Methods("GET")
	httpMux.router.HandleFunc("/observations", httpMux.addObservation).Methods("POST")
	httpMux.router.HandleFunc("/observations", httpMux.modifyObservation).Methods("PUT")
	httpMux.router.HandleFunc("/observations", httpMux.deleteObservation).Methods("DELETE")

	httpMux.router.HandleFunc("/stations/{stationID}", httpMux.getStation).Methods("GET")
	httpMux.router.HandleFunc("/stations", httpMux.getStations).Methods("GET")
	httpMux.router.HandleFunc("/stations", httpMux.addStation).Methods("POST")
	httpMux.router.HandleFunc("/stations/logConditions", httpMux.logConditions).Methods("POST")
	httpMux.router.HandleFunc("/stations/setCurrentConditions", httpMux.setCurrentConditions).Methods("POST")
	httpMux.router.HandleFunc("/stations", httpMux.modifyStation).Methods("PUT")
	httpMux.router.HandleFunc("/stations", httpMux.deleteStation).Methods("DELETE")

	httpMux.router.HandleFunc("/current/{stationName}", httpMux.getCurrentConditions).Methods("GET")
	*/

	// webserver functions
	s := http.StripPrefix("/website/", http.FileServer(http.Dir("./website/")))
	httpMux.router.PathPrefix("/website/").Handler(s)
}

// unmarshalToObject parses the http.request's body into the 'obj' variable
func (httpMux *HTTPMux) unmarshalToObject(webRequest *http.Request, obj interface{}) {
	body, err := ioutil.ReadAll(webRequest.Body)
	defer webRequest.Body.Close()

	if err != nil {
		fmt.Println(err)
	}

	errr := json.Unmarshal(body, &obj)
	if errr != nil {
		fmt.Println(errr)
	}

}

func (httpMux *HTTPMux) showHomePage(w http.ResponseWriter, r *http.Request) {
	// show the user some dosumentation on how the API is structured
	writeResponsePrettyfied(w, apiRoutes, "\t")

}

func writeResponsePrettyfied(webResponseWriter http.ResponseWriter, obj interface{}, indentationString string) {
	webResponseWriter.Header().Set("Content-Type", "application/json")
	objects, _ := json.MarshalIndent(obj, "", indentationString)
	webResponseWriter.Write(objects)
}

// these functions route the http requests to Database actions, then marshall the result to JSON and send it back to the client

func (httpMux *HTTPMux) getUnitType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	unitType := httpMux.db.GetUnitType(vars["unitTypeID"])
	writeResponsePrettyfied(w, unitType, "\t")
}

func (httpMux *HTTPMux) getUnitTypes(w http.ResponseWriter, r *http.Request) {

	unitTypes := httpMux.db.GetUnitTypes()
	writeResponsePrettyfied(w, unitTypes, "\t")
}

func (httpMux *HTTPMux) modifyUnitType(w http.ResponseWriter, r *http.Request) {
	var unit Interfaces.UnitType
	httpMux.unmarshalToObject(r, &unit)
	httpMux.db.AddOrUpdateUnitType(&unit)
}
func (httpMux *HTTPMux) deleteUnitType(w http.ResponseWriter, r *http.Request) {}
func (httpMux *HTTPMux) addUnitType(w http.ResponseWriter, r *http.Request) {
	var unit Interfaces.UnitType
	httpMux.unmarshalToObject(r, &unit)
	httpMux.db.AddOrUpdateUnitType(&unit)
}

func (httpMux *HTTPMux) getSensor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sensor := httpMux.db.GetSensor(vars["sensorID"])
	writeResponsePrettyfied(w, sensor, "\t")
}

func (httpMux *HTTPMux) getSensors(w http.ResponseWriter, r *http.Request) {
	sensors := httpMux.db.GetSensors()
	writeResponsePrettyfied(w, sensors, "\t")
}

func (httpMux *HTTPMux) modifySensor(w http.ResponseWriter, r *http.Request) {
	var unit Interfaces.Sensor
	httpMux.unmarshalToObject(r, &unit)
	httpMux.db.AddOrUpdateSensor(&unit)
}
func (httpMux *HTTPMux) deleteSensor(w http.ResponseWriter, r *http.Request) {}
func (httpMux *HTTPMux) addSensor(w http.ResponseWriter, r *http.Request) {
	var unit Interfaces.Sensor
	httpMux.unmarshalToObject(r, &unit)
	httpMux.db.AddOrUpdateSensor(&unit)
}

func (httpMux *HTTPMux) getObservedProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	observations := httpMux.db.GetObservedProperty(vars["propertyID"])
	writeResponsePrettyfied(w, observations, "\t")
}

func (httpMux *HTTPMux) getObservedProperties(w http.ResponseWriter, r *http.Request) {
	observations := httpMux.db.GetObservedProperties()
	writeResponsePrettyfied(w, observations, "\t")
}

func (httpMux *HTTPMux) modifyObservedProperty(w http.ResponseWriter, r *http.Request) {
	var unit Interfaces.ObservedProperty
	httpMux.unmarshalToObject(r, &unit)
	httpMux.db.AddOrUpdateObservedProperty(&unit)
}
func (httpMux *HTTPMux) deleteObservedProperty(w http.ResponseWriter, r *http.Request) {}
func (httpMux *HTTPMux) addObservedProperty(w http.ResponseWriter, r *http.Request) {
	var unit Interfaces.ObservedProperty
	httpMux.unmarshalToObject(r, &unit)
	httpMux.db.AddOrUpdateObservedProperty(&unit)
}

func (httpMux *HTTPMux) getDataStream(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	dataStream := httpMux.db.GetDataStream(vars["streamID"])
	writeResponsePrettyfied(w, dataStream, "\t")
}

func (httpMux *HTTPMux) getDataStreams(w http.ResponseWriter, r *http.Request) {
	dataStreams := httpMux.db.GetDataStreams()
	writeResponsePrettyfied(w, dataStreams, "\t")
}

func (httpMux *HTTPMux) modifyDataStream(w http.ResponseWriter, r *http.Request) {
	var unit Interfaces.DataStream
	httpMux.unmarshalToObject(r, &unit)
	httpMux.db.UpdateDataStream(&unit)
}
func (httpMux *HTTPMux) deleteDataStream(w http.ResponseWriter, r *http.Request) {}
func (httpMux *HTTPMux) addDataStream(w http.ResponseWriter, r *http.Request) {
	var unit Interfaces.DataStream
	httpMux.unmarshalToObject(r, &unit)
	httpMux.db.AddDataStream(&unit)
}

func (httpMux *HTTPMux) getObservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	observations := httpMux.db.GetObservation(vars["observationID"])
	writeResponsePrettyfied(w, observations, "\t")
}

func (httpMux *HTTPMux) getObservations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	observationParams := new(Interfaces.ObservationParameters)

	intval, _ := strconv.ParseInt(vars["sensorID"], 10, 0)
	observationParams.SensorID = int(intval)
	observations := httpMux.db.GetObservations(*observationParams)
	writeResponsePrettyfied(w, observations, "\t")
}

func (httpMux *HTTPMux) modifyObservation(w http.ResponseWriter, r *http.Request) {}
func (httpMux *HTTPMux) deleteObservation(w http.ResponseWriter, r *http.Request) {}
func (httpMux *HTTPMux) addObservation(w http.ResponseWriter, r *http.Request) {
	var unit Interfaces.Observation
	httpMux.unmarshalToObject(r, &unit)
	httpMux.db.AddObservation(&unit)
}

func (httpMux *HTTPMux) getStation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	wxstation := httpMux.db.GetStation(vars["stationID"])
	writeResponsePrettyfied(w, wxstation, "\t")
}

func (httpMux *HTTPMux) getCurrentConditions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	station := httpMux.db.GetCurrentSensorReadings(vars["stationName"])
	fmt.Println("Station: " + station.StationName)
	writeResponsePrettyfied(w, station, "\t")
}

func (httpMux *HTTPMux) getStations(w http.ResponseWriter, r *http.Request) {
	wxstations := httpMux.db.GetStations()
	writeResponsePrettyfied(w, wxstations, "\t")
}

func (httpMux *HTTPMux) modifyStation(w http.ResponseWriter, r *http.Request) {
	var unit Interfaces.Station
	httpMux.unmarshalToObject(r, &unit)
	httpMux.db.AddOrUpdateStation(&unit)
}
func (httpMux *HTTPMux) deleteStation(w http.ResponseWriter, r *http.Request) {}
func (httpMux *HTTPMux) addStation(w http.ResponseWriter, r *http.Request) {
	var unit Interfaces.Station
	httpMux.unmarshalToObject(r, &unit)
	httpMux.db.AddOrUpdateStation(&unit)
}

func (httpMux *HTTPMux) setCurrentConditions(w http.ResponseWriter, r *http.Request) {
	// convert the json to an object
	currentConditions := Interfaces.StationUploadTemplate{}
	httpMux.unmarshalToObject(r, &currentConditions)

	fmt.Println("Setting current weather conditions for station: " + currentConditions.StationName)

	//store the current station in our in-memory listing of current conditions.
	httpMux.db.SetCurrentSensorReadings(currentConditions)
}

//logConditions logs the sensor info to the database
func (httpMux *HTTPMux) logConditions(w http.ResponseWriter, r *http.Request) {
	// convert the json to an object
	var currentConditions Interfaces.StationUploadTemplate
	fmt.Println(r)
	httpMux.unmarshalToObject(r, &currentConditions)

	httpMux.db.LogConditions(&currentConditions)
}
