//go:build windows

package main

import (
	types "common"
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	getForegroundWindow = user32.NewProc("GetForegroundWindow")
	setWindowLong       = user32.NewProc("SetWindowLongW")
	getWindowLong       = user32.NewProc("GetWindowLongW")
	showWindow          = user32.NewProc("ShowWindow")
	setWindowPos        = user32.NewProc("SetWindowPos")
	enumWindows         = user32.NewProc("EnumWindows")
	getWindowText       = user32.NewProc("GetWindowTextW")
	getWindowTextLength = user32.NewProc("GetWindowTextLengthW")
)

const (
	GWL_STYLE = 0xFFFFFFFFFFFFFFF0
)

const kill = "TASKKILL.exe"
const explorer = "explorer.exe"
const cmd = "cmd.exe"
const word = "WINWORD.exe"
const excel = "EXCEL.exe"
const powerpoint = "POWERPNT.exe"
const notepad = "NOTEPAD.exe"

const windows = "windows"

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
	running_app     *exec.Cmd
	explorer_killed bool
}

func NewRunner(send chan<- types.Message) (*Runner, error) {
	runner := &Runner{}

	var err error
	runner.paths.cmd, err = exec.LookPath(cmd)
	if err != nil {
		return nil, err
	}
	runner.paths.kill, err = exec.LookPath(kill)
	if err != nil {
		return nil, err
	}
	runner.paths.explorer, err = exec.LookPath(explorer)
	if err != nil {
		return nil, err
	}
	runner.paths.notepad, err = exec.LookPath(notepad)
	if err != nil {
		send <- types.NewMessage(types.TExeNotFound{
			Name:   notepad,
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}
	runner.paths.word, err = exec.LookPath(word)
	if err != nil {
		send <- types.NewMessage(types.TExeNotFound{
			Name:   word,
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}
	runner.paths.excel, err = exec.LookPath(excel)
	if err != nil {
		send <- types.NewMessage(types.TExeNotFound{
			Name:   excel,
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}
	runner.paths.powerpoint, err = exec.LookPath(powerpoint)
	if err != nil {
		send <- types.NewMessage(types.TExeNotFound{
			Name:   powerpoint,
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}

	return runner, err
}

func (self *Runner) SetupEnv() error {
	err := self.killExplorer()
	if err != nil {
		return err
	}
	self.explorer_killed = true
	self.disableTitlebar()
	self.fullscreenForegroundWindow()
	return err
}

func (self *Runner) RestoreEnv() error {
	if self.explorer_killed {
		self.startExplorer()
	}
	self.explorer_killed = false
	return nil
}

// waits until app is finished runninig
func (self *Runner) OpenApp(typ AppType, file string) error {
	if self.running_app != nil {
		if self.running_app.Process != nil {
			return fmt.Errorf("an app is already running")
		} else {
			self.running_app = nil
		}
	}

	switch typ {
	case TXT:
		return self.open(self.paths.notepad, file)
	case DOCX:
		return self.open(self.paths.word, file)
	case PPTX:
		return self.open(self.paths.powerpoint, file)
	case XLSX:
		return self.open(self.paths.excel, file)
	default:
		return fmt.Errorf("invalid app type: %d", typ)
	}
}

func (self *Runner) KillApp() error {
	if self.running_app != nil {
		return nil
	}

	err := self.running_app.Process.Kill()
	if err != nil {
		self.running_app = nil
	}
	return err
}

func (self *Runner) killExplorer() error {
	return self.kill(explorer)
}

func (self *Runner) startExplorer() {
	// OOF: running explorer.exe always seems to return 1 :/
	command := exec.Command(self.paths.cmd, "/C", "start", self.paths.explorer)
	err := command.Run()
	if err != nil {
		log.Println(err)
	}
}

func (self *Runner) kill(name string) error {
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
	// cmd := exec.Command(self.paths.explorer, file)
	// cmd := exec.Command(self.paths.cmd, "/C", "start", file)
	cmd := exec.Command(exe, file)
	self.running_app = cmd
	out, err := cmd.CombinedOutput()
	log.Printf("%s\n", string(out))
	log.Println(err)
	if err != nil {
		log.Println(err)
	}
	return err
}

// func (self *Runner) fullscreenForegroundWindow() {

// 	const (
// 		SW_MAXIMIZE = 3
// 		SW_SHOW     = 5
// 	)

// 	hwnd, _, _ := getForegroundWindow.Call()
// 	showWindow.Call(hwnd, SW_MAXIMIZE)
// 	setWindowPos.Call(hwnd, 0, 0, 0, 1920, 1080, 0)
// }

func enumWindowsCallbackTest() {
	enumWindows.Call(syscall.NewCallback(EnumWindowsProc), 0)
}

type EnumWindowsCallback func(hwnd syscall.Handle, lParam uintptr) uintptr

func EnumWindowsProc(hwnd syscall.Handle, lParam uintptr) uintptr {
	length, _, _ := getWindowTextLength.Call(uintptr(hwnd))
	if length == 0 {
		return 1 // Continue enumeration
	}

	buf := make([]uint16, length+1)
	_, _, _ = getWindowText.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(length+1))

	title := syscall.UTF16ToString(buf)

	winhand := win.HWND(hwnd)
	if win.IsWindowVisible(winhand) {
		log.Printf("HWND: %v, Title: %s\n", hwnd, title)
	}

	return 1 // Continue enumeration
}

func (self *Runner) disableTitlebar() {
	hwnd, _, _ := getForegroundWindow.Call()

	style, _, _ := getWindowLong.Call(hwnd, GWL_STYLE)

	newStyle := style &^ (win.WS_CAPTION | win.WS_BORDER | win.WS_DLGFRAME)
	newStyle = 0x000000
	newStyle |= win.WS_POPUP | win.WS_VISIBLE

	_, _, err := setWindowLong.Call(hwnd, GWL_STYLE, newStyle)
	if err != nil && err.Error() != "The operation completed successfully." {
		log.Println("Error setting window style:", err)
	} else {
		log.Println("Title bar and borders removed successfully.")
	}

	// - [maybe keep setting this in a loop to keep the window in bg](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowpos)
	//   - HWND_BOTTOM
	// - MAYBE: keep looping through unknown windows and keep hiding them?
	syscall.SyscallN(uintptr(user32.NewProc("SetWindowPos").Addr()), 5,
		hwnd, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)
}
