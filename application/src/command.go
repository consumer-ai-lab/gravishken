package main

import (
	"fmt"
	assets "gravtest"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/go-vgo/robotgo"
)

const kill = "TASKKILL.exe"
const explorer = "explorer.exe"
const cmd = "cmd.exe"
const word = "WINWORD.exe"
const excel = "EXCEL.exe"
const powerpoint = "POWERPNT.exe"
const notepad = "NOTEPAD.exe"

const windows = "windows"

const template_docx = "template.docx"
const template_xlsx = "template.xlsx"
const template_pptx = "template.pptx"
const template_txt = "template.txt"
const tmp_docx = "tmp_*.docx"
const tmp_xlsx = "tmp_*.xlsx"
const tmp_pptx = "tmp_*.pptx"
const tmp_txt = "tmp_*.txt"

type Runner struct {
	paths struct {
		kill       string
		explorer   string
		cmd        string
		word       string
		excel      string
		notepad    string
		powerpoint string
	}
}

func newRunner() (*Runner, error) {
	runner := &Runner{}

	log.Println(runtime.GOOS)
	if runtime.GOOS != "windows" {
		return runner, nil
	}

	var err error
	runner.paths.cmd, err = exec.LookPath(cmd)
	log.Println(runner.paths.cmd)
	if err != nil {
		return nil, err
	}
	runner.paths.kill, err = exec.LookPath(kill)
	log.Println(runner.paths.kill)
	if err != nil {
		return nil, err
	}
	runner.paths.explorer, err = exec.LookPath(explorer)
	log.Println(runner.paths.explorer)
	if err != nil {
		return nil, err
	}
	runner.paths.notepad, err = exec.LookPath(notepad)
	log.Println(runner.paths.notepad)
	if err != nil {
		return nil, err
	}
	// runner.paths.word, err = exec.LookPath(word)
	// log.Println(runner.paths.word)
	// if err != nil {
	// 	return nil, err
	// }
	// runner.paths.excel, err = exec.LookPath(excel)
	// log.Println(runner.paths.excel)
	// if err != nil {
	// 	return nil, err
	// }
	// runner.paths.powerpoint, err = exec.LookPath(powerpoint)
	// log.Println(runner.paths.powerpoint)
	// if err != nil {
	// 	return nil, err
	// }

	return runner, nil
}

func (self *Runner) killExplorer() error {
	return self.kill(explorer)
}

func (self *Runner) startExplorer() {
	if runtime.GOOS != windows {
		return
	}

	// OOF: running explorer.exe always seems to return 1 :/
	command := exec.Command(self.paths.cmd, "/C", "start", self.paths.explorer)
	err := command.Run()
	if err != nil {
		log.Println(err)
	}
}

func (self *Runner) kill(name string) error {
	if runtime.GOOS != "windows" {
		return nil
	}

	// command := exec.Command(self.paths.cmd, "/C", self.paths.kill, "/F", "/IM", name)
	command := exec.Command(self.paths.kill, "/F", "/IM", name)
	out, err := command.CombinedOutput()
	log.Printf("%s\n", string(out))
	if err != nil {
		log.Println(err)
	}
	return err
}

func (self *Runner) open(exe string, file string) error {
	if runtime.GOOS != "windows" {
		return nil
	}

	// cmd := exec.Command(self.paths.explorer, file)
	// cmd := exec.Command(self.paths.cmd, "/C", "start", file)
	cmd := exec.Command(exe, file)
	out, err := cmd.CombinedOutput()
	log.Printf("%s\n", string(out))
	log.Println(err)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (self *Runner) newTemplate(name string, dest string) error {
	// NOTE: non os specific path separaters
	path := fmt.Sprintf("templates/%s", name)
	contents, err := assets.Templates.ReadFile(path)
	if err != nil {
		return err
	}
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(contents)
	return err
}

func (self *Runner) fullscreenForegroundWindow() {
	pid := robotgo.GetPid()
	robotgo.MaxWindow(pid)
}

func (self *Runner) run(name string, args ...string) error {
	if runtime.GOOS != "windows" {
		return nil
	}

	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	return err
}

func (self *Runner) maybe(err error) {
	if err != nil {
		log.Println(err)
		if self != nil {
			self.startExplorer()
		}
		panic(err)
	}
}

func test() {
	pids, _ := robotgo.Pids()
	for _, p := range pids {
		log.Println(p)
	}

	runner, err := newRunner()
	runner.maybe(err)

	err = runner.killExplorer()
	runner.maybe(err)

	file, err := os.CreateTemp("", tmp_txt)
	runner.maybe(err)
	file.Close()

	dest := file.Name()
	log.Println(dest)
	err = runner.newTemplate(template_txt, dest)
	runner.maybe(err)

	time.Sleep(2000 * time.Millisecond)

	disableTitlebar()
	runner.fullscreenForegroundWindow()
	gadsgadd()

	go (func() {
		_ = runner.open(runner.paths.notepad, dest)
	})()

	time.Sleep(3000 * time.Millisecond)

	err = runner.kill(notepad)
	runner.maybe(err)

	err = os.Remove(dest)
	runner.maybe(err)

	runner.startExplorer()
}
func asyncTest() {
	path, err := exec.LookPath("")
	command := exec.Command(path)
	command.Start()
	command.Process.Kill()
	err = command.Wait()
	if err != nil {

	}
}
