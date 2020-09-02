package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

func (m Message) ToJSONString() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (m Message) PrintOut() {
	ts := time.Unix(m.Timestamp, 0).String()
	fmt.Println("Here is the message, scheduled for " + ts + ":" + m.Value)
}
