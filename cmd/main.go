package main

import (
	"context"
	"git.ice.global/packages/hitrix"
	"git.ice.global/packages/hitrix/service/component/app"
	"git.ice.global/packages/hitrix/service/registry"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/gin-gonic/gin"
	"github.com/mozhdekzm/gqlgql/internal/config"
	"github.com/mozhdekzm/gqlgql/internal/domain"
	"github.com/mozhdekzm/gqlgql/internal/infrastructure/migrate"
	"github.com/mozhdekzm/gqlgql/internal/infrastructure/mysql"
	"github.com/mozhdekzm/gqlgql/internal/infrastructure/redis"
	"github.com/mozhdekzm/gqlgql/internal/infrastructure/worker"
	"github.com/mozhdekzm/gqlgql/internal/interface/graph"
	"github.com/mozhdekzm/gqlgql/internal/usecase"

	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	cfg := config.Load()

	// Run migrations
	migrator := migrate.NewMigrator(cfg)
	migrator.Up()

	// Hitrix Server
	s, deferFunc := hitrix.New(
		"gqlgql",
		"secret",
	).RegisterDIGlobalService(
		registry.ServiceProviderErrorLogger(),
		registry.ServiceProviderConfigDirectory("config"),
		registry.ServiceProviderOrmRegistry(domain.Init),
		registry.ServiceProviderOrmEngine(),
		registry.ServiceProviderClock(),
		registry.ServiceProviderJWT(),
		registry.ServiceProviderPassword(),
	).RegisterDIRequestService(
		registry.ServiceProviderOrmEngineForContext(false),
	).RegisterRedisPools(&app.RedisPools{
		Persistent: "default",
		Cache:      "default",
	}).Build()
	defer deferFunc()

	// Setup DB & Redis manually (می‌تونی اینا رو هم توی DI بندازی)
	db, err := mysql.NewBeeORMEngine(cfg)
	if err != nil {
		log.Fatalf("failed to init BeeORM engine: %v", err)
	}
	redisClient := redis.NewRedisClient(cfg)

	// Repositories & Services
	todoRepo := mysql.NewTodoRepository(*db)
	outboxRepo := mysql.NewOutboxRepository()
	streamPublisher := redis.NewStreamPublisher(*redisClient, cfg)
	todoService := usecase.NewTodoService(todoRepo, outboxRepo, streamPublisher)

	// GraphQL Resolver
	resolver := &graph.Resolver{
		TodoService: *todoService,
	}

	// gqlgen server
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.Use(extension.Introspection{})

	// Outbox Worker
	outboxWorker := worker.NewOutboxWorker(outboxRepo, streamPublisher)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go outboxWorker.Start(ctx)

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		cancel()
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	// Run Hitrix server
	port, _ := strconv.Atoi(cfg.ServerPort)
	s.RunServer(uint(port), graph.NewExecutableSchema(graph.Config{Resolvers: resolver}),
		func(ginEngine *gin.Engine) {
			// GraphQL routes
			ginEngine.GET("/", func(c *gin.Context) {
				playground.Handler("GraphQL playground", "/query").ServeHTTP(c.Writer, c.Request)
			})
			ginEngine.POST("/query", func(c *gin.Context) {
				srv.ServeHTTP(c.Writer, c.Request)
			})
		}, nil,
	)
}
