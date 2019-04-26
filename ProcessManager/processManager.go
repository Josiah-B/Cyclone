package ProcessManager

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

//ProcessMgr allows management of external processes as used by the extensions
type ProcessMgr struct {
	processes      map[int]*process
	configFilePath string
	configuration  Config
}

var newProcID = 0

type process struct {
	pathToExec    string
	command       *exec.Cmd
	Status        string
	errors        error
	Stdout        io.ReadCloser
	LastHeartbeat time.Time
}

type Config struct {
	StationConfigFolder string
	StationExecFolder   string
}

//stationConfig is a structure that we can load the station config into just long enough to find the executable name for the station type
type stationConfig struct {
	ExecName string
}

func (prcMgr *ProcessMgr) loadConfigurationFile(file string, config interface{}) error {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println("Open File error : ", err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	err2 := jsonParser.Decode(&config)
	fmt.Println("Decode error: ", err2)
	return err
}

func (prcMgr *ProcessMgr) LoadStatationConfigs() {
	// walk all files in directory
	files, err := ioutil.ReadDir(prcMgr.configuration.StationConfigFolder)
	if err != nil {
		log.Fatal(err)
	}

	// each file is a config file
	for _, f := range files {

		baseName := f.Name()
		var cExec stationConfig
		var fullConfigFilePath = prcMgr.configuration.StationConfigFolder + "/" + baseName
		prcMgr.loadConfigurationFile(fullConfigFilePath, &cExec)
		fmt.Println(fullConfigFilePath, ":", cExec.ExecName)
		var fullExecString = prcMgr.configuration.StationExecFolder + "/" + cExec.ExecName
		fmt.Println("Full executable string : ", fullExecString)
		prcMgr.CreateProc(fullExecString, "-config="+fullConfigFilePath)
	}
}

//NewProcessMgr creates a new Proccess Manager
func NewProcessMgr(configFilePath string) *ProcessMgr {
	var prcMgr ProcessMgr

	prcMgr.processes = make(map[int]*process)
	//load the configuration file
	err := prcMgr.loadConfigurationFile(configFilePath, &prcMgr.configuration)
	fmt.Println(prcMgr.configuration)
	prcMgr.LoadStatationConfigs()
	if err != nil {
		fmt.Println(err)
	}
	//start monitoring the processes
	go prcMgr.monitorProcesses()
	return &prcMgr
}

//monitorProcesses periodically checks the processes, restarts any that are hung or crashed
func (prcMgr *ProcessMgr) monitorProcesses() {
	for true {
		fmt.Println("Checking processes...")
		for procID, value := range prcMgr.processes {
			durationSinceHeartbeat := time.Now().Sub(value.LastHeartbeat)
			if value.Status == "Launching" {
				prcMgr.StartProc(procID)
				go value.listen()
				// restart the process if it has crashed or it has not reported a heartbeat in awhile (e.g. the process is hung)
			} else if (durationSinceHeartbeat.Minutes() > 3) || value.Status == "Stopped" {
				prcMgr.recreateProc(procID)          // recreate the proccess
				prcMgr.StartProc(procID)             // start the new process
				go prcMgr.processes[procID].listen() // finally, start listening for input from the new process
			}
		}
		time.Sleep(time.Second * 60) // wait a little before checking the processes again
	}
}

//listen listens to the process and updates the lastHeartbeat variable whenever something is recieved from the process.
func (prc *process) listen() {
	scanner := bufio.NewScanner(prc.Stdout)
	prc.Status = "Starting"
	for scanner.Scan() {
		if scanner.Err() != nil {
			fmt.Println(scanner.Err())
		}

		prc.Status = "Running"         //The process is running since we are getting output from it
		prc.LastHeartbeat = time.Now() // anything we get from the process we will say is a valid hearbeat
	}
	prc.Status = "Stopped"
}

func (prcMgr *ProcessMgr) recreateProc(procID int) {
	prcMgr.processes[procID].command.Process.Kill() // make sure the proccess is not running

	path := prcMgr.processes[procID].pathToExec

	var prog = process{
		pathToExec: path,
		command:    exec.Command(path),
		Status:     "Launching",
	}

	prcMgr.processes[procID] = &prog
}

func (prcMgr *ProcessMgr) CreateProc(path string, args ...string) {
	// create and start the new process
	var prog = process{
		pathToExec: path,
		command:    exec.Command(path, args...),
		Status:     "Launching",
	}

	// if the process started then add it to the list
	prcMgr.processes[newProcID] = &prog
	newProcID++
}

// Start launches a new process
func (prcMgr *ProcessMgr) StartProc(procID int) {
	proc := prcMgr.processes[procID]
	// open the output pipe
	proc.Stdout, proc.errors = proc.command.StdoutPipe()
	if proc.errors != nil {
		proc.Status = "Unknown"
		fmt.Println(proc.errors)
	}

	//start the process
	proc.errors = proc.command.Start()
	if proc.errors != nil {
		log.Println("Failed to start process: ", proc.errors)
		return
	}
}

// Stop ends the specified process
func (prcMgr *ProcessMgr) Stop(ID int) {
	log.Println("Killing process : ", ID)
	prcMgr.processes[ID].errors = prcMgr.processes[ID].command.Process.Kill()
	if prcMgr.processes[ID].errors != nil {
		log.Println("Failed to kill process: ", ID, " ; ", prcMgr.processes[ID].errors)
	} else {
		prcMgr.processes[ID].Status = "Stopped"
		delete(prcMgr.processes, ID) //remove the process from our list of proccess that should be running
	}

}

//Restart stops then starts a specific process
func (prcMgr *ProcessMgr) Restart() {
	//prcMgr.Stop()
	//prcMgr.Start()
}

//ListProcesses lists all the processes
func (prcMgr *ProcessMgr) ListProcesses() []byte {
	for PID, value := range prcMgr.processes {
		fmt.Println(PID, "Status: ", value.Status)
	}
	payload, _ := json.MarshalIndent(prcMgr.processes, "", "\t")

	return payload
}
