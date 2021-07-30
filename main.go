package main

import (
	"context"
	"fmt"
	"github.com/serhiihuberniuk/blog-api/configs"
	"github.com/serhiihuberniuk/blog-api/view/graphql/graph"
	"github.com/serhiihuberniuk/blog-api/view/graphql/graph/generated"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
	repository "github.com/serhiihuberniuk/blog-api/repository/postgresql"
	"github.com/serhiihuberniuk/blog-api/service"
	grpcHandlers "github.com/serhiihuberniuk/blog-api/view/grpc/handlers"
	"github.com/serhiihuberniuk/blog-api/view/grpc/pb"
	"github.com/serhiihuberniuk/blog-api/view/rest/handlers"
	"google.golang.org/grpc"
)

//const dbUrl = "postgres://serhii:serhii@localhost:5432/api"

func postgresConnPool(ctx context.Context, dbUrl string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("connection to database failed: %w", err)
	}

	return pool, nil
}

func main() {
	ctx := context.Background()

	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("error occured while initialisation configs: %v", err)
	}

	pool, err := postgresConnPool(ctx, config.PostgresUrl)
	if err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}

	repo := &repository.Repository{
		Db: pool,
	}

	serv := service.NewService(repo)

	handlerRest := handlers.NewRestHandlers(serv)

	restServer := http.Server{
		Addr:    ":" + config.HttpPort,
		Handler: handlerRest.ApiRouter(),
	}

	errs := make(chan error)

	// Rest server

	go func() {
		c := cors.New(cors.Options{
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		})
		handlerCors := c.Handler(restServer.Handler)

		if err := http.ListenAndServe(restServer.Addr, handlerCors); err != nil {
			errs <- err
		}
	}()

	log.Println(" Rest server is listening on ", restServer.Addr)

	// gRPC server

	address := ":" + config.GrpcPort
	grpcServer := grpc.NewServer()
	grpcHandler := grpcHandlers.NewGrpcHandlers(serv)

	go func() {
		lis, err := net.Listen("tcp", address)
		if err != nil {
			errs <- err

			return
		}

		pb.RegisterBlogApiServer(grpcServer, grpcHandler)

		if err := grpcServer.Serve(lis); err != nil {
			errs <- err
		}
	}()

	log.Println("gRPC server is listening on ", address)

	// GraphQl server
	resolver := graph.NewResolver(serv)
	srvGraphQl := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srvGraphQl)

	graphqlServer := http.Server{
		Addr:    ":" + config.GraphqlPort,
		Handler: srvGraphQl,
	}

	go func() {
		if err := http.ListenAndServe(graphqlServer.Addr, nil); err != nil {
			errs <- err
		}
	}()

	log.Printf("GraphQl server is listening on: %s with GraphQl playground", graphqlServer.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errs:
		log.Fatalf("error occurred while running server: %v", err)
	case <-quit:
	}

	log.Println("Service is shutting down...")

	pool.Close()

	if err := restServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error while shutting down: %v", err)
	}

	if err := graphqlServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error while shutting down: %v", err)
	}

	grpcServer.GracefulStop()

	log.Print("Done")
}
