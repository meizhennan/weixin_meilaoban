package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	http.HandleFunc("/", service.IndexHandler)
	http.HandleFunc("/api/count", service.CounterHandler)

	http.HandleFunc("/api/auto_reply", service.AutoReplyHandler)

	log.Fatal(http.ListenAndServe(":80", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
}
