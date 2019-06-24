package publisher

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

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

type Publisher struct {
	client *http.Client
}

func New() (*Publisher, error) {
	client, err := google.DefaultClient(context.Background(), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &Publisher{client: client}, nil
}

func NewWithClient(client *http.Client) (*Publisher, error) {
	return &Publisher{client: client}, nil
}

func (p *Publisher) Send(projectId string, topic string, message []byte) error {
	data := base64.StdEncoding.EncodeToString(message)
	m := pubsubMessage{Data: data}
	var messages []pubsubMessage
	messages = append(messages, m)
	request := publishRequest{Messages: messages}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(request)
	resp, err := p.client.Post("https://pubsub.googleapis.com/v1/projects/"+projectId+"/topics/"+topic+":publish", "application/json", buf)
	if err != nil {
		return err
	}
	log.Println(resp)
	resp.Body.Close()
	return nil
}
