package main

import (
	types "common"
	"fmt"
	"log"
	"os"

	webview "github.com/thrombe/webview_go"
)

type App struct {
	send    chan types.Message
	recv    chan types.Message
	runner  IRunner
	webview webview.WebView
	client  *Client

	state struct {
		webview_opened bool
	}
	jwt string
}

func (self *App) destroy() {
	close(self.send)

	err := self.runner.RestoreEnv()
	log.Println(err)

	if self.state.webview_opened {
		self.webview.Destroy()
	}
}

func newApp() (*App, error) {
	app := &App{
		send:    make(chan types.Message, 100),
		recv:    make(chan types.Message, 100),
		webview: nil,
	}
	var err error

	client, err := newClient()
	if err != nil {
		return nil, err
	}
	app.client = client

	app.runner, err = NewRunner(app.send)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (self *App) login(user_login types.TUserLogin) error {
	jwt, err := self.client.login(user_login)
	if err != nil {
		errorMessage := types.NewMessage(types.TErr{
			Message: "Failed to log in user: " + err.Error(),
		})
		self.send <- errorMessage
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

func (self *App) startTest(testData types.TGetTest) error {
	questionPaper, err := self.client.getTest(testData)

	if err != nil {
		return err
	}

	fmt.Println("Question paper: ", questionPaper)

	routeMessage := types.TLoadRoute{
		Route: "/tests/1",
	}
	message := types.NewMessage(routeMessage)

	self.send <- message

	return nil
}

// here we have to start the microsoft word, excel, powerpoint application
// func (self *App) StartMicrosoftApps(microSoftApp types.TMicrosoftApps) error {
// 	runner, err := NewRunner()
// 	if err != nil {
// 		return err
// 	}

// 	err = self.runner.StartMicrosoftApps(runner, microSoftApp)
// 	if err != nil {
// 		return err
// 	}

// 	routeMessage := types.TLoadRoute{
// 		Route: "/tests/2",
// 	}
// 	message := types.NewMessage(routeMessage)

// 	self.send <- message

// 	return nil
// }

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

	err := self.runner.SetupEnv()
	self.notifyErr(err)
}
