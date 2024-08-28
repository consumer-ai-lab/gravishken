package main

import (
	types "common"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	webview "github.com/thrombe/webview_go"
)

type App struct {
	send    chan types.Message
	recv    chan types.Message
	runner  Runner
	webview webview.WebView
	client  *Client

	state struct {
		webview_opened  bool
		explorer_killed bool
	}
	jwt string
}

func (self *App) destroy() {
	close(self.send)
	if self.state.explorer_killed {
		self.runner.startExplorer()
	}
	if self.state.webview_opened {
		self.webview.Destroy()
	}
}

func newApp() (*App, error) {
	app := &App{
		send:    make(chan types.Message, 100),
		recv:    make(chan types.Message, 100),
		runner:  Runner{},
		webview: nil,
	}
	var err error

	client, err := newClient()
	if err != nil {
		return nil, err
	}
	app.client = client

	if runtime.GOOS != "windows" {
		return app, nil
	}

	app.runner.paths.cmd, err = exec.LookPath(cmd)
	if err != nil {
		return nil, err
	}
	app.runner.paths.kill, err = exec.LookPath(kill)
	if err != nil {
		return nil, err
	}
	app.runner.paths.explorer, err = exec.LookPath(explorer)
	if err != nil {
		return nil, err
	}
	app.runner.paths.notepad, err = exec.LookPath(notepad)
	if err != nil {
		app.send <- types.NewMessage(types.TExeNotFound{
			Name:   notepad,
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}
	app.runner.paths.word, err = exec.LookPath(word)
	if err != nil {
		app.send <- types.NewMessage(types.TExeNotFound{
			Name:   word,
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}
	app.runner.paths.excel, err = exec.LookPath(excel)
	if err != nil {
		app.send <- types.NewMessage(types.TExeNotFound{
			Name:   excel,
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}
	app.runner.paths.powerpoint, err = exec.LookPath(powerpoint)
	if err != nil {
		app.send <- types.NewMessage(types.TExeNotFound{
			Name:   powerpoint,
			ErrMsg: fmt.Sprintf("%s", err),
		})
		err = nil
	}

	return app, nil
}

func (self *App) login(user_login types.TUserLogin) error {
	jwt, err := self.client.login(user_login)
	if err != nil {
		return err
	}
	self.jwt = jwt

    routeMessage := types.TLoadRoute{
        Route: "/instructions",
    }
    message := types.NewMessage(routeMessage)

    self.send <- message

	return nil
}

func (self *App) openWv() {
	self.webview = webview.New(build_mode == "DEV")
	self.state.webview_opened = true
}

func (self *App) wait() {
	self.webview.Run()
}

func (self *App) notifyErr(err error) {
	if err != nil {
		self.send <- types.NewMessage(types.TErr{
			Message: fmt.Sprintf("Error: %s", err),
		})
		log.Printf("Error: %s\n", err)
	}
}

func (self *App) prepareEnv() {
	self.webview.SetTitle("gravishken")

	if build_mode == "DEV" {
		url := fmt.Sprintf("http://localhost:%s/", os.Getenv("DEV_PORT"))
		self.webview.Navigate(url)
	} else {
		url := fmt.Sprintf("http://localhost:%s/", port)
		self.webview.Navigate(url)
	}

	if runtime.GOOS != windows {
		return
	}

	err := self.runner.killExplorer()
	self.notifyErr(err)
	self.state.explorer_killed = err == nil

	self.runner.disableTitlebar()
	self.runner.fullscreenForegroundWindow()
}
