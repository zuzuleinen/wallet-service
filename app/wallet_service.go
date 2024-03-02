package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"wallet-service/domain"
	"wallet-service/infra/db"
	"wallet-service/infra/pubsub"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/google/uuid"
)

type WalletService struct {
	producer         pulsar.Producer
	transactionsRepo *db.TransactionRepository
	logger           *log.Logger
}

func NewWalletService(producer pulsar.Producer, transactionsRepo *db.TransactionRepository, logger *log.Logger) *WalletService {
	return &WalletService{
		producer:         producer,
		transactionsRepo: transactionsRepo,
		logger:           logger,
	}
}

// HandleFunds send the data to Pulsar broker
func (s *WalletService) HandleFunds(reference string, amount int64, userID string) error {
	data, err := json.Marshal(&pubsub.TransactionPayload{
		Reference: reference,
		Amount:    amount,
		UserID:    userID,
	})
	if err != nil {
		return fmt.Errorf("error marshalling for producer: %s", err)
	}

	f := func(id pulsar.MessageID, message *pulsar.ProducerMessage, err error) {
		if err != nil {
			s.logger.Printf("error publishing message: %s", err)
		}
	}
	s.producer.SendAsync(context.TODO(), &pulsar.ProducerMessage{
		Payload: data,
	}, f)

	return nil
}

func (s *WalletService) HasWallet(userID string) bool {
	ts := s.transactionsRepo.UserTransactions(userID)
	return len(ts) > 0
}

func (s *WalletService) GetWallet(userID string) *domain.Wallet {
	ts := s.transactionsRepo.UserTransactions(userID)
	w := domain.NewWallet(userID)
	for _, v := range ts {
		w.AddFunds(v.Amount) // todo handle error
	}
	return w
}

func (s *WalletService) CreateWallet(userID string, amount int64) (*domain.Wallet, error) {
	w := s.GetWallet(userID)
	err := w.AddFunds(amount)
	if err != nil {
		return nil, fmt.Errorf("error adding funds: %s", err)
	}

	if err = s.transactionsRepo.AddTransaction(uuid.New().String(), userID, fmt.Sprintf("initialTopup-%s", userID), amount); err != nil {
		return nil, fmt.Errorf("error adding transaction: %s", err)
	}
	return w, nil
}
