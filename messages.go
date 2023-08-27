package main

import (
	"encoding/json"
)

type ReceivedMessage struct {
	Type    string          `json:"type"`
	Id      string          `json:"id"`
	Payload json.RawMessage `json:"payload"`
}

// String is to override the default serialization of [ReceivedMessage.Payload]
// as a set of bytes to a string. This enables easier reading of the logs.
func (rm ReceivedMessage) String() string {
	serialized := `{"type":"` +
		rm.Type + `","id":"` +
		rm.Id + `","payload":` +
		string(rm.Payload) + `}`
	return serialized
}

type PayloadRegisterKey struct {
	ClientKey string `json:"client-key"`
}
