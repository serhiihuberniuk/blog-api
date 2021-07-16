package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	repository "github.com/serhiihuberniuk/blog-api/repository/postgresql"
	"github.com/serhiihuberniuk/blog-api/service"
	handlers2 "github.com/serhiihuberniuk/blog-api/view/rest/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const dbUrl = "postgres://serhii:serhii@localhost:5432/api"

func postgresConnPool(ctx context.Context, dbUrl string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("connection to database failed: %w", err)
	}

	return pool, nil
}

func main() {
	ctx := context.Background()

	pool, err := postgresConnPool(ctx, dbUrl)
	if err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}

	repo := &repository.Repository{
		Db: pool,
	}

	serv := service.NewService(repo)

	handler := handlers2.NewHandlers(serv)

	srv := http.Server{
		Addr:    ":8080",
		Handler: handler.ApiRouter(),
	}

	errs := make(chan error)
	go func() {
		if err := http.ListenAndServe(srv.Addr, srv.Handler); err != nil {
			errs <- err
		}
	}()

	log.Println("server is listening")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errs:
		log.Fatalf("error occured while running HTTP server: %v", err)
	case <-quit:
	}

	log.Println("Service is shutting down...")

	pool.Close()

	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error while shutting down: %v", err)
	}

	log.Print("Done")
}
