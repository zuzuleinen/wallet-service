package handlers

import (
	"net/http"

	"wallet-service/app"
)

type getWalletResponse struct {
	Balance int64 `json:"balance"`
}

func GetWalletHandler(ws *app.WalletService) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.PathValue("userId")

			if !ws.HasWallet(userId) {
				jsonError(w, "wallet not found", http.StatusNotFound)
				return
			}

			var resp getWalletResponse
			resp.Balance = ws.GetWallet(userId).Balance()

			jsonResponse(w, resp)
		},
	)
}
