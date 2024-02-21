package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"wallet-service/application"
)

type createWalletRequest struct {
	Amount int64 `json:"amount"`
}

type createWalletResponse struct {
	UserID  string `json:"userId"`
	Balance int64  `json:"balance"`
}

func CreateWalletHandler(ws *application.WalletService, logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var req createWalletRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "Failed to parse request body", http.StatusBadRequest)
				return
			}

			if req.Amount < 0 {
				jsonError(w, "`amount` must be positive", http.StatusBadRequest)
				return
			}

			userId := r.PathValue("userId")
			if ws.HasWallet(userId) {
				jsonError(w, "user already has a wallet created.", http.StatusConflict)
				return
			}
			wallet, err := ws.CreateWallet(userId, req.Amount)
			if err != nil {
				logger.Printf("error creating wallet: %s\n", err)
				jsonError(w, "something went wrong", http.StatusInternalServerError)
				return
			}

			resp := createWalletResponse{
				UserID:  userId,
				Balance: wallet.Balance(),
			}

			w.WriteHeader(http.StatusCreated)
			jsonResponse(w, resp)
		},
	)
}
