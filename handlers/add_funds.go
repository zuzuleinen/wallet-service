package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"wallet-service/app"
)

type addFundsRequest struct {
	Amount    int64  `json:"amount"`
	Reference string `json:"reference"`
}

func AddFundsHandler(ws *app.WalletService, logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.PathValue("userId")

			var req addFundsRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				logger.Printf("error decoding request: %s\n", err)
				jsonError(w, "invalid request body", http.StatusBadRequest)
				return
			}
			if req.Amount < 0 {
				jsonError(w, "`amount` must be positive", http.StatusBadRequest)
				return
			}
			err = ws.HandleFunds(req.Reference, req.Amount, userId)
			if err != nil {
				logger.Printf("error adding funds: %s\n", err)
				jsonError(w, "something went wrong", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			jsonSuccess(w, "Created", http.StatusCreated)
		},
	)
}
