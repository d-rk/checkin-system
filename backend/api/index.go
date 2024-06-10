package api

import (
	"github.com/d-rk/checkin-system/pkg/server"
	"net/http"
)

// Handler is the entrypoint for the vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	db := server.NewDB(false)
	defer db.Close()

	router := server.NewRouter(r.Context(), db)
	router.ServeHTTP(w, r)
}
