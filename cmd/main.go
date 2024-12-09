package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"my-go-project/internal/config"
	"my-go-project/internal/http/handlers/student"
	"my-go-project/internal/storage/sqlite"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//load config
	cfg := config.MustLoad()

	// database set up

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	//set up router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetbyID(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	//set up server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("server start ", slog.String("Address", cfg.Addr))
	fmt.Printf("server start %s", cfg.Addr)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start")
		}

	}()
	<-done
	slog.Info("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown serever", slog.String("error", err.Error()))
	}
	slog.Info("server shut down succesfully!")

}
