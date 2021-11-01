package broker

import (
	"fmt"
	"log"
	"net"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

type BrokerConfig struct {
	port uint
	host string
}

type MqBroker struct {
}

func NewBroker() *MqBroker {
	return &MqBroker{}
}

func (*MqBroker) ListenAndServe() error {

	ln, err := net.Listen("tcp", "127.0.0.1:1883")
	if err != nil {
		// handle error
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println(err.Error())
			continue
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {

	packet, err := packets.ReadPacket(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("CONNECT MESSAGE")
	fmt.Println(packet.String())

	msg, ok := packet.(*packets.ConnectPacket)
	if !ok {
		fmt.Println("received msg that was not Connect")
		return
	}

	connack := packets.NewControlPacket(packets.Connack).(*packets.ConnackPacket)
	connack.SessionPresent = msg.CleanSession
	connack.ReturnCode = msg.Validate()

	fmt.Println(connack)

	if connack.ReturnCode != packets.Accepted {
		pub := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
		err = pub.Write(conn)
		if err != nil {
			log.Fatal("send connack error, ", msg.ClientIdentifier)
			return
		}

		return
	}

	err = connack.Write(conn)
	if err != nil {
		log.Fatal("send connack error, ", msg.ClientIdentifier)
		return
	}

	fmt.Println("HELLO ", connack)
}
