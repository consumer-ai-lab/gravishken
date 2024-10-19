//go:build windows

package main

import (
	types "common"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"golang.org/x/sys/windows/registry"
)

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	getForegroundWindow          = user32.NewProc("GetForegroundWindow")
	setWindowLong                = user32.NewProc("SetWindowLongW")
	getWindowLong                = user32.NewProc("GetWindowLongW")
	showWindow                   = user32.NewProc("ShowWindow")
	setWindowPos                 = user32.NewProc("SetWindowPos")
	enumWindows                  = user32.NewProc("EnumWindows")
	getWindowText                = user32.NewProc("GetWindowTextW")
	getWindowTextLength          = user32.NewProc("GetWindowTextLengthW")
	procGetMessageExtraInfo      = user32.NewProc("GetMessageExtraInfo")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
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
	send         chan<- types.Message
	webview_hwnd win.HWND
	this_pid     uint32
	state        struct {
		running_typ types.AppType
		running_app *exec.Cmd
		file        string
		hwnd        win.HWND
	}
	explorer_killed bool
}

func NewRunner(send chan<- types.Message) (*Runner, error) {
	runner := &Runner{}
	runner.send = send

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
	runner.paths.word, err = findMicrosoftExe(word)
	if err != nil {
		send <- types.NewMessage(types.TExeNotFound{
			Name:   word,
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}
	runner.paths.excel, err = findMicrosoftExe(excel)
	if err != nil {
		send <- types.NewMessage(types.TExeNotFound{
			Name:   excel,
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}
	runner.paths.powerpoint, err = findMicrosoftExe(powerpoint)
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
	self.webview_hwnd = win.GetForegroundWindow()
	_ = win.GetWindowThreadProcessId(self.webview_hwnd, &self.this_pid)
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
func (self *Runner) OpenApp(typ types.AppType, file string) error {
	if self.isOpen() {
		return fmt.Errorf("an app is already running")
	} else {
		self.resetState()
	}
	defer self.resetState()

	// go (func() {
	// 	timeout := time.After(time.Second * 100)
	// 	for {
	// 		err := self.forceSaveInApp()
	// 		log.Println(err)
	// 		select {
	// 		case <-timeout:
	// 			return
	// 		default:
	// 			time.Sleep(time.Millisecond * 50)
	// 		}
	// 	}
	// })()

	self.state.running_typ = typ
	self.state.file = file
	switch typ {
	case types.TXT:
		return self.open(self.paths.notepad, file)
	case types.DOCX:
		return self.open(self.paths.word, file)
	case types.PPTX:
		return self.open(self.paths.powerpoint, file)
	case types.XLSX:
		return self.open(self.paths.excel, file)
	default:
		return fmt.Errorf("invalid app type: %d", typ)
	}
}

func (self *Runner) KillApp() error {
	if self.state.running_app == nil {
		return nil
	}

	// self.FocusOpenApp()
	// self.waitForFocusApp(30)
	// self.sendCtrlS()
	// _ = win.PostMessage(self.state.hwnd, win.WM_CLOSE, 0, 0)
	// self.resetState()
	// return nil

	err := self.state.running_app.Process.Kill()
	if err != nil {
		self.resetState()
	}
	return err
}

func (self *Runner) FocusOpenApp() error {
	log.Println("trying to focus open app")
	if self.state.running_app == nil {
		return nil
	}
	if self.state.running_app.Process == nil {
		return nil
	}
	log.Println("focusing app...")

	_ = win.SetForegroundWindow(self.state.hwnd)
	return nil
}

func (self *Runner) waitForFocusApp(timeout_sec int) error {
	timeout := time.After(time.Second * 30)

	fg := win.GetForegroundWindow()
	for fg != self.state.hwnd {
		select {
		case <-timeout:
			log.Println()
			return fmt.Errorf("ERROR: open app timeout")
		default:
			time.Sleep(time.Millisecond * 50)
		}

		fg = win.GetForegroundWindow()
	}

	return nil
}

func (self *Runner) sendCtrlS() {
	// VK_S := uint16(win.VkKeyScan('s'))
	// _ = win.SendMessage(self.state.hwnd, win.WM_KEYDOWN, win.VK_LCONTROL, 0)
	// _ = win.SendMessage(self.state.hwnd, win.WM_KEYUP, VK_S, 0)

	// _ = win.SendMessage(self.state.hwnd, win.WM_KEYUP, VK_S, 0)
	// _ = win.SendMessage(self.state.hwnd, win.WM_KEYUP, win.VK_LCONTROL, 0)

	// text := "Hello, World!"
	// textPtr, _ := syscall.UTF16PtrFromString(text)
	// ptr := uintptr(unsafe.Pointer(textPtr))
	// _ = win.PostMessage(self.state.hwnd, win.WM_CHAR, ptr, 0)

	// {
	// 	k1, _, _ := procGetMessageExtraInfo.Call()
	// 	k2, _, _ := procGetMessageExtraInfo.Call()
	// 	inputs := []win.KEYBD_INPUT{
	// 		{
	// 			Type: win.INPUT_KEYBOARD,
	// 			Ki: win.KEYBDINPUT{
	// 				WVk:         win.VK_CONTROL,
	// 				WScan:       0,
	// 				DwFlags:     0,
	// 				Time:        0,
	// 				DwExtraInfo: k2,
	// 			},
	// 		},
	// 		{
	// 			Type: win.INPUT_KEYBOARD,
	// 			Ki: win.KEYBDINPUT{
	// 				WVk:         0,
	// 				WScan:       's',
	// 				DwFlags:     win.KEYEVENTF_UNICODE,
	// 				Time:        0,
	// 				DwExtraInfo: k1,
	// 			},
	// 		},
	// 		// {
	// 		// 	Type: win.INPUT_KEYBOARD,
	// 		// 	Ki: win.KEYBDINPUT{
	// 		// 		WVk:         0,
	// 		// 		WScan:       's',
	// 		// 		DwFlags:     win.KEYEVENTF_UNICODE | win.KEYEVENTF_KEYUP,
	// 		// 		Time:        0,
	// 		// 		DwExtraInfo: k2,
	// 		// 	},
	// 		// },
	// 		// {
	// 		// 	Type: win.INPUT_KEYBOARD,
	// 		// 	Ki: win.KEYBDINPUT{
	// 		// 		WVk:         win.VK_CONTROL,
	// 		// 		WScan:       0,
	// 		// 		DwFlags:     win.KEYEVENTF_KEYUP,
	// 		// 		Time:        0,
	// 		// 		DwExtraInfo: 0,
	// 		// 	},
	// 		// },
	// 	}

	// 	log.Println(inputs)

	// 	// Send the input
	// 	num := win.SendInput(uint32(len(inputs)), unsafe.Pointer(&inputs[0]), int32(unsafe.Sizeof(win.KEYBDINPUT{})))
	// 	log.Println(num)

	// 	time.Sleep(time.Millisecond * 100)
	// }
	// {
	// 	// k1, _, _ := procGetMessageExtraInfo.Call()
	// 	// k2, _, _ := procGetMessageExtraInfo.Call()
	// 	k3, _, _ := procGetMessageExtraInfo.Call()
	// 	k4, _, _ := procGetMessageExtraInfo.Call()
	// 	inputs := []win.KEYBD_INPUT{
	// 		// {
	// 		// 	Type: win.INPUT_KEYBOARD,
	// 		// 	Ki: win.KEYBDINPUT{
	// 		// 		WVk:         win.VK_CONTROL,
	// 		// 		WScan:       0,
	// 		// 		DwFlags:     0,
	// 		// 		Time:        0,
	// 		// 		DwExtraInfo: k1,
	// 		// 	},
	// 		// },
	// 		// {
	// 		// 	Type: win.INPUT_KEYBOARD,
	// 		// 	Ki: win.KEYBDINPUT{
	// 		// 		WVk:         0,
	// 		// 		WScan:       's',
	// 		// 		DwFlags:     win.KEYEVENTF_UNICODE,
	// 		// 		Time:        0,
	// 		// 		DwExtraInfo: k2,
	// 		// 	},
	// 		// },
	// 		{
	// 			Type: win.INPUT_KEYBOARD,
	// 			Ki: win.KEYBDINPUT{
	// 				WVk:         0,
	// 				WScan:       's',
	// 				DwFlags:     win.KEYEVENTF_UNICODE | win.KEYEVENTF_KEYUP,
	// 				Time:        0,
	// 				DwExtraInfo: k3,
	// 			},
	// 		},
	// 		{
	// 			Type: win.INPUT_KEYBOARD,
	// 			Ki: win.KEYBDINPUT{
	// 				WVk:         win.VK_CONTROL,
	// 				WScan:       0,
	// 				DwFlags:     win.KEYEVENTF_KEYUP,
	// 				Time:        0,
	// 				DwExtraInfo: k4,
	// 			},
	// 		},
	// 	}

	// 	log.Println(inputs)

	// 	// Send the input
	// 	num := win.SendInput(uint32(len(inputs)), unsafe.Pointer(&inputs[0]), int32(unsafe.Sizeof(win.KEYBDINPUT{})))
	// 	log.Println(num)
	// }

	// text := "Hello, World!"
	// textPtr, _ := syscall.UTF16PtrFromString(text)
	// ptr := uintptr(unsafe.Pointer(textPtr))
	// _ = win.SendMessage(self.state.hwnd, win.WM_SETTEXT, 0, ptr)

	// robotgo.TypeStr("Hello, World!")
	robotgo.KeyTap("s", "ctrl")
	robotgo.KeyTap("s", "ctrl")
}

func (self *Runner) forceSaveInApp() error {
	err := self.FocusOpenApp()
	if err != nil {
		return err
	}
	err = self.waitForFocusApp(30)
	if err != nil {
		return err
	}
	self.sendCtrlS()
	log.Println("ctrl s sent")

	return nil
}

func (self *Runner) FocusOrOpenApp(typ types.AppType, file string) error {
	// if self.isOpen() && self.state.running_typ == typ && self.state.file == file {
	if self.isOpen() && self.state.running_typ == typ {
		return self.FocusOpenApp()
	} else {
		return self.OpenApp(typ, file)
	}
}

func (self *Runner) resetState() {
	self.state.file = ""
	self.state.hwnd = 0
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

// TODO: disallow alt + tab
func (self *Runner) preventDistractions(ctx context.Context) {
	killApps := func() {
		processes, err := self.ListAllProcess()
		if err != nil {
			fmt.Printf("Error listing processes: %v\n", err)
			return
		}

		fmt.Println("Running Processes (Visible on Taskbar):")
		for pid, windowText := range processes {
			fmt.Printf("PID: %d, Window Title: %s\n", pid, windowText)
		}

		appsToKill := []string{"Chrome", "Firefox", "Brave"}

		for pid, cmdline := range processes {
			for _, app := range appsToKill {
				if strings.Contains(cmdline, app) {
					fmt.Printf("Killing process %d (%s)\n", pid, cmdline)

					cmd := exec.Command("taskkill", "/PID", strconv.Itoa(int(pid)), "/F")
					err = cmd.Run()
					if err != nil {
						fmt.Printf("Error killing process %d: %v\n", pid, err)
					}
				}
			}
		}
	}

	hideActiveApps := func() {
		hwnd := win.GetForegroundWindow()
		title, _ := getWindowTitle(hwnd)
		child := win.GetParent(hwnd)
		var pid uint32
		_ = win.GetWindowThreadProcessId(hwnd, &pid)
		log.Println(title)
		if hwnd == self.webview_hwnd || child == self.webview_hwnd {
			return
		}
		if hwnd == self.state.hwnd {
			return
		}
		// this check allows any windows created within the same process
		if child == self.state.hwnd {
			return
		}

		// NOTE: this allows user to create/open new windows from the currently open app
		// if pid == uint32(self.state.running_app.Process.Pid) || pid == self.this_pid {
		// 	return
		// }

		log.Println("bad window detected")
		_ = win.SetForegroundWindow(self.webview_hwnd)
		_ = win.BringWindowToTop(self.webview_hwnd)
		_ = win.ShowWindow(hwnd, win.SW_SHOWMINIMIZED)
		_ = win.SetWindowPos(hwnd, 1, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)
		// _ = win.PostMessage(hwnd, win.WM_CLOSE, 0, 0)
		self.send <- types.NewMessage(types.TWarnUser{Message: "Unknown open application detected. Please do not open any other application during test"})
	}

	for {
		time.Sleep(time.Millisecond * 100)

		select {
		case <-ctx.Done():
			return
		default:
			_ = hideActiveApps
			// hideActiveApps()
			killApps()
		}
	}
}

func (self *Runner) open(exe string, file string) error {
	if exe == "" {
		return fmt.Errorf("executable unspecified")
	}
	if file == "" {
		return fmt.Errorf("file path unspecified")
	}

	ctx, close := context.WithCancel(context.Background())
	defer close()

	// wait for app to open and assign the hwnd to self.state
	go (func() {
		timeout := time.After(time.Second * 30)
		for {
			hwnd := win.GetForegroundWindow()
			title, _ := getWindowTitle(hwnd)
			if strings.Contains(title, tmp_prefix) {
				self.state.hwnd = hwnd
				go self.preventDistractions(ctx)
				break
			}

			select {
			case <-timeout:
				log.Println("ERROR: open app timeout")
				return
			default:
				time.Sleep(time.Millisecond * 50)
			}
		}
	})()
	// go (func() {
	// 	timeout := time.After(time.Second * 400000)
	// 	for {
	// 		select {
	// 		case <-timeout:
	// 			log.Println("ERROR: open app timeout")
	// 			return
	// 		default:
	// 			hwnd := win.GetForegroundWindow()
	// 			title, _ := getWindowTitle(hwnd)
	// 			log.Println(hwnd, title)
	// 			time.Sleep(time.Millisecond * 50)
	// 		}
	// 	}
	// })()

	// cmd := exec.Command(self.paths.explorer, file)
	// cmd := exec.Command(self.paths.cmd, "/C", "start", file)
	cmd := exec.Command(exe, file)
	self.state.running_app = cmd
	out, err := cmd.CombinedOutput()
	log.Printf("%s\n", string(out))
	log.Println(err)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (self *Runner) fullscreenForegroundWindow() {
	fg := win.GetForegroundWindow()
	_ = win.ShowWindow(fg, win.SW_MAXIMIZE)
}

func findMicrosoftExe(name string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\`+name, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue("")
	if err != nil {
		return "", err
	}

	return s, err
}

func getWindowTitle(hwnd win.HWND) (string, error) {
	win := win.GetForegroundWindow()

	// NOTE: error returned by these calls is never nil (i might be wrong, but that is what i see)
	length, _, _ := getWindowTextLength.Call(uintptr(win))
	if length == 0 {
		return "", fmt.Errorf("could not get title")
	}

	buf := make([]uint16, length+1)
	_, _, _ = getWindowText.Call(uintptr(win), uintptr(unsafe.Pointer(&buf[0])), uintptr(length+1))

	title := syscall.UTF16ToString(buf)
	return title, nil
}

// var hwnd win.HWND
// func enumWindowsProc(hWnd win.HWND, lParam uintptr) uintptr {
// 	var pid uint32
// 	win.GetWindowThreadProcessId(hWnd, &pid)
// 	name, _ := robotgo.FindName(int(pid))
// 	log.Printf("window process: %s %d\n", name, pid)
// 	if pid == uint32(lParam) {
// 		hwnd = hWnd
// 		return 0    // Stop enumeration
// 	}
// 	return 1 // Continue enumeration
// }

// func GetHWNDFromPID(pid int) (win.HWND, error) {
// 	hwnd = 0 // Reset hwnd
// 	enumWindows.Call(syscall.NewCallback(enumWindowsProc), uintptr(uint32(pid)))
// 	if hwnd == 0 {
// 		return 0, fmt.Errorf("no window found for PID %d", pid)
// 	}
// 	return hwnd, nil
// }

// func (self *Runner) fullscreenForegroundWindow() {
// 	const (
// 		SW_MAXIMIZE = 3
// 		SW_SHOW     = 5
// 	)
// 	hwnd, _, _ := getForegroundWindow.Call()
// 	showWindow.Call(hwnd, win.SW_MAXIMIZE)
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

func IsWindowVisible(hwnd syscall.Handle) bool {
	r1, _, _ := procIsWindowVisible.Call(uintptr(hwnd))
	return r1 != 0
}

func GetWindowThreadProcessId(hwnd syscall.Handle) (uint32, error) {
	var pid uint32
	_, _, _ = procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))
	return pid, nil
}

func GetWindowText(hwnd syscall.Handle) string {
	buf := make([]uint16, 200)
	getWindowText.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf)
}

type EnumWindowsProc1 func(hwnd syscall.Handle, lParam uintptr) uintptr

func EnumWindows(enumFunc EnumWindowsProc1, lParam uintptr) error {
	r1, _, err := enumWindows.Call(syscall.NewCallback(enumFunc), lParam)
	if r1 == 0 {
		return err
	}
	return nil
}

func (self *Runner) ListAllProcess() (map[uint32]string, error) {
	processes := make(map[uint32]string)

	cb := func(hwnd syscall.Handle, lParam uintptr) uintptr {
		if IsWindowVisible(hwnd) {
			pid, _ := GetWindowThreadProcessId(hwnd)
			windowText := GetWindowText(hwnd)
			if windowText != "" {
				processes[pid] = windowText
			}
		}
		return 1 // Continue enumeration
	}

	err := EnumWindows(cb, 0)
	if err != nil {
		return nil, err
	}

	return processes, nil
}
