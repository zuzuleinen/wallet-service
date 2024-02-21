package domain

type Wallet struct {
	userID  string
	balance int64
}

func NewWallet(userID string) *Wallet {
	return &Wallet{
		userID:  userID,
		balance: 0,
	}
}

func (w *Wallet) Balance() int64 {
	return w.balance
}

func (w *Wallet) UserID() string {
	return w.userID
}

type NegativeBalancerErr struct {
}

func (e *NegativeBalancerErr) Error() string {
	return "negative balance not allowed"
}

// AddFunds updates Wallet balance with amount
//
// To add funds you can use positive amount
// To remove funds you can use negative amount
// Adding funds that will result in a negative balance returns a NegativeBalancerErr
func (w *Wallet) AddFunds(amount int64) error {
	finalBalance := w.balance + amount
	if finalBalance < 0 {
		return &NegativeBalancerErr{}
	}
	w.balance = finalBalance
	return nil
}
