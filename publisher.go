package publisher

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2/google"
)

type pubsubMessage struct {
	Data        string            `json:"data"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	MessageId   string            `json:"messageId,omitempty"`
	PublishTime string            `json:"publishTime,omitempty"`
}

type publishRequest struct {
	Messages []pubsubMessage `json:"messages"`
}

type publishResponse struct {
	MessageIds []string `json:"messageIds`
}

type Publisher struct {
	client *http.Client
}

func New() (*Publisher, error) {
	client, err := google.DefaultClient(context.Background(), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, err
	}
	return &Publisher{client: client}, nil
}

func WithClient(client *http.Client) (*Publisher, error) {
	return &Publisher{client: client}, nil
}

func (p *Publisher) Send(ctx context.Context, projectId string, topic string, message []byte) (string, error) {
	data := base64.StdEncoding.EncodeToString(message)
	m := pubsubMessage{Data: data}
	var messages []pubsubMessage
	messages = append(messages, m)
	request := publishRequest{Messages: messages}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(request)

	var sb strings.Builder
	sb.WriteString("https://pubsub.googleapis.com/v1/projects/")
	sb.WriteString(projectId)
	sb.WriteString("/topics/")
	sb.WriteString(topic)
	sb.WriteString(":publish")

	req, err := http.NewRequest("POST", sb.String(), buf)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)
	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var publishResponse publishResponse
	if err := json.NewDecoder(resp.Body).Decode(&publishResponse); err != nil {
		return "", err
	}
	if len(publishResponse.MessageIds) < 1 {
		log.Println(resp)
		return "", errors.New("missing return message id")
	}

	return publishResponse.MessageIds[0], nil
}
