package pubsub

const TopicTransactions = "transactions"

// TransactionPayload is the JSON representation that is saved in the Tr
type TransactionPayload struct {
	Reference string `json:"reference"`
	UserID    string `json:"userId"`
	Amount    int64  `json:"amount"`
}
