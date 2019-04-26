package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Josiah-B/Cyclone/Logger"
	"github.com/Josiah-B/Cyclone/MemoryCache"
	"github.com/Josiah-B/Cyclone/ProcessManager"

	"github.com/Josiah-B/Cyclone/Interfaces"

	"github.com/Josiah-B/Cyclone/ConfigAPI"
)

// structure for holding the settings: this is a temporary measure until we can get a standard in place
type Settings struct {
	configAPIport             string
	hostPort                  string
	dbPath                    string
	procManagerConfigFilePath string
}

var (
	httpMuxRouter HTTPMux
	dataStore     Interfaces.Storage

	//Configuration options
	settings = Settings{
		configAPIport: "8000",
		hostPort:      "8080",
		dbPath:        "./database.db",
		procManagerConfigFilePath: "./config/process manager.json"}

	//Proccess Management API
	confAPI ConfigAPI.HTTPMux

	//This handles pushing current sensor readings to the database
	logger *Logger.Logger

	//this monitors our station processes and restarts them if they crash
	processManager *ProcessManager.ProcessMgr
)

func main() {
	setupDatabase()

	//setup the process manager
	processManager = ProcessManager.NewProcessMgr(settings.procManagerConfigFilePath)
	//setup the configAPI
	confAPI.Create(processManager)
	go http.ListenAndServe(":"+settings.configAPIport, confAPI.Router)

	httpMuxRouter.Create(dataStore)
	logger = new(Logger.Logger)
	logger.Interval = 15 //Logging interval in Minutes
	logger.Initilize(dataStore)

	fmt.Println("Data Logger started")
	//fmt.Println("Closing Database...")
	//dataStore.Close()

	fmt.Println("Starting Webserver")

	// start the web server
	go http.ListenAndServe(":"+settings.hostPort, httpMuxRouter.router)

	fmt.Println("Webserver started")
	fmt.Println("Starting Logger")

	// wait for the user to enter the exit commend
	listenForExit()

	fmt.Println("Program Completed!")
}

func setupDatabase() {
	var temp = new(MemoryCache.Cache)
	temp.DataBasePath = settings.dbPath
	dataStore = temp

	dataStore.Initilize()

}

// listenForExit listens to the keyboard for the 'exit' command before returning.
func listenForExit() {
	fmt.Println("Type 'Exit' shutdown the program")
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")

	var continueToWait = true

	for continueToWait == true {
		sentence, err := buf.ReadBytes('\n')

		if err != nil {
			fmt.Println(err)
		} else {
			var trimmedString = string(sentence[:])
			if strings.Contains(strings.ToLower(trimmedString), "exit") {
				fmt.Println("Exiting program...")
				continueToWait = false
			} else {
				fmt.Println("Unknown Command: " + trimmedString)
				fmt.Print("> ")
			}
		}
	}

}
