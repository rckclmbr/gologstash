package main


import (
	"log"
	"os/signal"
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

type LogstashEvents struct {
	filters []filters.FilterType
	inputs []input.InputType
	outputs []output.OutputType
}

func readConfig(filename string) (interface{}) {
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

   	return jsontype
}

func (le *LogstashEvents) parseConfig() {
	args := readConfig("./test.json").(map[string]interface{})

   	filters_config := args["filters"].([]interface{})
   	inputs_config := args["inputs"].([]interface{})
   	outputs_config := args["outputs"].([]interface{})

   	for _, filter_config := range filters_config {
   		f := filter_config.(map[string]interface{})

   		switch f["type"] {
   		case "grok":
			g := filters.NewGrokFilter()
			g.Register(f)
			le.filters = append(le.filters, g)
		}
	}

   	for _, input_config := range inputs_config {
   		f := input_config.(map[string]interface{})

   		switch f["type"] {
   		case "zeromq":
			g := input.NewZeroMQ()
			g.Register(f)
			le.inputs = append(le.inputs, g)
		}
	}

   	for _, output_config := range outputs_config {
   		f := output_config.(map[string]interface{})
   		var g output.OutputType
   		switch f["type"] {
   		case "elasticsearch": g = output.NewElasticSearch()
		case "debug": g = output.NewDebugOutput()
		}
		if g != nil {
			g.Register(f)
			le.outputs = append(le.outputs, g)
		}
	}

	fmt.Printf("%v", le.outputs[0])
}

func (le *LogstashEvents) Input(output chan *event.Event) {
	tmpchan := make(chan []byte)		
	for _, input := range le.inputs {
		go input.Receive(tmpchan)
	}
	for {
		msg := <- tmpchan
	  	evt, err := event.NewFromJSON(msg)
	  	if err != nil {
			evt = event.New(msg)
		}
		output <- evt
	}
}

func (le *LogstashEvents) Filter(input chan *event.Event, output chan *event.Event) {
	start := input
	end := make(chan *event.Event)		
		
	for i, f := range le.filters {
		if i == len(le.filters)-1 {
			end = output
		}
		go func(s chan *event.Event, e chan *event.Event, f filters.FilterType) {
			for {
				evt := <- s
				f.Filter(evt)
				e <- evt
			}
		}(start, end, f)
		start = end // Start at the last end
		end = make(chan *event.Event)
		i += 1
	}
}

func (le *LogstashEvents) Output(out chan *event.Event) {

	outputs := make([]chan *event.Event, 0)
	for _, out := range le.outputs {
		tmp := make(chan *event.Event)
		outputs = append(outputs, tmp)
		go func(out output.OutputType, ch chan *event.Event) {
			for { out.Output(<- ch) }
		}(out, tmp)
	}

	for {
		evt := <- out
		for _, o := range outputs {
			o <- evt
		}
	}
}

func main() {
	fmt.Println("Starting logstash")
		
	inchan := make(chan *event.Event)
	outchan := make(chan *event.Event)

	le := &LogstashEvents{
		make([]filters.FilterType, 0), 
		make([]input.InputType, 0),
		make([]output.OutputType, 0),
	}
	le.parseConfig()
	go le.Input(inchan)
	go le.Filter(inchan, outchan)
	go le.Output(outchan)

	// Wait for SIGINT
	c := make(chan os.Signal, 1)                                       
	signal.Notify(c, os.Interrupt)
	t0 := time.Now()
    log.Printf("Captured %v, Stopping and exiting.", <- c)
    log.Printf("Runtime: %v", time.Now().Sub(t0))
}
