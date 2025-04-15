package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cloudducoeur/PowerDNS-WebUI/internal/handlers"
)

func StartServer(listenAddress, port string) error {
	http.HandleFunc("/", handlers.ListZonesHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	address := fmt.Sprintf("%s:%s", listenAddress, port)
	log.Printf("Server started on %s", address)
	return http.ListenAndServe(address, nil)
}
