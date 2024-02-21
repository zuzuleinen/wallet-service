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

func CreateWalletHandler(ws *application.WalletService) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.PathValue("userId")

			var req createWalletRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "Failed to parse request body", http.StatusBadRequest)
				return
			}

			if ws.HasWallet(userId) {
				jsonError(w, "User already has a wallet created.", http.StatusConflict)
				return
			}

			wallet, err := ws.CreateWallet(userId, req.Amount)
			if err != nil {
				// todo log error
				log.Println(err)
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
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
