package output

import (
	"github.com/rckclmbr/gologstash/logstash/event"
)

type OutputType interface {
    Output(evt *event.Event) error
    Register(args map[string]interface{}) (error)
}