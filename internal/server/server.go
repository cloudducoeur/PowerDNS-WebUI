package server

import (
	"log"
	"net/http"

	"github.com/cloudducoeur/PowerDNS-WebUI/internal/handlers"
)

func StartServer(port string) error {
	http.HandleFunc("/", handlers.ListZonesHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Server started on port %s", port)
	return http.ListenAndServe(":"+port, nil)
}
