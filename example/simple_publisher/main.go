package main

import (
	"crypto/tls"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//const SERVER = "tcp://broker.hivemq.com:1883"

const SERVER = "tcp://localhost:1883"
const TOPIC = "nn/sensors"
const QOS = 1

func main() {

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ClientAuth:         tls.NoClientCert,
	}

	opts := mqtt.NewClientOptions().AddBroker(SERVER)
	opts.SetClientID("mac-go")
	opts.SetTLSConfig(tlsConfig)

	client := mqtt.NewClient(opts)
	defer client.Disconnect(250)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error: ", token.Error().Error())
		return
	}

	token := client.Publish(TOPIC, QOS, false, "mymessage")
	if token.Wait() && token.Error() != nil {
		fmt.Println("Error publish: ", token.Error().Error())
	}

}
