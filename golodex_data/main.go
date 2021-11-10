package main

import (
	"context"
	"internal/golodexdata"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/readiness", readinessHandler)
	http.HandleFunc("/api/login", golodexdata.Login)
	http.HandleFunc("/", golodexdata.DataPage)
	//http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("./pages"))))
	//http.ListenAndServe(golodexdata.Property("GOLODEX_PORT", ":8095"), nil)
	srv := &http.Server{
		Handler:      http.DefaultServeMux,
		Addr:         ":"+golodexdata.Property("GOLODEX_PORT", "8095"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("running..."+golodexdata.Property("GOLODEX_PORT", ":8095"))
	//Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	//Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}