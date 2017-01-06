package dingtalk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrServerError = errors.New("server error")
	ErrBadRequest  = errors.New("bad request")
)

const (
	robotEndpointFormatter = "https://oapi.dingtalk.com/robot/send?access_token=%s"
	messageTypeText        = "text"
	messageTypeMarkdown    = "markdown"
)

type Client struct {
	robotEndpoint string
}

func getErrFromHttpResponse(resp *http.Response) error {
	if resp.StatusCode > http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
		return ErrBadRequest
	}
	if resp.StatusCode > http.StatusInternalServerError {
		return ErrServerError
	}
	return nil
}

func NewClient(accessToken string) *Client {
	return &Client{robotEndpoint: fmt.Sprintf(robotEndpointFormatter, accessToken)}
}

func (c *Client) SendText(content string, atMobiles []string, isAtAll bool) error {
	message := Message{
		Type: messageTypeText,
		Text: Text{
			Content: content,
		},
		At: At{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	}
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	resp, err := http.Post(c.robotEndpoint, "application/json;charset=utf-8", bytes.NewReader(b))
	if err != nil {
		return err
	}
	return getErrFromHttpResponse(resp)
}

func (c *Client) SendMarkdown(title, text string, atMobiles []string, isAtAll bool) error {
	message := Message{
		Type: messageTypeMarkdown,
		Markdown: Markdown{
			Title: title,
			Text:  text,
		},
		At: At{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	}
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	resp, err := http.Post(c.robotEndpoint, "application/json;charset=utf-8", bytes.NewReader(b))
	if err != nil {
		return err
	}
	return getErrFromHttpResponse(resp)
}
