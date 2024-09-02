//go:build linux

package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

func (self *Runner) disableTitlebar() {
	// empty (for conditional compilation)
}

const libre_office_path = "/usr/bin/libreoffice"      // Excel equivalent in Linux
const text_editor_path = "/usr/bin/gnome-text-editor" // Notepad equivalent in Linux
const libre_impress_path = "/usr/bin/ooimpress"       // PowerPoint equivalent in Linux

const kill_cmd = "kill -9"

type LinuxRunner struct {
	paths struct {
		kill       string
		excel      string
		notepad    string
		powerpoint string
	}
}

func (self *LinuxRunner) newRunner() (*LinuxRunner, error) {
	runner := &LinuxRunner{}

	runner.paths.kill = kill_cmd
	runner.paths.excel = libre_office_path
	runner.paths.notepad = text_editor_path
	runner.paths.powerpoint = libre_impress_path

	return runner, nil
}

func (self *LinuxRunner) killTasks(pids []string) error {
	for _, pid := range pids {
		cmd := self.paths.kill + " " + pid
		err := exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			return err
		}
		log.Printf("Killed process with PID %s", pid)
	}
	return nil
}

func (self *LinuxRunner) findPIDs(processName string) ([]string, error) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("pgrep -f %s", processName))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	pids := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(pids) == 0 || pids[0] == "" {
		return nil, fmt.Errorf("process not found")
	}

	return pids, nil
}

func (self *LinuxRunner) killLibreOffice() error {
	pids, err := self.findPIDs("libreoffice")
	if err != nil {
		return fmt.Errorf("error finding LibreOffice process: %v", err)
	}

	err = self.killTasks(pids)
	if err != nil {
		return fmt.Errorf("error killing LibreOffice process: %v", err)
	}

	log.Printf("LibreOffice processes killed successfully")
	return nil
}

func (self *LinuxRunner) runLibreOffice() error {
	// Create a new command to run LibreOffice
	cmd := exec.Command(self.paths.excel)

	// Start the command (don't use cmd.Wait, let it run in the background)
	err := cmd.Start()
	if err != nil {
		return err
	}

	log.Println("LibreOffice started successfully")
	return nil
}

func (self *LinuxRunner) runNotepad() error {
	cmd := exec.Command(self.paths.notepad)

	err := cmd.Start()
	if err != nil {
		return err
	}

	log.Println("Notepad started successfully")
	return nil
}

func (self *LinuxRunner) runPowerPoint() error {
	cmd := exec.Command(self.paths.powerpoint)

	err := cmd.Start()
	if err != nil {
		return err
	}

	log.Println("PowerPoint started successfully")
	return nil
}

func TestLinux() {
	runner, err := new(LinuxRunner).newRunner()
	if err != nil {
		log.Fatalf("Failed to initialize LinuxRunner: %v", err)
	}

	// Run LibreOffice
	err = runner.runLibreOffice()
	if err != nil {
		log.Fatalf("Failed to start LibreOffice: %v", err)
	}

	// Sleep for 15 seconds
	time.Sleep(10 * time.Second)

	// Kill LibreOffice
	err = runner.killLibreOffice()
	if err != nil {
		log.Fatalf("Failed to kill LibreOffice: %v", err)
	}
}
