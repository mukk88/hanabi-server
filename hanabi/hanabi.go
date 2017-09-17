package main

import (
	"log"
	"net/http"
	"github.com/mukk88/hanabi-server/sockethandler"
)

func main() {
	sh := sockethandler.NewSocketHandler()
	sh.HandleConnections()
	log.Println("starting server..")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Origin", "http://hanabi.markwooym.com")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		sh.ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(":5050", nil))
}
