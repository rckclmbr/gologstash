package output

import (
	// "encoding/json"
	"github.com/rckclmbr/gologstash/logstash/event"
	"fmt"
	"log"
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
	fmt.Printf(".")
	return nil
	//fmt.Println(string(j))
}
