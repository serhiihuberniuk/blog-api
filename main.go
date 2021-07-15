package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	repository "github.com/serhiihuberniuk/blog-api/repository/postgresql"
	"github.com/serhiihuberniuk/blog-api/service"
	"github.com/serhiihuberniuk/blog-api/view/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const dbUrl = "postgres://serhii:serhii@localhost:5432/api"

func PostgresConnPool(ctx context.Context, dbUrl string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("connection to database failed: %w", err)
	}

	return pool, nil
}

func main() {
	ctx := context.Background()

	pool, err := PostgresConnPool(ctx, dbUrl)
	if err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}

	repo := &repository.Repository{
		Db: pool,
	}

	serv := service.NewService(repo)

	handler := handlers.NewHandlers(serv)

	srv := http.Server{
		Addr:    ":8080",
		Handler: handler.ApiRouter(),
	}

	go func() {
		if err := http.ListenAndServe(srv.Addr, srv.Handler); err != nil {
			log.Fatalf("error while starting server: %v", err)
		}
	}()

	log.Println("server is listening")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	killSignal := <-quit

	switch killSignal {
	case syscall.SIGINT:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}

	log.Println("Service is shutting down...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("error while shutting down: %v", err)
	}

	pool.Close()

	log.Print("Done")
}
