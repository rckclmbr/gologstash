package main


import (
	"fmt"
	"encoding/json"
	"time"
	"io/ioutil"
	"os"
	"github.com/rckclmbr/gologstash/logstash/filters"
	"github.com/rckclmbr/gologstash/logstash/input"
	"github.com/rckclmbr/gologstash/logstash/output"
	"github.com/rckclmbr/gologstash/logstash/event"
)

type Message struct {
    name string
    test string
    time int64
}

func ReadConfig(filename string) (interface{}) {
	file, e := ioutil.ReadFile(filename)
    if e != nil {
       fmt.Printf("File error: %v\n", e)
       os.Exit(1)
   }
   var jsontype interface{}
   e = json.Unmarshal(file, &jsontype)
   if e != nil {
      fmt.Printf("Error reading config file '%v': %v\n", filename, e)
      os.Exit(1)
  }
   fmt.Printf("Results: %v\n", jsontype)
   return jsontype
}

func Input(output chan *event.Event) {
		
	filter := make(chan []byte)
	
	zmq := input.NewZeroMQ()
	go zmq.Receive(filter)
		
	for {
		msg := <- filter
	  	evt, err := event.NewFromJSON(msg)
	  	if err != nil {
			evt = event.New(msg)
		}
		output <- evt
	}
}

func Filter(input chan *event.Event, output chan *event.Event) {
	
	
	args := ReadConfig("./test.json").(map[string]interface{})
	grok := filters.NewGrokFilter()
	grok.Register(args)
	
	runMany(input, output, func(evt *event.Event) {
		grok.Filter(evt)
	}, 20)

}

func runMany2(input chan *event.Event, output chan *event.Event, f func(*event.Event), iterations int) {
	for {
		evt := <- input
		for i:=0; i<iterations; i++ {
			f(evt)
		}
		output <- evt
	}
}

func runMany(input chan *event.Event, output chan *event.Event, f func(*event.Event), iterations int) {
	start := input
	end := make(chan *event.Event)		
		
	for i:=0; i<iterations; i++ {
		if i == iterations-1 {
			end = output
		}
		go func(s chan *event.Event, e chan *event.Event) {
			for {
				evt := <- s
				f(evt)
				e <- evt
			}
		}(start, end)
		start = end // Start at the last end
		end = make(chan *event.Event)
	}
}

func Output(out chan *event.Event) {
	es := output.NewElasticSearch()
	
	t0 := time.Now()
	for i:=0; i<10000; i++ {
		evt := <- out
		es.Output(evt)
	}	
	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))

}

func main() {
	fmt.Println("Hello, 世界")
	
	ReadConfig("test.json")
		
	input := make(chan *event.Event)
	output := make(chan *event.Event)
	go Input(input)
	go Filter(input, output)
	Output(output)
}
