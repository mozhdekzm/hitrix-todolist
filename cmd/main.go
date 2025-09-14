package main

import (
	"context"
	"git.ice.global/packages/hitrix"
	"git.ice.global/packages/hitrix/service"
	"git.ice.global/packages/hitrix/service/component/app"
	"git.ice.global/packages/hitrix/service/registry"
	"github.com/mozhdekzm/hitrix-todolist/config"
	"github.com/mozhdekzm/hitrix-todolist/internal/domain"
	"github.com/mozhdekzm/hitrix-todolist/internal/interface/graph"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/mozhdekzm/hitrix-todolist/internal/infrastructure/migrate"
	"github.com/mozhdekzm/hitrix-todolist/internal/infrastructure/mysql"
	"github.com/mozhdekzm/hitrix-todolist/internal/infrastructure/redis"
	"github.com/mozhdekzm/hitrix-todolist/internal/infrastructure/worker"
	"github.com/mozhdekzm/hitrix-todolist/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Run migrations
	migrator := migrate.NewMigrator(cfg)
	migrator.Up()

	// Setup DB
	//db, err := mysql.NewBeeORMEngine(cfg)
	//if err != nil {
	//	log.Fatalf("failed to init BeeORM engine: %v", err)
	//}

	// Setup Redis
	redisClient := redis.NewRedisClient(cfg)

	// Build Hitrix server and register global services including Redis pools
	s, deferFunc := hitrix.New(
		"my-app", "secret",
	).RegisterDIGlobalService(
		registry.ServiceProviderGenerator(),
		registry.ServiceProviderErrorLogger(),
		registry.ServiceProviderConfigDirectory("config"),
		registry.ServiceProviderClock(),
		registry.ServiceProviderOrmRegistry(domain.Init),
		registry.ServiceProviderOrmEngine(),
	).RegisterDIRequestService(
		registry.ServiceProviderOrmEngineForContext(false),
	).RegisterRedisPools(&app.RedisPools{
		Persistent: "default",
		Cache:      "default",
	}).Build()
	defer deferFunc()

	// Start background processors **after Redis pool is registered**
	b := &hitrix.BackgroundProcessor{Server: s}
	b.RunAsyncOrmConsumer()
	b.RunAsyncRequestLoggerCleaner()

	// Setup repositories and services
	ormengine := service.DI().OrmEngine()
	todoRepo := mysql.NewTodoRepository(ormengine)
	outboxRepo := mysql.NewOutboxRepository(ormengine)
	streamPublisher := redis.NewStreamPublisher(*redisClient, cfg)
	todoService := usecase.NewTodoService(todoRepo, outboxRepo, streamPublisher, ormengine)

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
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	port, _ := strconv.Atoi(cfg.ServerPort)
	s.RunServer(uint(port), graph.NewExecutableSchema(graph.Config{Resolvers: resolver}), nil, nil)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
