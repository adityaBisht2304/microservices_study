package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/microservices_study/handlers"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// Initializing Handlers
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	sm := http.NewServeMux()

	// Defining urls and handlers which will handle those requests
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	// Initializing Server
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// Moving to a goroutine so that ListenAndServe do not block the execution in main
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// OS Signals like Ctrl+Z, Ctrl+C
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// This channel will block here until OS signals are received
	sig := <-sigChan
	l.Println("Received terminate now - next step is graceful shutdown", sig)

	// Graceful Shutdown after receiving signal
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
