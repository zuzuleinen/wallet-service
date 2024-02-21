package handlers

import (
	"net/http"

	"wallet-service/application"
)

type getWalletResponse struct {
	Balance int64 `json:"balance"`
}

func GetWalletHandler(ws *application.WalletService) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.PathValue("userId")

			wallet := ws.GetWallet(userId)

			var resp getWalletResponse
			resp.Balance = wallet.Balance()

			jsonResponse(w, resp)
		},
	)
}
