package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/faysal0x1/go-students-api/internal/config"
)

func main() {

	// load config

	cfg := config.MustLoad()

	// database setup

	//setup router

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, request *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// setup server

	server := http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Server started on", cfg.HttpServer.Addr)

	<-done

	slog.Info("Server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shut down", slog.String("error", err.Error()))
	}

	slog.Info("Server shut down successfully")

}
