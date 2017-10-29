package stream

import (
	"encoding/json"
	"time"
)

// Message defines the message interface
type Message interface {
	Get(key string) interface{}
	Put(key string, value interface{})
	Replace(payload map[string]interface{})
	SetTime(timestamp int64)
	ToMap() map[string]interface{}
	String() string
}

// SimpleMessage is a default Message implementation
type SimpleMessage struct {
	Timestamp int64                  `json:"timestamp"` // In miliseconds
	Payload   map[string]interface{} `json:"payload"`
}

// Get retrieves value by key from payload
func (m *SimpleMessage) Get(key string) interface{} {
	return m.Payload[key]
}

// Put push data into message payload
func (m *SimpleMessage) Put(key string, value interface{}) {
	m.Payload[key] = value
}

// Replace set the message payload with the provided payload
func (m *SimpleMessage) Replace(payload map[string]interface{}) {
	m.Payload = payload
}

func (m *SimpleMessage) SetTime(timestamp int64) {
	m.Timestamp = timestamp
}

func (m *SimpleMessage) ToMap() map[string]interface{} {
	msg := make(map[string]interface{})
	msg["timestamp"] = m.Timestamp
	msg["payload"] = m.Payload
	// result, err := json.Marshal(m)
	// if err == nil {
	// 	return string(result)
	// }
	// panic(err)
	return msg
}

func (m *SimpleMessage) String() string {
	result, err := json.Marshal(m)
	if err == nil {
		return string(result)
	}
	panic(err)
}

// NewMessage creates a Message
func NewMessage() Message {
	return &SimpleMessage{
		Timestamp: time.Now().UnixNano() / 1000000, // In miliseconds
		Payload:   make(map[string]interface{}),
	}
}
