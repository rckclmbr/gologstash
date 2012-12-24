package output

import (
	// "encoding/json"
	"github.com/rckclmbr/gologstash/logstash/event"
	"fmt"
	"log"
)

type DebugOutput struct {
	something int
}

func NewDebugOutput() (*DebugOutput) {
	return &DebugOutput{2}
}

func (o *DebugOutput) Register(args map[string]interface{}) (error) {
	return nil
}

func (o *DebugOutput) Output(evt *event.Event) (error) {
	data, err := evt.ToJSON()
	if err != nil {
		log.Printf("Error generating json: %v\n", err)
	}
	fmt.Println(string(data))
	return nil
}
