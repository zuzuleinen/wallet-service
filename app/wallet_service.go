package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"wallet-service/db"
	"wallet-service/domain"

	"github.com/google/uuid"
)

type WalletService struct {
	transactionsRepo *db.TransactionRepository
	doneWG           sync.WaitGroup
	logger           *log.Logger
}

func NewWalletService(transactionsRepo *db.TransactionRepository, logger *log.Logger) *WalletService {
	return &WalletService{
		transactionsRepo: transactionsRepo,
		logger:           logger,
	}
}

func (s *WalletService) HandleFunds(reference string, amount int64, userID string) error {
	s.doneWG.Add(1)
	errChan := make(chan error, 1)

	go func() {
		defer s.doneWG.Done()
		defer close(errChan)

		w := s.GetWallet(userID)
		err := w.AddFunds(amount)
		if err != nil {
			errChan <- err
			return
		}

		err = s.transactionsRepo.AddTransaction(uuid.New().String(), userID, reference, amount)
		if err != nil {
			errChan <- err
		}
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