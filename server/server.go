package server

import (
	"net/http"
	"os"

	"github.com/ckaminer/schedule-api/router"
)

func StartServer() {
	r := router.InitializeRouter()

	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	http.ListenAndServe(":"+port, r)
}
