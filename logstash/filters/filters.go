package filters

import (
	"github.com/rckclmbr/gologstash/logstash/event"
)

type FilterType interface {
    Filter(evt *event.Event) error
}