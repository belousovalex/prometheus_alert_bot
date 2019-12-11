package main

import (
	"encoding/json"
	"fmt"
	"github.com/belousovalex/prometheus_alert_bot/pkg/core"
	"github.com/belousovalex/prometheus_alert_bot/pkg/parser"
	"github.com/belousovalex/prometheus_alert_bot/pkg/sender_factory"
	"log"
	"net/http"
	"regexp"
)

func processAlertMessage(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var message core.PrometheusAlertManagerMessage
	err := decoder.Decode(&message)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
	fmt.Printf("%+v\n", message)
	alertReceiver.Receive(message.Alerts)
	w.WriteHeader(http.StatusOK)
}

func main() {
	botFlags := parser.GetFlags()
	if config, err := parser.ParseConfigFile(botFlags.ConfigPath); err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Config: %+v\n", *config)
		alertReceiver = makeReceiver(config)
	}
	http.HandleFunc("/", processAlertMessage)
	if err := http.ListenAndServe(":8030", nil); err != nil {
		print(err.Error())
	}
}

func makeReceiver(config *parser.Config) core.Receiver {
	var err error
	rec := &core.AlertManagerReceiver{}
	for _, senderConfig := range config.Senders {
		var patterns *core.Patterns = nil
		if senderConfig.Patterns != nil {
			patterns = new(core.Patterns)
			patterns.Labels = make(map[string]*regexp.Regexp)
			for label, patternValue := range senderConfig.Patterns.Labels {
				patterns.Labels[label], err = regexp.Compile(patternValue)
				if err != nil {
					panic(err.Error())
				}
			}
		}
		rec.Subscribe(sender_factory.SenderFactory(senderConfig), patterns)
	}
	return rec
}

var alertReceiver core.Receiver
