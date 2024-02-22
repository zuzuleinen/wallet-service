package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"wallet-service/domain"
	"wallet-service/infra/db"
	"wallet-service/infra/pubsub"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/google/uuid"
)

type WalletService struct {
	producer         pulsar.Producer
	transactionsRepo *db.TransactionRepository
	doneWG           sync.WaitGroup
	logger           *log.Logger
}

func NewWalletService(producer pulsar.Producer, transactionsRepo *db.TransactionRepository, logger *log.Logger) *WalletService {
	return &WalletService{
		producer:         producer,
		transactionsRepo: transactionsRepo,
		logger:           logger,
	}
}

// HandleFundsWithPulsar send the data to Pulsar broker
func (s *WalletService) HandleFundsWithPulsar(reference string, amount int64, userID string) error {
	s.doneWG.Add(1)
	errChan := make(chan error, 1)

	go func() {
		defer s.doneWG.Done()
		defer close(errChan)

		data, err := json.Marshal(&pubsub.TransactionPayload{
			Reference: reference,
			Amount:    amount,
			UserID:    userID,
		})
		if err != nil {
			errChan <- fmt.Errorf("error marshalling for producer: %s", err)
			return
		}

		_, err = s.producer.Send(context.TODO(), &pulsar.ProducerMessage{
			Payload: data,
		})
		if err != nil {
			errChan <- fmt.Errorf("error producing on pulsar: %s", err)
		}
		return
	}()

	select {
	case err := <-errChan:
		return err
	}
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
	s.doneWG.Add(1)

	walletChan := make(chan *domain.Wallet, 1)
	errChan := make(chan error, 1)

	go func() {
		defer s.doneWG.Done()

		w := s.GetWallet(userID)
		err := w.AddFunds(amount)
		if err != nil {
			errChan <- err
			close(errChan)
		}

		if err = s.transactionsRepo.AddTransaction(uuid.New().String(), userID, fmt.Sprintf("initialTopup-%s", userID), amount); err != nil {
			errChan <- err
			close(errChan)
		}
		walletChan <- s.GetWallet(userID)
		close(walletChan)
	}()

	select {
	case wallet := <-walletChan:
		return wallet, nil
	case err := <-errChan:
		return nil, err
	}
}

func (s *WalletService) Stop(ctx context.Context) {
	s.logger.Println("waiting for wallet service to finish")

	doneChan := make(chan struct{})
	go func() {
		s.doneWG.Wait()
		close(doneChan)
	}()

	select {
	case <-ctx.Done():
		s.logger.Println("context done earlier")
	case <-doneChan:
		s.logger.Println("wallet service stopped")
	}
}
