package output

import (
	"github.com/rckclmbr/gologstash/logstash/event"
	"log"
	"github.com/mattbaird/elastigo/core"
	"time"
	"fmt"
)

type ElasticSearch struct {
	something int
}

func NewElasticSearch() (*ElasticSearch) {
	return &ElasticSearch{1}
}

func (es *ElasticSearch) Register(args map[string]interface{}) (error) {
	return nil
}

func (es *ElasticSearch) Output(evt *event.Event) (error) {
	_, err := evt.ToJSON()
	if err != nil {
		log.Printf("Error generating json: %v\n", err)
	}

	t := time.Now()
	index := fmt.Sprintf("logstash-%d.%02d.%02d",
	    t.Year(),
	    t.Month(),
	    t.Day())

	data, err := evt.ToJSON()

	response, err := core.Index(true, index, evt.Type, "", string(data))
	if err != nil {
		log.Printf("Error: %+v %v\n", response, err)
		return err
	}
	return nil
	//fmt.Println(string(j))
}