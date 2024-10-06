package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

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

func main() {
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
