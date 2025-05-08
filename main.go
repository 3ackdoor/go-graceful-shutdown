package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	start time.Time
)

func init() {
	start = time.Now()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		vars := r.URL.Query()
		sleep := vars.Get("sleep")

		fmt.Printf("sleep val=%v\n", cmp.Or(sleep, "0"))
		if len(sleep) > 0 {
			t, err := strconv.Atoi(sleep)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			time.Sleep(time.Duration(t) * time.Second)
		}
		w.Write([]byte("Hello World"))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cmp.Or(os.Getenv("PORT"), "8080")),
		Handler: mux,
		// Enforce server timeout
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	if len(os.Getenv("NO_SIGNALS")) > 0 {
		log.Printf("Started server in %v", time.Since(start))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v\n", err)
		}
	} else {
		done := make(chan bool, 1)
		go gracefulShutdown(done, server)

		log.Printf("Started server in %v", time.Since(start))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v\n", err)
		}

		// Wait for the shutdown process to complete
		<-done
		log.Println("Shutdown completed")
	}
}

func gracefulShutdown(done chan bool, server *http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal
	<-ctx.Done()

	gracePeriodDuration, _ := strconv.Atoi(cmp.Or(os.Getenv("GRACE_PERIOD_DURATION"), "30"))
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(gracePeriodDuration)*time.Second, /* Kubernetes termination grace period time is 30 seconds by default */
	)
	defer cancel()

	log.Println("Server is shutting down...")

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server is forced to shutdown with error: %v\n", err)
	} else {
		log.Printf("Server has been shutdown\n")
	}

	// Notify the main goroutine that shutdown process is completed
	done <- true
}
