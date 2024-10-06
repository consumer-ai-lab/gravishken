//go:build linux

package main

import (
	types "common"
	"fmt"
	"log"
	"os/exec"

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

	var err error
	runner.paths.kill = "kill"
	runner.paths.excel, err = exec.LookPath("libreoffice")
	if err != nil {
		send <- types.NewMessage(types.TExeNotFound{
			Name:   "libreoffice",
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}
	runner.paths.notepad, err = exec.LookPath("gedit")
	if err != nil {
		send <- types.NewMessage(types.TExeNotFound{
			Name:   "gedit",
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}
	runner.paths.powerpoint = runner.paths.excel
	runner.paths.word = runner.paths.excel

	return runner, nil
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