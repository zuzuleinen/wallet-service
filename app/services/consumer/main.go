package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"wallet-service/infra/db"
	"wallet-service/infra/pubsub"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/google/uuid"
)

const TopicTransactions = "transactions"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Init database
	database, err := db.InitDatabase("dev.db")
	if err != nil {
		log.Fatalf("connecting to database: %s", err)
	}
	transactionsRepo := db.NewTransactionRepository(database)
	log.Println("database started")

	defer func() {
		log.Println("stopping database")
		sqlDB, _ := database.DB()
		sqlDB.Close()
	}()

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	defer client.Close()
	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            TopicTransactions,
		SubscriptionName: "wallet-transactions",
		Type:             pulsar.Exclusive,
		Name:             fmt.Sprintf("wallet-service-consumer"),
	})

	defer consumer.Close()

	// keep message stats
	msgReceived := int64(0)
	bytesReceived := int64(0)

	// print stats of the consume rate every 10 secs
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()

	for {
		select {
		case cm, ok := <-consumer.Chan():
			if !ok {
				break
			}

			var msg pubsub.TransactionPayload
			err := json.Unmarshal(cm.Message.Payload(), &msg)
			if err != nil {
				log.Printf("error on unmarshal: %s\n", err)
				break
			}

			err = transactionsRepo.AddTransactionWithCtx(ctx,
				uuid.New().String(),
				msg.UserID,
				msg.Reference,
				msg.Amount,
			)
			if err != nil {
				log.Printf("error saving to db: %s\n", err)
				break
			}

			// measure
			msgReceived++
			bytesReceived += int64(len(cm.Message.Payload()))

			err = consumer.Ack(cm.Message)
			if err != nil {
				log.Printf("error ack message: %s", cm.Message.ID())
			}
		case <-tick.C:
			currentMsgReceived := atomic.SwapInt64(&msgReceived, 0)
			currentBytesReceived := atomic.SwapInt64(&bytesReceived, 0)
			msgRate := float64(currentMsgReceived) / float64(10)
			bytesRate := float64(currentBytesReceived) / float64(10)

			log.Printf(`Stats - Consume rate: %6.1f msg/s - %6.1f Mbps`,
				msgRate, bytesRate*8/1024/1024)
		case <-ctx.Done():
			fmt.Println("Context is done. Exiting...")
			return
		}
	}
}
