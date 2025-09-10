package main

import (
	"context"
	"github.com/mozhdekzm/gqlgql/internal/interface/graph"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/mozhdekzm/gqlgql/internal/config"
	"github.com/mozhdekzm/gqlgql/internal/infrastructure/migrate"
	"github.com/mozhdekzm/gqlgql/internal/infrastructure/mysql"
	"github.com/mozhdekzm/gqlgql/internal/infrastructure/redis"
	"github.com/mozhdekzm/gqlgql/internal/infrastructure/worker"
	"github.com/mozhdekzm/gqlgql/internal/usecase"
)

func main() {
	cfg := config.Load()

	// Run migrations
	migrator := migrate.NewMigrator(cfg)
	migrator.Up()

	// Setup DB and Redis
	db, err := mysql.NewBeeORMEngine(cfg)
	if err != nil {
		log.Fatalf("failed to init BeeORM engine: %v", err)
	}
	redisClient := redis.NewRedisClient(cfg)

	// Setup repositories and services
	todoRepo := mysql.NewTodoRepository(db)
	outboxRepo := mysql.NewOutboxRepository(db)
	streamPublisher := redis.NewStreamPublisher(*redisClient, cfg)
	todoService := usecase.NewTodoService(todoRepo, outboxRepo, streamPublisher, db)

	// Setup GraphQL resolver
	resolver := &graph.Resolver{
		TodoService: *todoService,
	}

	// Setup gqlgen server
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.Use(extension.Introspection{})

	// Start outbox worker
	outboxWorker := worker.NewOutboxWorker(outboxRepo, streamPublisher)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go outboxWorker.Start(ctx)

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		cancel()
		time.Sleep(2 * time.Second) // Give worker time to finish
		os.Exit(0)
	}()

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
