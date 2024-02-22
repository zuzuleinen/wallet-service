package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"wallet-service/app"
	"wallet-service/domain"
)

type removeFundsRequest struct {
	Amount    int64  `json:"amount"`
	Reference string `json:"reference"`
}

func RemoveFundsHandler(ws *app.WalletService) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.PathValue("userId")
			var req removeFundsRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "Failed to parse request body", http.StatusBadRequest)
				return
			}

			if req.Amount < 0 {
				jsonError(w, "`amount` must be positive", http.StatusBadRequest)
				return
			}

			err = ws.HandleFunds(req.Reference, req.Amount*-1, userId)
			if err != nil {
				var negativeBalancerErr *domain.NegativeBalancerErr
				if errors.As(err, &negativeBalancerErr) {
					jsonError(w, "adding `amount` will result in a negative balance", http.StatusBadRequest)
					return
				}
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			jsonSuccess(w, "Removed", http.StatusOK)
		},
	)
}
