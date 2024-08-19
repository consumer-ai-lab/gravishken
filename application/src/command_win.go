//go:build windows
// +build windows

package main

import (
	"fmt"
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

// func (self *Runner) fullscreenForegroundWindow() {

// 	const (
// 		SW_MAXIMIZE = 3
// 		SW_SHOW     = 5
// 	)

// 	hwnd, _, _ := getForegroundWindow.Call()
// 	showWindow.Call(hwnd, SW_MAXIMIZE)
// 	setWindowPos.Call(hwnd, 0, 0, 0, 1920, 1080, 0)
// }

func gadsgadd() {
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
		fmt.Printf("HWND: %v, Title: %s\n", hwnd, title)
	}

	return 1 // Continue enumeration
}

func disableTitlebar() {
	hwnd, _, _ := getForegroundWindow.Call()

	style, _, _ := getWindowLong.Call(hwnd, GWL_STYLE)

	newStyle := style &^ (win.WS_CAPTION | win.WS_BORDER | win.WS_DLGFRAME)
	newStyle = 0x000000
	newStyle |= win.WS_POPUP | win.WS_VISIBLE

	_, _, err := setWindowLong.Call(hwnd, GWL_STYLE, newStyle)
	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("Error setting window style:", err)
	} else {
		fmt.Println("Title bar and borders removed successfully.")
	}

	syscall.SyscallN(uintptr(user32.NewProc("SetWindowPos").Addr()), 5,
		hwnd, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)
}
