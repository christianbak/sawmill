package datadog

import (
	"fmt"
	"github.com/pagerduty/godspeed"
	"github.com/phemmer/sawmill/event"
)

// All of the exported attribues are safe to replace before the handler has been added into a logger.
type DDHandler struct {
	client *godspeed.Godspeed
	tags   []string
}

func New(ddUrl string, tags []string) (*DDHandler, error) {
	sw := &DDHandler{}
	sw.tags = tags

	var err error
	sw.client, err = godspeed.New(ddUrl, godspeed.DefaultPort, false)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to datadog agent. %s", err)
	}

	return sw, nil
}

func (sw *DDHandler) Event(logEvent *event.Event) error {

	title := fmt.Sprintf("[%s] %s", logEvent.Level, logEvent.Message)
	text := fmt.Sprintf("%s", logEvent.FlatFields)

	// // the optionals are for the optional arguments available for an event
	// // http://docs.datadoghq.com/guides/dogstatsd/#fields
	optionals := make(map[string]string)

	switch logEvent.Level {
	case event.Debug:
	case event.Alert:
	case event.Info:
		optionals["alert_type"] = "info"
	case event.Emergency:
	case event.Error:
	case event.Critical:
		optionals["alert_type"] = "err"
	case event.Notice:
	case event.Warning:
		optionals["alert_type"] = "warning"
	}
	optionals["source_type_name"] = "application"

	//addlTags := []string{"service:firebirds"}

	err := sw.client.Event(title, text, optionals, sw.tags)

	fmt.Println(title)
	fmt.Println(text)

	if err != nil {
		return err
	}

	return nil
}

func (s *DDHandler) Stop() {
	s.client.Conn.Close()
}
