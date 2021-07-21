package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
	repository "github.com/serhiihuberniuk/blog-api/repository/postgresql"
	"github.com/serhiihuberniuk/blog-api/service"
	grpcHandlers "github.com/serhiihuberniuk/blog-api/view/grpc/handlers"
	"github.com/serhiihuberniuk/blog-api/view/grpc/pb"
	"github.com/serhiihuberniuk/blog-api/view/rest/handlers"
	"google.golang.org/grpc"
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

	handler := handlers.NewRestHandlers(serv)

	srv := http.Server{
		Addr:    ":8080",
		Handler: handler.ApiRouter(),
	}

	errsHTTP := make(chan error)

	// Rest server

	go func() {
		c := cors.New(cors.Options{
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		})
		handlerCors := c.Handler(srv.Handler)

		if err := http.ListenAndServe(srv.Addr, handlerCors); err != nil {
			errsHTTP <- err
		}
	}()

	log.Println(" Rest server is listening on ", srv.Addr)

	// gRPC server

	address := ":8081"
	grpcServer := grpc.NewServer()
	grpcHandler := grpcHandlers.NewGrpcHandlers(serv)

	errsGRPC := make(chan error)

	go func() {
		lis, err := net.Listen("tcp", address)
		if err != nil {
			errsGRPC <- err
		}

		pb.RegisterBlogApiServer(grpcServer, grpcHandler)

		if err := grpcServer.Serve(lis); err != nil {
			errsGRPC <- err
		}
	}()

	log.Println("gRPC server is listening on ", address)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errsHTTP:
		log.Fatalf("error occurred while running HTTP server: %v", err)
	case err := <-errsGRPC:
		log.Fatalf("error occurred while running gRPC server: %v", err)
	case <-quit:
	}

	log.Println("Service is shutting down...")

	pool.Close()

	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error while shutting down: %v", err)
	}

	grpcServer.GracefulStop()

	log.Print("Done")
}
