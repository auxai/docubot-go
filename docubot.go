package docubotlib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const unknownErrorMessage string = "Unknown error occurred"

// DocumentTree is a data model
type DocumentTree struct {
	ID            string        `json:"id"`
	DocumentName  string        `json:"documentName"`
	EntryQuestion *QuestionNode `json:"entryQuestion,omitempty"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	CreatedAt     time.Time     `json:"createdAt"`
}

// QuestionCondition is a data model
type QuestionCondition struct {
	VariableName string `json:"variableName"`
	Comparator   string `json:"comparator"`
	Value        string `json:"value"`
}

// QuestionNode is a data model
type QuestionNode struct {
	VariableName    string                `json:"variableName"`
	Question        string                `json:"question"`
	LogicalOperator string                `json:"logicalOperator"`
	Conditions      []QuestionCondition   `json:"conditions"`
	EntityType      string                `json:"entityType"`
	ChildQuestions  []QuestionNode        `json:"childQuestions"`
	MetaData        *QuestionNodeMetaData `json:"metaData,omitempty"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	CreatedAt       time.Time             `json:"createdAt"`
}

// QuestionNodeMetaData is a data model
type QuestionNodeMetaData struct {
	// Choices is what holds the choices of a multiple choice entity
	Choices map[string]string `json:"choices,omitempty"`
}

// Document is a data model
type Document struct {
	ID             string    `json:"id"`
	DocumentTreeID string    `json:"documentTreeId"`
	HeaderHTML     string    `json:"headerHtml,omitempty"`
	BodyHTML       string    `json:"bodyHtml,omitempty"`
	FooterHTML     string    `json:"footerHtml,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedAt      time.Time `json:"createdAt"`
}

// Client represents a Docubot API Client
type Client struct {
	DocubotAPIURLBase        string
	DocubotPreviewAPIURLBase string
	DocubotAPIKey            string
	DocubotAPISecret         string
}

// NewClient initializes a docubot client struct
func NewClient(url string, key string, secret string) *Client {
	return &Client{
		DocubotAPIURLBase:        url,
		DocubotPreviewAPIURLBase: url,
		DocubotAPIKey:            key,
		DocubotAPISecret:         secret,
	}
}

// PreviewMessageResponse is the response received from a preview message sent to docubot
type PreviewMessageResponse struct {
	Data PreviewMessageResponseData `json:"data"`
	Meta MessageResponseMeta        `json:"meta"`
}

// PreviewMessageResponseData is the data reveived from a preview message sent to docubot
type PreviewMessageResponseData struct {
	Messages    []string               `json:"messages"`
	Complete    bool                   `json:"complete"`
	HasDocument bool                   `json:"hasDocument"`
	Variables   map[string]interface{} `json:"variables"`
}

// MessageResponse is the response received from a message sent to docubot
type MessageResponse struct {
	Data MessageResponseData `json:"data"`
	Meta MessageResponseMeta `json:"meta"`
}

// MessageResponseData is the data received from a message sent to docubot
type MessageResponseData struct {
	Messages    []string `json:"messages"`
	HasDocument bool     `json:"hasDocument"`
	Complete    bool     `json:"complete"`
}

// MessageResponseMeta is the meta received from a message sent to docubot
type MessageResponseMeta struct {
	ThreadID        string                            `json:"threadId"`
	UserID          string                            `json:"userId"`
	DocumentName    string                            `json:"documentName"`
	MessageMetaData map[string]map[string]interface{} `json:"messageMetaData"`
}

// MessageResponseError is the response when there is an error
type MessageResponseError struct {
	Errors []string `json:"errors"`
}

// DocumentURLResponse is the response received from getting a document's URL from docubot
type DocumentURLResponse struct {
	Data DocumentURLData        `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

// DocumentURLData is the response data received from getting a document's URL from docubot
type DocumentURLData struct {
	URL string `json:"url"`
}

// DocumentVariablesResponse is the response received from getting a document's Variables from docubot
type DocumentVariablesResponse struct {
	Data DocumentVariablesData  `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

// DocumentVariablesData is the response data received from getting a document's Variables from docubot
type DocumentVariablesData struct {
	Variables map[string]interface{} `json:"variables"`
}

// SendMessage sends a message to docubot
func (c *Client) SendMessage(message string, thread string, sender string, docTreeID string) (*MessageResponse, error) {
	jsonStr, _ := json.Marshal(
		map[string]interface{}{
			"message":   message,
			"thread":    thread,
			"sender":    sender,
			"docTreeId": docTreeID,
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
		e := unknownErrorMessage
		if len(error.Errors) > 0 {
			e = error.Errors[0]
		}
		return nil, errors.New(e)
	}
	var response MessageResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	return &response, err
}

// SendPreviewMessage sends a preview message to docubot, this is a message that isn't stored on docubot at all
func (c *Client) SendPreviewMessage(message string, variables map[string]interface{}, docTree *DocumentTree) (*PreviewMessageResponse, error) {
	jsonStr, _ := json.Marshal(
		map[string]interface{}{
			"message":   message,
			"docTree":   docTree,
			"variables": variables,
		},
	)
	url := fmt.Sprintf("%v/api/v1/preview", c.DocubotPreviewAPIURLBase)
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
		e := unknownErrorMessage
		if len(error.Errors) > 0 {
			e = error.Errors[0]
		}
		return nil, errors.New(e)
	}
	var response PreviewMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	return &response, err
}

// GetPreviewDoc gets a preview document that isn't stored permanently
func (c *Client) GetPreviewDoc(variables map[string]interface{}, document *Document) (io.ReadCloser, error) {
	jsonStr, _ := json.Marshal(
		map[string]interface{}{
			"document":  document,
			"variables": variables,
		},
	)
	url := fmt.Sprintf(
		"%v/api/v1/preview/doc",
		c.DocubotPreviewAPIURLBase,
	)
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
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()
		var error MessageResponseError
		json.NewDecoder(resp.Body).Decode(&error)
		e := unknownErrorMessage
		if len(error.Errors) > 0 {
			e = error.Errors[0]
		}
		return nil, errors.New(e)
	}
	return resp.Body, nil
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
		e := unknownErrorMessage
		if len(error.Errors) > 0 {
			e = error.Errors[0]
		}
		return nil, errors.New(e)
	}
	return resp.Body, nil
}

// GetDocubotDocURL gets the docubot document url
func (c *Client) GetDocubotDocURL(thread string, user string, exp time.Duration) (*DocumentURLResponse, error) {
	params := url.Values{}
	params.Set("user", user)
	params.Set("duration", fmt.Sprintf("%v", int(exp.Seconds())))
	url := fmt.Sprintf(
		"%v/api/v1/docubot/%v/doc/url?%v",
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
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()
		var error MessageResponseError
		json.NewDecoder(resp.Body).Decode(&error)
		e := unknownErrorMessage
		if len(error.Errors) > 0 {
			e = error.Errors[0]
		}
		return nil, errors.New(e)
	}
	var response DocumentURLResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	return &response, err
}

// GetDocubotVariables gets the docubot variables for the provided user in the provided thread
func (c *Client) GetDocubotVariables(thread string, user string) (*DocumentVariablesResponse, error) {
	params := url.Values{}
	params.Set("user", user)
	url := fmt.Sprintf(
		"%v/api/v1/docubot/%v/variables?%v",
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
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()
		var error MessageResponseError
		json.NewDecoder(resp.Body).Decode(&error)
		e := unknownErrorMessage
		if len(error.Errors) > 0 {
			e = error.Errors[0]
		}
		return nil, errors.New(e)
	}
	var response DocumentVariablesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	return &response, err
}
