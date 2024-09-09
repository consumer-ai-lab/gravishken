package main

import (
	"bytes"
	types "common"
	TEST "common/models/test"
	user "common/models/user"
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

	server struct {
		conn *websocket.Conn
		send chan types.Message
		recv chan types.Message
	}
}

func newClient() (*Client, error) {
	self := &Client{}

	self.client = http.Client{}

	self.server.send = make(chan types.Message, 100)
	self.server.recv = make(chan types.Message, 100)

	return self, nil
}

func (self *Client) destroy() {
	self.closeServerConn()
}

func (self *Client) closeServerConn() {
	self.server.conn.Close()
	close(self.server.send)
}

func (self *Client) login(user_login *types.TUserLogin) error {
	login_req := user.UserLoginRequest{
		Username:     user_login.Username,
		Password:     user_login.Password,
		TestPassword: user_login.TestCode,
	}

	json_data, err := json.Marshal(login_req)
	if err != nil {
		return err
	}

	url := server_url + "user/login"
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

	var result struct {
		Message  string `json:"message"`
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	self.jwt = result.Response

	return nil
}

// TODO: check if this auto reconnects or reconnection has to be handled separately
func (self *Client) connect(username string) error {
	// TODO: with jwt auth
	url, err := url.Parse(server_url)
	if err != nil {
		return err
	}
	url.Scheme = "ws"
	url.Path = "/ws"
	url.RawQuery = "username=" + username

	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return err
	}
	self.server.conn = conn

	go func() {
		for {
			msg, ok := <-self.server.send
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
	}()
	go func() {
		defer close(self.server.recv)
		for {
			var msg types.Message
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

func (self *Client) getTest(testData types.TGetTest) (TEST.Test, error) {
	test_code := testData.TestPassword
	url := server_url + "test/get_question_paper/" + test_code
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return TEST.Test{}, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("test_code", test_code)

	resp, err := self.client.Do(req)
	if err != nil {
		return TEST.Test{}, err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return TEST.Test{}, fmt.Errorf("%s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TEST.Test{}, err
	}

	var result struct {
		Message  string    `json:"message"`
		Response TEST.Test `json:"response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return TEST.Test{}, err
	}

	return result.Response, nil
}
