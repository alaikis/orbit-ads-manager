package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tenantHandlers "orbit/apps/api/internal/tenant/handlers"
	adHandlers "orbit/apps/api/internal/adplatform/handlers"
	reportHandlers "orbit/apps/api/internal/report/handlers"
	agentHandlers "orbit/apps/api/internal/agent/handlers"
	ruleHandlers "orbit/apps/api/internal/rule/handlers"
	feedHandlers "orbit/apps/api/internal/feed/handlers"
	notifHandlers "orbit/apps/api/internal/notify/handlers"
	auditHandlers "orbit/apps/api/internal/audit/handlers"
	oauthHandlers "orbit/apps/api/internal/oauth/handlers"
	intelHandlers "orbit/apps/api/internal/intelligence/handlers"
	workspaceHandlers "orbit/apps/api/internal/workspace"
	providerHandlers "orbit/apps/api/internal/provider"
	"orbit/apps/api/internal/auth"
	worker "orbit/apps/api/internal/sync"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/config"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/robfig/cron/v3"
)

func main() {
	config.Load()
	database.Connect()
	database.ConnectRedis()
	model.Migrate()

	// Shared shutdown context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start scheduler in background
	schedulerDone := make(chan struct{})
	go startScheduler(ctx, schedulerDone)

	// Start worker in background
	workerDone := make(chan struct{})
	go startWorker(ctx, workerDone)

	// Start API server (blocking)
	startServer(cancel)

	// Wait for background components to finish
	<-schedulerDone
	<-workerDone

	log.Println("all components stopped")
}

func startServer(cancel context.CancelFunc) {
	r := gin.Default()
	r.Use(httputil.CORSMiddleware())
	r.Use(httputil.Recovery())

	api := r.Group("/api/v1")
	{
		auth.RegisterRoutes(api)
		tenantHandlers.RegisterStoreRoutes(api)
		tenantHandlers.RegisterProductRoutes(api)
		adHandlers.RegisterAdAccountRoutes(api)
		adHandlers.RegisterCampaignRoutes(api)
		reportHandlers.RegisterMetricsRoutes(api)
		agentHandlers.RegisterAgentRoutes(api)
		ruleHandlers.RegisterRuleRoutes(api)
		feedHandlers.RegisterFeedRoutes(api)
		reportHandlers.RegisterReportRoutes(api)
		notifHandlers.RegisterNotificationRoutes(api)
		auditHandlers.RegisterAuditRoutes(api)
		oauthHandlers.RegisterOAuthRoutes(api)
		intelHandlers.RegisterIntelligenceRoutes(api)
		workspaceHandlers.RegisterWorkspaceRoutes(api)
		providerHandlers.RegisterProviderRoutes(api)
		httputil.RegisterHealthRoutes(api)
	}

	port := config.AppCfg.App.Port
	if port == "" {
		port = "8080"
	}
	fmt.Printf("server starting on :%s\n", port)

	// Graceful shutdown for HTTP server
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		fmt.Println("shutting down api server...")
		cancel()
	}()

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func startWorker(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.AppCfg.Redis.Addr()},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 5,
				"default":  10,
				"low":      10,
			},
		},
	)

	mux := asynq.NewServeMux()
	worker.RegisterHandlers(mux)

	go func() {
		<-ctx.Done()
		fmt.Println("shutting down worker...")
		srv.Shutdown()
	}()

	if err := srv.Start(mux); err != nil {
		log.Printf("worker stopped: %v", err)
	}
}

func startScheduler(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	c := cron.New()

	c.AddFunc("0 */3 * * *", func() { fmt.Println("scheduled: ad_meta sync") })
	c.AddFunc("*/30 * * * *", func() { fmt.Println("scheduled: product/order sync") })
	c.AddFunc("0 */4 * * *", func() { fmt.Println("scheduled: feed regenerate") })
	c.AddFunc("*/10 * * * *", func() { fmt.Println("scheduled: token refresh scan") })
	c.AddFunc("0 9 * * *", func() { fmt.Println("scheduled: daily rule evaluation") })
	c.AddFunc("0 8 * * *", func() { fmt.Println("scheduled: daily ad metrics sync") })

	c.Start()
	fmt.Println("scheduler started")

	go func() {
		<-ctx.Done()
		fmt.Println("shutting down scheduler...")
		c.Stop()
	}()

	<-ctx.Done()
	// Give scheduler a moment to finish running jobs
	time.Sleep(500 * time.Millisecond)
}
