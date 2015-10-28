package main

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/cloudfoundry/noaa"
	"github.com/cloudfoundry/sonde-go/events"
	"metricRevealer/metricnode"
	"net/http"
	"net"
	"html/template"
	"metricRevealer/handlers"
)

var dopplerAddress = os.Getenv("DOPPLER_ADDR") // should look like ws://host:port
var authToken = os.Getenv("CF_ACCESS_TOKEN")   // use $(cf oauth-token | grep bearer)

const firehoseSubscriptionId = "firehose-a"

func main() {
	nodes := metricnode.NewNodeContainer()
	connection := noaa.NewConsumer(dopplerAddress, &tls.Config{InsecureSkipVerify: true}, nil)
	connection.SetDebugPrinter(ConsoleDebugPrinter{})

	fmt.Println("===== Streaming Firehose (will only succeed if you have admin credentials)")

	msgChan := make(chan *events.Envelope)
	go func() {
		defer close(msgChan)
		errorChan := make(chan error)
		go connection.Firehose(firehoseSubscriptionId, authToken, msgChan, errorChan)

		for err := range errorChan {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		}
	}()

	go startHttp(nodes)

	//TODO: process all event metrics that come from loggregator
	//TODO: add type of event to table
	for msg := range msgChan {
		switch msg.GetEventType() {
		case events.Envelope_ValueMetric:
			vm := msg.GetValueMetric()
			if vm == nil { continue }

			node := metricnode.NewMetricNode(*vm.Name)
			nodes.Store(node)

			fmt.Printf(nodes.Nodes[0].Value)

		case events.Envelope_CounterEvent:
			ce := msg.GetCounterEvent()
			if ce == nil { continue }
		}
	}
}

func startHttp(nodes metricnode.NodeContainer) {
	http.Handle("/messages", handlers.NewListing(nodes))
	err := http.ListenAndServe(net.JoinHostPort("", fmt.Sprintf("%d", 8818)), nil)
	if err != nil {
		println("Proxy Server Error", err)
		panic("We could not start the HTTP listener")
	}
}

type ConsoleDebugPrinter struct{}

func (c ConsoleDebugPrinter) Print(title, dump string) {
	println(title)
	println(dump)
}

