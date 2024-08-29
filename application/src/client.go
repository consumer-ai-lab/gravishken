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

)

var server_url string

type Client struct {
	client http.Client
}

func newClient() (*Client, error) {
	self := &Client{}

	self.client = http.Client{}

	return self, nil
}

func (self *Client) login(user_login types.TUserLogin) (string, error) {
	login_req := user.UserLoginRequest{
		Username:     user_login.Username,
		Password:     user_login.Password,
		TestPassword: user_login.TestCode,
	}

	json_data, err := json.Marshal(login_req)
	if err != nil {
		return "", err
	}

	url := server_url + "user/login"
	log.Println(url, string(json_data))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")

	resp, err := self.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return "", fmt.Errorf("%s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Message  string `json:"message"`
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Response, nil
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
		Message  string `json:"message"`
		Response TEST.Test `json:"response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return TEST.Test{}, err
	}

	return result.Response, nil
}


func (self *Client) StartMicrosoftApps(runner *LinuxRunner,  microsftApp types.TMicrosoftApps) error{

	switch microsftApp.AppName {
	case "Word":
		runner.runLibreOffice()
	case "NotePad":
		runner.runNotepad()
	case "PowerPoint":
		runner.runPowerPoint()
	default:
		return fmt.Errorf("Invalid Microsoft App")
	}

	return nil
	
}