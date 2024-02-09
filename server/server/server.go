package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	port = ":443"
	ip   = "127.0.0.1"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	logFilePath := filepath.Join(currentDir, "JmaLive", "server", "logs", "server.log")

	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("os.OpenFile() filed with '%s\n", err)
	}
	defer logFile.Close()

	log := log.New(logFile, "Server: ", log.LstdFlags|log.Lshortfile)
	log.SetOutput(logFile)
	log.Println("Log entry")

	// Handle Funcs

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test successful!"))
	})

	// Handle Funcs

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v", sig)
		cancel()
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	fmt.Println("Starting the HTTP server...")
	log.Printf("HTTPS server is starting with address: %s%s", ip, port)
	server := &http.Server{Addr: port}
	go func() {
		if err := server.ListenAndServeTLS("JmaLive/server/server/cert.pem", "JmaLive/server/server/key.pem"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
			fmt.Printf("Error starting server: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down the server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	// Some more end function
}
