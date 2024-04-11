package testlogger

import (
	"encoding/json"
)

type Message struct {
	Level     int            `json:"level"`
	Message   string         `json:"message"`
	KeyValues map[string]any `json:"-"`
}

func (msg Message) MarshalJSON() ([]byte, error) {
	type alias Message
	msgJSON, err := json.Marshal(alias(msg))
	if err != nil {
		return nil, err
	}

	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(msgJSON, &raw); err != nil {
		return nil, err
	}

	for k, v := range msg.KeyValues {
		value, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		raw[k] = value
	}

	return json.Marshal(raw)
}

func (msg Message) String() string {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(msgJSON)
}

func (msg Message) Keys() []string {
	keys := make([]string, 0, len(msg.KeyValues))
	for k := range msg.KeyValues {
		keys = append(keys, k)
	}
	return keys
}

func (msg Message) HasKey(key string) bool {
	_, ok := msg.KeyValues[key]
	return ok
}

func newMessage(level int, message string, keyValues ...any) Message {
	msg := Message{
		Level:     level,
		Message:   message,
		KeyValues: make(map[string]any),
	}
	for i := 0; i < len(keyValues); i += 2 {
		key, ok := keyValues[i].(string)
		if !ok {
			panic("key must be a string")
		}
		msg.KeyValues[key] = keyValues[i+1]
	}
	return msg
}
