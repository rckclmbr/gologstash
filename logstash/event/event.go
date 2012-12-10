package event

import (
	"encoding/json"
	"log"
	"bytes"
	"strings"
//	"fmt"
	"reflect"
)

type Event struct {
	Tags []string
	Fields map[string]string
	Source string
	Type string
	Source_host string
	Source_path string
	Timestamp string
	Message string
}

func NewFromJSON(data []byte) (*Event, error) {
	
	var evt Event
	evt.Fields = make(map[string]string)
	dec := json.NewDecoder(bytes.NewReader(data))
	s_Struct := reflect.ValueOf(&evt).Elem()
	
	// Will be overwritten below if @message exists
	evt.Message = string(data)
	
	for {
		var v map[string]interface{}
		if err := dec.Decode(&v); err != nil {
			break
		} else if err != nil {
			log.Println(err)
			break
		}
		
		for k, v := range v {
			if strings.HasPrefix(k, "@") {
				k = k[1:]
			}
			if k == "fields" {
				for j, w := range v.(map[string]interface{}) {
					evt.Fields[j] = w.(string)
				}
			} else if k == "tags" {
				//evt.Tags = v.([]interface{})
				for _, j := range v.([]interface{}) {
					evt.Tags = append(evt.Tags, j.(string))
				}
			} else {
				var field = string(strings.ToUpper(string(k[0]))) + k[1:]
				f := s_Struct.FieldByName(field)
				if f.IsValid() {
					f.SetString(v.(string))
				} else {
					evt.Fields[k] = v.(string)					
				}
			}
		}
	}
	
	return &evt, nil
}

func New(data []byte) (*Event) {
	var evt Event
	evt.Message = string(data)
	return &evt
}

func (evt *Event) ToJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"@tags": evt.Tags,
		"@fields": evt.Fields,
		"@source": evt.Source,
		"@type": evt.Type,
		"@source_host": evt.Source_host,
		"@source_path": evt.Source_path,
		"@timestamp": evt.Timestamp,
		"@message": evt.Message,
	})
}

func (evt *Event) GetMessage() (string) {
	return evt.Message
}

func (evt *Event) SetField(key, value string) {
	evt.Fields[key] = value
}