package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"wallet-service/infra/db"
	"wallet-service/infra/pubsub"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/google/uuid"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Init database
	database, err := db.InitDatabase("dev.db")
	if err != nil {
		log.Fatalf("connecting to database: %s", err)
	}

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
		Topic:            pubsub.TopicTransactions,
		SubscriptionName: "wallet-transactions",
		Type:             pulsar.Exclusive,
		Name:             fmt.Sprintf("wallet-service-consumer"),
	})

	defer consumer.Close()

	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()

	var ts []db.Transaction

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

			ts = append(ts, db.Transaction{
				ID:        uuid.New().String(),
				Reference: msg.Reference,
				UserID:    msg.UserID,
				Amount:    msg.Amount,
				CreatedAt: time.Now(),
			})

			err = consumer.AckID(cm.ID()) // todo
			if err != nil {
				log.Printf("error on ACK: %s", err)
			}
		case <-tick.C:
			if len(ts) > 0 {
				chunkSize := 300
				for i := 0; i < len(ts); i += chunkSize {
					end := i + chunkSize
					if end > len(ts) {
						end = len(ts)
					}
					chunk := ts[i:end]

					tx := database.Begin()
					if tx.Error != nil {
						panic(tx.Error)
					}

					fmt.Println("flush", len(chunk))
					if err := tx.Create(&chunk).Error; err != nil {
						tx.Rollback()
						panic(err)
					}

					tx.Commit()
				}
				ts = make([]db.Transaction, 0)
			} else {
				fmt.Println("Nothing to flush...", len(ts))
			}
		case <-ctx.Done():
			fmt.Println("Context is done. Exiting...")
			return
		}
	}
}
