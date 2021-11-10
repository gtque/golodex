package main

import (
	"context"
	"fmt"
	"internal/golodex"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/readiness", readinessHandler)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/", golodex.Page)
	//http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("./pages"))))
	//http.ListenAndServe(golodex.Property("GOLODEX_PORT", ":8090"), nil)
	fmt.Println("running...")

	srv := &http.Server{
		Handler:      http.DefaultServeMux,
		Addr:         ":"+golodex.Property("GOLODEX_PORT", "8090"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("running..."+golodex.Property("GOLODEX_PORT", ":8090"))
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