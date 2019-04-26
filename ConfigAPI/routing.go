package ConfigAPI

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/Josiah-B/Cyclone/Interfaces"
	"github.com/Josiah-B/Cyclone/ProcessManager"
	"github.com/gorilla/mux"
)

type HTTPMux struct {
	Router  *mux.Router
	procMgr *ProcessManager.ProcessMgr
}

var (
	apiRoutes          []Interfaces.APIRoute
	configURL          = "/config/"
	configWebsiteFiles = "configsite/"
)

// Create the url mappings for the REST operations
func (httpMux *HTTPMux) Create(processManager *ProcessManager.ProcessMgr) {
	httpMux.procMgr = processManager
	httpMux.Router = mux.NewRouter()

	//setup the api mapping
	apiRoutes = []Interfaces.APIRoute{
		Interfaces.APIRoute{
			Route:         "/processes",
			HandlerMethod: httpMux.GetProcesses,
			HTTPMethod:    "GET",
			Description:   "Returns a all processes"},

		Interfaces.APIRoute{
			Route:         "/process/start/{processPath}",
			HandlerMethod: httpMux.StartProcess,
			HTTPMethod:    "GET",
			Description:   "Starts a process"},
		Interfaces.APIRoute{
			Route:         "/process/stop/{processID}",
			HandlerMethod: httpMux.StartProcess,
			HTTPMethod:    "GET",
			Description:   "Stops a process"},
	}

	//fmt.Println(apiRoutes)

	httpMux.Router.HandleFunc("/", httpMux.showAPIRouting)

	//Setup the webserver to handle the apiRoutes
	for _, v := range apiRoutes {
		httpMux.Router.HandleFunc(v.Route, v.HandlerMethod).Methods(v.HTTPMethod)
	}

	// webserver functions
	s := http.StripPrefix("/"+configWebsiteFiles, http.FileServer(http.Dir("./"+configWebsiteFiles)))
	httpMux.Router.PathPrefix("/" + configWebsiteFiles).Handler(s)
}

func (httpMux *HTTPMux) showAPIRouting(w http.ResponseWriter, r *http.Request) {
	// show the user some dosumentation on how the API is structured
	writeResponsePrettyfied(w, apiRoutes, "\t")
}

func writeResponsePrettyfied(webResponseWriter http.ResponseWriter, obj interface{}, indentationString string) {
	webResponseWriter.Header().Set("Content-Type", "application/json")
	objects, _ := json.MarshalIndent(obj, "", indentationString)
	webResponseWriter.Write(objects)
}

func (httpMux *HTTPMux) GetProcesses(webResponseWriter http.ResponseWriter, r *http.Request) {
	webResponseWriter.Header().Set("Content-Type", "application/json")
	webResponseWriter.Write(httpMux.procMgr.ListProcesses())
}

//StartProcess runs the specified command
func (httpMux *HTTPMux) StartProcess(webResponseWriter http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	httpMux.procMgr.CreateProc(vars["processPath"])
}

//StopProcess kills the specified process
func (httpMux *HTTPMux) StopProcess(webResponseWriter http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	procID, err := strconv.Atoi(vars["processID"])
	if err != nil {
		webResponseWriter.Write([]byte("That Process ID does not exist"))
		return
	}
	httpMux.procMgr.Stop(procID)
}
