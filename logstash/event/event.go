package event

import (
	"encoding/json"
	// "log"
	// "bytes"
	// "strings"
//	"fmt"
	// "reflect"
)

type Event struct {
	Tags []string `json:"@tags"`
	Fields map[string]string `json:"@fields"`
	Source string `json:"@source"`
	Type string `json:"@type"`
	Source_host string `json:"@source_host"`
	Source_path string `json:"@source_path"`
	Timestamp string `json:"@timestamp"`
	Message string `json:"@message"`
}

func NewFromJSON(data []byte) (*Event, error) {
	
	var evt Event

	jsonErr := json.Unmarshal(data, &evt)
	if jsonErr != nil {
		return nil, jsonErr
	}

	if evt.Message == "" {
		evt.Message = string(data)
	}
	
	return &evt, nil
}

func New(data []byte) (*Event) {
	var evt Event
	evt.Message = string(data)
	return &evt
}

func (evt *Event) ToJSON() ([]byte, error) {
	return json.Marshal(evt)
}

func (evt *Event) GetMessage() (string) {
	return evt.Message
}

func (evt *Event) SetField(key, value string) {
	evt.Fields[key] = value
}