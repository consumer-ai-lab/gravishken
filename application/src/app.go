package main

import (
	"common"
	"context"
	"fmt"
	"log"
	"os"
)

type App struct {
	send   chan common.Message
	recv   chan common.Message
	runner IRunner
	client *Client

	exitCtx context.Context
	exitFn  context.CancelFunc

	state struct {
		webview_opened     bool
		connection_started bool
	}
	test_state struct {
		submitted map[common.ID]bool
		tests     map[common.ID]string
	}
}

func (self *App) destroy() {
	close(self.send)
	self.client.destroy()

	err := self.runner.RestoreEnv()
	log.Println(err)
}

func newApp() (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())
	app := &App{
		send:    make(chan common.Message, 100),
		recv:    make(chan common.Message, 100),
		exitCtx: ctx,
		exitFn:  cancel,
	}
	var err error

	client, err := newClient(app.send)
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

func (self *App) login(user_login *common.TUserLoginRequest) error {
	err := self.client.login(user_login)
	if err != nil {
		errorMessage := common.NewMessage(common.TErr{
			Message: "Failed to log in user: " + err.Error(),
		})
		self.send <- errorMessage
		return err
	}

	return nil
}

func (self *App) maintainConnection() {
	if self.state.connection_started {
		return
	}
	self.state.connection_started = true

	go self.client.maintainConn()
	go self.handleServerMessages()
}

func (self *App) handleServerMessages() {
	for {
		var msg common.Message
		var ok bool
		select {
		case <-self.exitCtx.Done():
			return
		case msg, ok = <-self.client.server.recv:
			if !ok {
				return
			}
		}

		switch msg.Typ {
		default:
			log.Printf("message type '%s' not handled ('%s')\n", msg.Typ.TSName(), msg.Val)
		}
	}
}

func (self *App) startTest() error {
	tests, err := self.client.getTests(self.client.user.Batch)
	self.client.tests = tests

	if err != nil {
		return err
	}

	routeMessage := common.TLoadRoute{
		Route: "/tests",
	}
	message := common.NewMessage(routeMessage)

	self.send <- message

	return nil
}

func (self *App) openWv() {
	var url string
	if build_mode == "DEV" {
		url = fmt.Sprintf("http://localhost:%s/", os.Getenv("DEV_PORT"))
	} else {
		url = fmt.Sprintf("http://localhost:%s/", port)
	}
	self.state.webview_opened = true

	go func() {
		openWv(url)
		self.exitFn()
	}()
}

// must be called from the main thread :/
func (self *App) wait() {
	<-self.exitCtx.Done()
}

func (self *App) notifyErr(err error) {
	if err != nil {
		self.send <- common.NewMessage(common.TErr{
			Message: fmt.Sprintf("Error: %s", err),
		})
		log.Printf("Error: %s\n", err)
	}
}

func (self *App) prepareEnv() {
	err := self.runner.SetupEnv()
	self.notifyErr(err)
}
