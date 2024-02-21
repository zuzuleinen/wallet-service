package handlers

import (
	"fmt"
	"net/http"
)

func HealthHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "OK")
		},
	)
}
