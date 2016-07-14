package docubotlib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client represents a Docubot API Client
type Client struct {
	DocubotAPIURLBase string
	DocubotAPIKey     string
	DocubotAPISecret  string
}

// NewClient initializes a docubot client struct
func NewClient(url string, key string, secret string) *Client {
	return &Client{
		DocubotAPIURLBase: url,
		DocubotAPIKey:     key,
		DocubotAPISecret:  secret,
	}
}

// MessageResponse is the response received from a message sent to docubot
type MessageResponse struct {
	Data MessageResponseData `json:"data"`
	Meta MessageResponseMeta `json:"meta"`
}

// MessageResponseData is the data received from a message sent to docubot
type MessageResponseData struct {
	Messages []string `json:"messages"`
	Complete bool     `json:"complete"`
}

// MessageResponseMeta is the meta received from a message sent to docubot
type MessageResponseMeta struct {
	ThreadID string `json:"threadId"`
	UserID   string `json:"userId"`
}

// MessageResponseError is the response when there is an error
type MessageResponseError struct {
	Errors []string `json:"errors"`
}

// SendMessage sends a message to docubot
func (c *Client) SendMessage(message string, thread string, sender string) (*MessageResponse, error) {
	jsonStr, _ := json.Marshal(
		map[string]interface{}{
			"message": message,
			"thread":  thread,
			"sender":  sender,
		},
	)
	url := fmt.Sprintf("%v/api/v1/docubot", c.DocubotAPIURLBase)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.DocubotAPIKey, c.DocubotAPISecret)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var error MessageResponseError
		json.NewDecoder(resp.Body).Decode(&error)
		e := "Unknown error occurred"
		if len(error.Errors) > 0 {
			e = error.Errors[0]
		}
		return nil, errors.New(e)
	}
	var response MessageResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	return &response, err
}

// GetDocubotDoc gets the docubot document
func (c *Client) GetDocubotDoc(thread string, user string) (io.ReadCloser, error) {
	params := url.Values{}
	params.Set("user", user)
	url := fmt.Sprintf(
		"%v/api/v1/docubot/%v/doc/download?%v",
		c.DocubotAPIURLBase,
		thread,
		params.Encode(),
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.DocubotAPIKey, c.DocubotAPISecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()
		var error MessageResponseError
		json.NewDecoder(resp.Body).Decode(&error)
		e := "Unknown error occurred"
		if len(error.Errors) > 0 {
			e = error.Errors[0]
		}
		return nil, errors.New(e)
	}
	return resp.Body, nil
}
