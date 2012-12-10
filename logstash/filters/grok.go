package filters

import (
	"strings"
	cgrock "github.com/rckclmbr/gogrok/grok"
	"github.com/rckclmbr/gologstash/logstash/event"
	// "fmt"
)

type GrokData struct {
	grok *cgrock.Grok
	pattern string
	pattern_files []string
}

// Called once when the application initializes
func NewGrokFilter() (*GrokData)  {
	grok := cgrock.New()
	return &GrokData{grok, "", nil}
}

func (gd *GrokData) parseConfig(args map[string]interface{}) {
	gd.pattern = args["pattern"].(string)
	for _, i := range args["pattern_files"].([]interface{}) {
		gd.pattern_files = append(gd.pattern_files, i.(string))
	}
}

// Called once for every time a "grok" config is used
func (gd *GrokData) Register(args map[string]interface{}) (error) {

	gd.parseConfig(args)
	
	err := gd.grok.AddPattern("WORD", "\\w+")	
	if err != nil {
		return err
	}
	err = gd.grok.Compile(gd.pattern)
	if err != nil {
		return err
	}
	return nil
}

// Called once for every line that will be processed
func (gd *GrokData) Filter(evt *event.Event) (error) {
	//fmt.Println(evt.GetMessage())
	match, err := gd.grok.Match(evt.Message)
	if err != nil {
		return err
	}
	
	for k, v := range match {
		if strings.Contains(k, ":") {
			newkey := strings.SplitN(k, ":", 2)[1]
			evt.SetField(newkey, v)
		}
	}
	return nil
}
