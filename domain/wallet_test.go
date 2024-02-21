package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const userID = "some-id"

func TestNewWalletInitialization(t *testing.T) {
	w := NewWallet(userID)

	assert.Equal(t, int64(0), w.Balance(), "a new wallet should have balance zero")
	assert.Equal(t, userID, w.userID, "a new wallet is correctly assigned to a user")
}

func TestAddFundsToWallet(t *testing.T) {
	w := NewWallet(userID)

	w.AddFunds(100)
	w.AddFunds(200)

	assert.Equal(t, int64(300), w.Balance(), "funds can be added to a wallet")
}

func TestRemoveFundsFromWallet(t *testing.T) {
	w := NewWallet(userID)
	w.AddFunds(100)
	w.AddFunds(-20)

	assert.Equal(t, int64(80), w.Balance(), "funds can be removed from wallet")
}
func TestBalanceOfWalletCannotBeNegative(t *testing.T) {
	w := NewWallet(userID)
	err := w.AddFunds(-30)

	assert.Error(t, err, "negative balance not allowed")
	assert.EqualError(t, err, "negative balance not allowed")
}
