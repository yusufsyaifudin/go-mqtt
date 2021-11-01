package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var f = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Println(string(msg.Payload()))
}

const SERVER = "tcp://broker.hivemq.com:1883"

//const SERVER = "tcp://localhost:1883"
const TOPIC = "nn/sensors"
const QOS = 1

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	opts := MQTT.NewClientOptions().AddBroker(SERVER)
	opts.SetClientID("mac-go-subscriber") // this must be different than publisher client id

	opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(TOPIC, QOS, f); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}

	<-c
}
