package app

import (
	"fmt"
	"net/http"
)

func (a *App) loadRoutes() {
	a.router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Welcome to the basic http server")
		if err != nil {
			return
		}
	})
}
