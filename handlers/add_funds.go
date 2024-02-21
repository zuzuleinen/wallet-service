package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"wallet-service/application"
)

type addFundsRequest struct {
	Amount    int64  `json:"amount"`
	Reference string `json:"reference"`
}

func AddFundsHandler(ws *application.WalletService) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.PathValue("userId")
			var req addFundsRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "Failed to parse request body", http.StatusBadRequest)
				return
			}
			if req.Amount < 0 {
				jsonError(w, "`amount` must be positive", http.StatusBadRequest)
				return
			}
			err = ws.HandleFunds(req.Reference, req.Amount, userId)
			if err != nil {
				log.Println(err)
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			jsonResponse(w, nil)
		},
	)
}
