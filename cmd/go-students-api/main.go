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
	"github.com/faysal0x1/go-students-api/internal/http/handlers/student"
	"github.com/faysal0x1/go-students-api/internal/storage/sqlite"
)

func main() {

	// load config

	cfg := config.MustLoad()

	// database setup

	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Database setup completed", slog.String("env", cfg.Env), slog.String("storage", cfg.StoragePath),
		slog.String("storage_path", cfg.StoragePath))

	//setup router

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students{id},student.GetById(storage)")
	router.HandleFunc("GET /api/students,student.GetList(storage)")

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

	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shut down", slog.String("error", err.Error()))
	}

	slog.Info("Server shut down successfully")

}
