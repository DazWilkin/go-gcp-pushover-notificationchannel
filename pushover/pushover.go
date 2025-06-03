package pushover

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	url string = "https://api.pushover.net/1/messages.json"
)

type Client struct {
	Client  *http.Client
	Device  string
	Token   string
	UserKey string
}

func New(device, userKey, token string) *Client {
	return &Client{
		Client:  &http.Client{},
		Device:  device,
		Token:   token,
		UserKey: userKey,
	}
}

func (c *Client) SendMessage(title, message string) error {
	log.Println("[pushover:Client.NewSendMessage] Entered")

	rqstMsg := &SendMessageRequest{
		Device:  c.Device,
		Token:   c.Token,
		User:    c.UserKey,
		Title:   title,
		Message: message,
		HTML:    "1",
	}
	rqstBytes, err := json.Marshal(rqstMsg)
	if err != nil {
		return err
	}

	log.Printf("[pushover:Client.NewSendMessage] SendMessageRequest:\n%s", string(rqstBytes))

	rqstBody := bytes.NewReader(rqstBytes)
	rqst, err := http.NewRequest("POST", url, rqstBody)
	if err != nil {
		log.Println("[pushover:Client.NewSendMessage] Unable to create new request")
		return err
	}

	rqst.Header.Add("Content-Type", "application/json")

	resp, err := c.Client.Do(rqst)
	if err != nil {
		log.Println("[pushover:Client.NewSendMessage] Unable to do request")
		return err
	}

	log.Printf("[pushover:Client.NewSendMessage] Response code: %d", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[pushover:Client.NewSendMessage] Unable to read response body")
		return err
	}

	log.Printf("[pushover:Client.NewSendMessage] SendMessageResponse:\n%s", string(respBody))

	if resp.StatusCode != http.StatusOK {
		respErr := &SendMessageErrorResponse{}
		if err := json.Unmarshal(respBody, respErr); err != nil {
			log.Println("[pushover:Client.NewSendMessage] Unable to unmarshal error response body")
			return err
		}

		log.Printf("[pushover:Client.NewSendMessage] Request: %s\nStatus: %+v", respErr.Request, respErr.Errors)
	}

	respMsg := &SendMessageResponse{}
	if err := json.Unmarshal(respBody, respMsg); err != nil {
		log.Println("[pushover:Client.NewSendMessage] Unable to unmarshal response body")
		return err
	}

	log.Printf("[pushover:Client.NewSendMessage] Request: %s (Status: %d)", respMsg.Request, respMsg.Status)
	return nil
}
