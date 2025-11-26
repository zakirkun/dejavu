package queue

import (
	"encoding/json"
	"os"

	"github.com/nats-io/nats.go"
)

type Queue struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func Connect() (*Queue, error) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, err
	}

	// Create streams
	streams := []string{
		"DEPLOYMENTS",
		"BUILDS",
	}

	for _, stream := range streams {
		_, err := js.StreamInfo(stream)
		if err != nil {
			// Stream doesn't exist, create it
			_, err = js.AddStream(&nats.StreamConfig{
				Name:     stream,
				Subjects: []string{stream + ".*"},
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return &Queue{conn: conn, js: js}, nil
}

func (q *Queue) Publish(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = q.js.Publish(subject, payload)
	return err
}

func (q *Queue) Subscribe(subject string, handler func([]byte)) (*nats.Subscription, error) {
	return q.js.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
		msg.Ack()
	})
}

func (q *Queue) Close() error {
	q.conn.Close()
	return nil
}
