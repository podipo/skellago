package be

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

/*
	A client for interacting with the Skella back end web API
*/
type Client struct {
	BaseURL string
	Schema  Schema
	Session string
	User    User
}

/*
	Create a client for interacting with the Skella back end web API
	baseURL: a fully qualified URL to the API like http://127.0.0.1:9000/api/0.1.0
*/
func NewClient(baseURL string) (*Client, error) {
	client := &Client{
		BaseURL: baseURL,
	}
	err := client.fetchSchema()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (client *Client) Authenticate(email string, password string) error {
	// Post the login info
	loginData := LoginData{
		Email:    email,
		Password: password,
	}
	resp, err := client.PostJSON("/user/current", loginData)
	if err != nil {
		return err
	}

	// Look for the session cookie
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == TEST_SESSION_COOKIE {
			client.Session = cookie.Value
		}
	}
	if client.Session == "" {
		return errors.New("No session cookie on the authentication response")
	}

	// Read the User data
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&client.User)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) Deauthenticate() error {
	if client.Session == "" {
		return nil
	}
	client.Session = ""
	req, err := client.prepJSONRequest("DELETE", client.BaseURL+"/user/current", nil)
	if err != nil {
		return err
	}
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	logger.Print("Header ", resp.Header)
	return nil
}

func (client *Client) prepJSONRequest(method string, url string, data []byte) (req *http.Request, err error) {
	if data == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewReader(data))
	}
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", AcceptHeaderPrefix+client.Schema.API.Version)
	if client.Session != "" {
		req.AddCookie(&http.Cookie{
			Name:  TEST_SESSION_COOKIE,
			Value: client.Session,
		})
	}
	return req, nil
}

func (client *Client) GetList(url string) (*APIList, error) {
	c := &http.Client{}
	req, err := client.prepJSONRequest("GET", client.BaseURL+url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Non-200 error " + strconv.Itoa(resp.StatusCode) + " getting list from " + url)
	}
	defer resp.Body.Close()
	list := new(APIList)
	err = json.NewDecoder(resp.Body).Decode(list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (client *Client) GetJSON(url string, target interface{}) error {
	c := &http.Client{}
	req, err := client.prepJSONRequest("GET", client.BaseURL+url, nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Non-200 error " + strconv.Itoa(resp.StatusCode) + " getting JSON from " + url)
	}
	// Read the User data
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) PostJSON(url string, data interface{}) (resp *http.Response, err error) {
	return client.SendJSON("POST", url, data)
}

func (client *Client) PutJSON(url string, data interface{}) (resp *http.Response, err error) {
	return client.SendJSON("PUT", url, data)
}

func (client *Client) SendJSON(method string, url string, data interface{}) (resp *http.Response, err error) {
	c := &http.Client{}
	dataBuff, err := json.Marshal(data)
	if err != nil {
		return
	}
	req, err := client.prepJSONRequest(method, client.BaseURL+url, dataBuff)
	if err != nil {
		return
	}
	resp, err = c.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return resp, errors.New("Non-200 error " + strconv.Itoa(resp.StatusCode) + " " + method + "ing JSON to " + url)
	}
	return
}

func (client *Client) UpdateUser(user *User) error {
	resp, err := client.PutJSON("/user/"+user.UUID, user)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) fetchSchema() error {
	resp, err := http.Get(client.BaseURL + "/schema")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&client.Schema)
	if err != nil {
		return err
	}
	return nil
}
