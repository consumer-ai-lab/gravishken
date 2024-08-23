package main

import (
	"bytes"
	types "common"
	user "common/models/user"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var server_url string

type Client struct {
	client http.Client
	jwt    string
}

func newClient() (*Client, error) {
	self := &Client{}

	self.client = http.Client{}

	return self, nil
}

func (self *Client) login(user_login types.TUserLogin) error {
	login_req := user.UserLoginRequest{
		Username:     user_login.Username,
		Password:     user_login.Password,
		TestPassword: user_login.TestCode,
	}

	json, err := json.Marshal(login_req)
	if err != nil {
		return err
	}

	url := server_url + "user/login"
	log.Println(url, string(json))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
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

	// TODO: maybe use the cookie jar for jwt??

	return nil
}
