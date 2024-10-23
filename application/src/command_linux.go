//go:build linux

package main

import (
	types "common"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
)

type Runner struct {
	send  chan<- types.Message
	paths struct {
		kill       string
		excel      string
		notepad    string
		powerpoint string
		word       string
	}
	state struct {
		running_typ types.AppType
		running_app *exec.Cmd
		file        string
		pid         int
	}
}

func NewRunner(send chan<- types.Message) (*Runner, error) {
	runner := &Runner{
		send: send,
	}

	runner.CheckApps()

	return runner, nil
}

func (self *Runner) CheckApps() {
	var err error
	self.paths.kill = "kill"
	self.paths.excel, err = exec.LookPath("libreoffice")
	if err != nil {
		self.send <- types.NewMessage(types.TExeNotFound{
			Name:   "libreoffice",
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}
	self.paths.notepad, err = exec.LookPath("gedit")
	if err != nil {
		self.send <- types.NewMessage(types.TExeNotFound{
			Name:   "gedit",
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}
	self.paths.powerpoint = self.paths.excel
	self.paths.word = self.paths.excel
}

func (self *Runner) SetupEnv() error {
	self.fullscreenForegroundWindow()
	return nil
}

func (self *Runner) RestoreEnv() error {
	return nil
}

func (self *Runner) OpenApp(typ types.AppType, file string) error {
	if self.isOpen() {
		return fmt.Errorf("an app is already running")
	}
	self.resetState()
	defer self.resetState()

	self.state.running_typ = typ
	self.state.file = file

	var cmd *exec.Cmd

	switch typ {
	case types.XLSX:
		cmd = exec.Command(self.paths.excel, "--calc", file)
	case types.TXT:
		cmd = exec.Command(self.paths.notepad, file)
	case types.PPTX:
		cmd = exec.Command(self.paths.powerpoint, "--impress", file)
	case types.DOCX:
		cmd = exec.Command(self.paths.word, "--writer", file)
	default:
		return fmt.Errorf("unsupported app type")
	}

	self.state.running_app = cmd

	err := cmd.Start()
	if err != nil {
		return err
	}

	self.state.pid = cmd.Process.Pid

	err = cmd.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return err
}

func (self *Runner) KillApp() error {
	if self.state.running_app == nil {
		return nil
	}

	err := self.state.running_app.Process.Kill()
	if err != nil {
		return err
	}

	self.resetState()
	return nil
}

func (self *Runner) IsAppOpen() bool {
	return self.isOpen()
}

func (self *Runner) FocusOrOpenApp(typ types.AppType, file string) error {
	if self.isOpen() && self.state.running_typ == typ {
		return self.FocusOpenApp()
	} else {
		return self.OpenApp(typ, file)
	}
}

func (self *Runner) FocusOpenApp() error {
	if self.state.pid == 0 {
		return fmt.Errorf("no app is currently open")
	}

	err := robotgo.ActivePid(self.state.pid)
	if err != nil {
		return fmt.Errorf("failed to focus app: %v", err)
	}

	self.fullscreenForegroundWindow()
	return nil
}

func (self *Runner) fullscreenForegroundWindow() {
	if self.state.pid == 0 {
		return
	}

	robotgo.MaxWindow(self.state.pid)
}

func (self *Runner) resetState() {
	self.state.file = ""
	self.state.pid = 0
	self.state.running_app = nil
	self.state.running_typ = 0
}

func (self *Runner) isOpen() bool {
	app := self.state.running_app
	if app != nil {
		if app.ProcessState == nil || !app.ProcessState.Exited() {
			return true
		}
	}
	return false
}

func (self *Runner) ListAllProcess() (map[uint32]string, error) {
	processes := make(map[uint32]string)

	fmt.Println("List All Processes called in linux")

	return processes, nil
}

// ListProcesses lists all running processes and returns their names and PIDs.
func ListProcesses() (map[int]string, error) {
	processes := make(map[int]string)

	// Read the /proc directory to list all processes
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			// Each process is represented by a numerical directory under /proc
			pid, err := strconv.Atoi(file.Name())
			if err != nil {
				continue // Skip non-numeric directories
			}

			// Read the command line of the process
			cmdlinePath := filepath.Join("/proc", file.Name(), "cmdline")
			cmdlineBytes, err := ioutil.ReadFile(cmdlinePath)
			if err != nil {
				continue
			}

			// cmdline is null-separated, split it to get the command name
			cmdline := strings.TrimSpace(string(cmdlineBytes))
			if cmdline != "" {
				processes[pid] = cmdline
			}
		}
	}

	return processes, nil
}

// KillProcess kills a process based on its PID.
func KillProcess(pid int) error {
	cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
	return cmd.Run()
}

func TestAppKills() {
	// List all running processes
	processes, err := ListProcesses()
	if err != nil {
		fmt.Printf("Error listing processes: %v\n", err)
		return
	}

	// Print running processes
	fmt.Println("Running Processes:")
	for pid, cmdline := range processes {
		fmt.Printf("PID: %d, Command: %s\n", pid, cmdline)
	}

	// List of apps to kill to prevent cheating
	appsToKill := []string{"notepad"}

	// Iterate over all processes and kill the ones that match appsToKill
	for pid, cmdline := range processes {
		for _, app := range appsToKill {
			if strings.Contains(cmdline, app) {
				fmt.Printf("Killing process %d (%s)\n", pid, cmdline)
				if err := KillProcess(pid); err != nil {
					fmt.Printf("Error killing process %d: %v\n", pid, err)
				}
			}
		}
	}
}
