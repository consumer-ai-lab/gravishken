package main

import (
	"bytes"
	"common"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var server_url string

type Client struct {
	client http.Client
	jwt    string
	user   *common.User

	server struct {
		conn         *websocket.Conn
		send         chan common.Message
		recv         chan common.Message
		conn_started bool
	}

	exit struct {
		ctx     context.Context
		destroy context.CancelFunc
	}

	frontend struct {
		send chan<- common.Message
	}
}

func newClient(send chan<- common.Message) (*Client, error) {
	self := &Client{}
	self.frontend.send = send

	self.client = http.Client{}

	self.server.send = make(chan common.Message, 100)
	self.server.recv = make(chan common.Message, 100)

	ctx, destroy := context.WithCancel(context.Background())
	self.exit.ctx = ctx
	self.exit.destroy = destroy

	return self, nil
}

func (self *Client) destroy() {
	self.closeServerConn()
	close(self.server.send)
	self.exit.destroy()
}

func (self *Client) notifyErr(err error) {
	if err != nil {
		self.frontend.send <- common.NewMessage(common.TErr{
			Message: fmt.Sprintf("Error: %s", err),
		})
		log.Printf("Error: %s\n", err)
	}
}

func (self *Client) closeServerConn() {
	if self.server.conn == nil {
		return
	}
	self.server.conn.Close()
}

func (self *Client) login(user_login *common.TUserLoginRequest) error {
	json_data, err := json.Marshal(user_login)
	if err != nil {
		return err
	}

	url := server_url + "/user/login"
	log.Println(url, string(json_data))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")

	resp, err := self.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return fmt.Errorf("%s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result common.UserLoginResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	self.jwt = result.Jwt
	self.user = &result.User
	log.Println(self.user)

	return nil
}

func (self *Client) maintainConn() {
	for {
		ctx, close := context.WithCancel(context.Background())

		err := self.connect(ctx, close)

		if err != nil {
			log.Println(err)
			close()
		}

		// block till connection breaks
		<-ctx.Done()

		msg := "server disconnected. trying reconnection in 5 seconds..."
		self.notifyErr(fmt.Errorf(msg))
		log.Println(msg)
		select {
		case <-self.exit.ctx.Done():
			log.Println("terminating connection with server")
			return
		case <-time.After(time.Second * 5):
			continue
		}
	}
}

func (self *Client) connect(exit context.Context, cancel context.CancelFunc) error {
	url, err := url.Parse(server_url)
	if err != nil {
		return err
	}
	url.Scheme = "ws"
	url.Path = "/ws"

	header := http.Header{}
	header.Add("Authorization", "Bearer "+self.jwt)

	conn, _, err := websocket.DefaultDialer.Dial(url.String(), header)
	if err != nil {
		return err
	}
	self.server.conn = conn

	go func() {
		defer cancel()
		for {
			select {
			case <-exit.Done():
				return
			case msg, ok := <-self.server.send:
				if !ok {
					return
				}
				log.Println(msg)
				self.server.conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
				err := self.server.conn.WriteJSON(msg)
				if err != nil {
					log.Println(err)
					return
				}

			}
		}
	}()
	go func() {
		defer cancel()
		for {
			var msg common.Message
			err := self.server.conn.ReadJSON(&msg)
			if err != nil {
				log.Println(err)
				return
			}
			self.server.recv <- msg
		}
	}()

	return nil
}

func (self *Client) getTests(batchName string) ([]common.Test, error) {
	// TODO: fetch all tests
	url := server_url + "/test/get_question_paper/" + batchName
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []common.Test{}, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+self.jwt)

	resp, err := self.client.Do(req)
	if err != nil {
		return []common.Test{}, err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return []common.Test{}, fmt.Errorf("%s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []common.Test{}, err
	}

	var result struct {
		Message  string      `json:"message"`
		Response common.Test `json:"response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return []common.Test{}, err
	}

	return []common.Test{result.Response}, nil
}
