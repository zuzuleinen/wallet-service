package application

import (
	"fmt"
	"sync"

	"wallet-service/domain"
	"wallet-service/infrastructure"

	"github.com/google/uuid"
)

type WalletService struct {
	transactionsRepo *infrastructure.TransactionRepository
	doneWG           sync.WaitGroup
}

func NewWalletService(transactionsRepo *infrastructure.TransactionRepository) *WalletService {
	return &WalletService{
		transactionsRepo: transactionsRepo,
	}
}

func (ws *WalletService) HandleFunds(reference string, amount int64, userID string) error {
	w := ws.GetWallet(userID)
	err := w.AddFunds(amount)
	if err != nil {
		return err
	}
	return ws.transactionsRepo.AddTransaction(uuid.New().String(), userID, reference, amount)
}

func (ws *WalletService) HasWallet(userID string) bool {
	ts := ws.transactionsRepo.UserTransactions(userID)
	return len(ts) > 0 // todo can do a count here instead of fetching all
}

func (ws *WalletService) GetWallet(userID string) *domain.Wallet {
	ts := ws.transactionsRepo.UserTransactions(userID)
	w := domain.NewWallet(userID)
	for _, v := range ts {
		w.AddFunds(v.Amount)
	}
	return w
}

func (ws *WalletService) CreateWallet(userID string, amount int64) (*domain.Wallet, error) {
	w := ws.GetWallet(userID)
	err := w.AddFunds(amount)
	if err != nil {
		return &domain.Wallet{}, err
	}

	if err := ws.transactionsRepo.AddTransaction(uuid.New().String(), userID, fmt.Sprintf("initialTopup-%s", userID), amount); err != nil {
		return &domain.Wallet{}, err
	}

	return ws.GetWallet(userID), nil
}
