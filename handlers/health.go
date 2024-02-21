package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func HealthHandler(logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Println(fmt.Sprintf("%s %s", r.Method, r.RequestURI))
			fmt.Fprintln(w, "OK")
		},
	)
}
