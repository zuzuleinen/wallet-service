package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
)

type Payload struct {
	Reference string `json:"reference"`
	UserID    string `json:"userId"`
	Amount    int64  `json:"amount"`
}

const TopicTransactions = "transactions"

func main() {
	producer()
}

func producer() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	defer client.Close()
	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}

	p := Payload{
		Reference: "won-bet-333",
		UserID:    "andrei",
		Amount:    600,
	}
	data, err := json.Marshal(&p)
	if err != nil {
		log.Fatalf("error on marshal: %s", err)
	}

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: TopicTransactions,
	})
	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: data,
	})
	defer producer.Close()

	if err != nil {
		fmt.Println("Failed to publish message", err)
	} else {
		fmt.Println("Published message")
	}
}
