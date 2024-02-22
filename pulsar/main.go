package main

import (
	"context"
	"fmt"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
)

const TopicTransactions = "transactions"

func main() {
	consumer()
}

func consumer() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	defer client.Close()

	cons, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            TopicTransactions,
		SubscriptionName: "my-sub",
		Type:             pulsar.Exclusive,
	})
	// defer cons.Close()

	msg, err := cons.Receive(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
		msg.ID(), string(msg.Payload()))

	fmt.Println("ACK")
	err = cons.Ack(msg)
	if err != nil {
		log.Fatal(err)
	}

}
