package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/golang-jwt/jwt"
	"github.com/rs/cors"
	"github.com/serhiihuberniuk/blog-api/configs"
	"github.com/serhiihuberniuk/blog-api/health"
	"github.com/serhiihuberniuk/blog-api/providers"
	"github.com/serhiihuberniuk/blog-api/repository/decorator"
	repository "github.com/serhiihuberniuk/blog-api/repository/postgresql"
	"github.com/serhiihuberniuk/blog-api/service"
	"github.com/serhiihuberniuk/blog-api/view/graphql/graph"
	"github.com/serhiihuberniuk/blog-api/view/graphql/graph/generated"
	graphqlMiddlewares "github.com/serhiihuberniuk/blog-api/view/graphql/middlewares"
	interceptors "github.com/serhiihuberniuk/blog-api/view/grpc/Interceptors"
	grpcHandlers "github.com/serhiihuberniuk/blog-api/view/grpc/handlers"
	"github.com/serhiihuberniuk/blog-api/view/grpc/pb"
	"github.com/serhiihuberniuk/blog-api/view/rest/handlers"
	"github.com/serhiihuberniuk/blog-api/view/rest/middlewares"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("error occurred while initialisation configs: %v", err)
	}

	pool, err := repository.NewPostgresDb(ctx, config.PostgresUrl, config.PostgresMigrationsPath, config.PostgresDatabaseVersion)
	if err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}

	repo := &repository.Repository{
		Db: pool,
	}

	repoWithCache := decorator.NewRepositoryCacheDecorator(repo, config.RedisAddress)

	privateKey, err := ioutil.ReadFile(config.PrivateKeyFile)
	if err != nil {
		log.Fatalf("cannot read Private Key from file: %v", err)
	}

	privateRSA, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		log.Fatalf("error occurred whole parsing privat key: %v", err)
	}

	userInfoProvider := providers.NewCurrentUserInformationProvider()

	serv, err := service.NewService(repoWithCache, privateRSA, userInfoProvider)
	if err != nil {
		log.Fatalf("error occurred while creating service: %v", err)
	}

	authMiddlewareProvider := providers.NewAuthInfoProvider(serv, userInfoProvider)

	errs := make(chan error)

	// Health check server

	healthHandler := health.NewHandlerHealth(repo.HealthCheck, repoWithCache.HealthCheck)

	healthServer := http.Server{
		Addr:    ":" + config.HealthcheckPort,
		Handler: healthHandler.HealthRouter(),
	}

	go func() {
		if err := http.ListenAndServe(healthServer.Addr, healthServer.Handler); err != nil {
			errs <- err
		}
	}()

	log.Println(" Health check server is listening on ", healthServer.Addr)

	middleware := middlewares.NewAuthMiddleware(authMiddlewareProvider)
	handlerRest := handlers.NewRestHandlers(serv, middleware, userInfoProvider)

	restServer := http.Server{
		Addr:    ":" + config.HttpPort,
		Handler: handlerRest.ApiRouter(),
	}

	// Rest server

	go func() {
		c := cors.New(cors.Options{
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders: []string{"Authorization", "content-type"},
		})
		handlerCors := c.Handler(restServer.Handler)

		if err := http.ListenAndServe(restServer.Addr, handlerCors); err != nil {
			errs <- err
		}
	}()

	log.Println(" Rest server is listening on ", restServer.Addr)

	// gRPC server

	address := ":" + config.GrpcPort
	authInterceptor := interceptors.NewAuthInterceptor(serv, userInfoProvider)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.UnaryAuthInterceptor),
		grpc.StreamInterceptor(authInterceptor.StreamAuthInterceptor))
	grpcHandler := grpcHandlers.NewGrpcHandlers(serv, userInfoProvider)

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
	resolverConfig := graph.NewResolverConfig(serv, userInfoProvider)

	graphMiddleware := graphqlMiddlewares.NewAuthMiddleware(authMiddlewareProvider).Auth

	srvGraphQl := handler.NewDefaultServer(generated.NewExecutableSchema(resolverConfig))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", graphMiddleware(srvGraphQl))

	graphqlServer := http.Server{
		Addr:    ":" + config.GraphqlPort,
		Handler: graphMiddleware(srvGraphQl),
	}

	go func() {
		if err := http.ListenAndServe(graphqlServer.Addr, graphqlServer.Handler); err != nil {
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

	if err := healthServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error while shutting down: %v", err)
	}

	grpcServer.GracefulStop()

	log.Print("Done")
}
