package structs

import uuid "github.com/satori/go.uuid"

// MessageStruct represents the common values for all messages
type MessageStruct struct {
	Channel       string    `json:"channel"`
	SentTimestamp int64     `json:"sent_timestamp"`
	UUID          uuid.UUID `json:"uuid"`
}

// SystemMessageStruct represents the base for a message on the system channel
type SystemMessageStruct struct {
	MessageStruct
	Cmd string `json:"cmd"`
}

// InitResponseMessageStruct represents the initial message from the system upon connection
type InitResponseMessageStruct struct {
	SystemMessageStruct
	SubscriberID   uuid.UUID    `json:"subscriber_id"`
	ReconnectToken uuid.UUID    `json:"reconnect_token"`
	Subscription   Subscription `json:"subscription"`
	Reconnected    bool         `json:"reconnected"`
}

// PushMessageStruct represents a message received from the socket-server
type PushMessageStruct struct {
	MessageStruct
	CreatedTimestamp int64                  `json:"created_timestamp"`
	Payload          map[string]interface{} `json:"payload"`
}

// SeriesPayloadStruct represents the Payload in the series channel
type SeriesPayloadStruct struct {
	Type   string       `json:"type"`
	Events []string     `json:"events"`
	State  SeriesStruct `json:"state"`
	Diff   []struct {
		Attribute string `json:"attribute"`
		Before    string `json:"before"`
		After     string `json:"after"`
	} `json:"diff"`
}

type MatchPayloadStruct struct {
	Type   string      `json:"type"`
	Events []string    `json:"events"`
	State  MatchStruct `json:"state"`
	Diff   []struct {
		Attribute string `json:"attribute"`
		Before    string `json:"before"`
		After     string `json:"after"`
	} `json:"diff"`
}
