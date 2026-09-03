package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"
	"orbit/apps/api/internal/sync"
	"github.com/hibiken/asynq"
)

func main() {
	config.Load()
	database.Connect()
	database.ConnectRedis()
	model.Migrate()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.AppCfg.Redis.URL},
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

	if err := srv.Start(mux); err != nil {
		log.Fatalf("worker failed to start: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	srv.Shutdown()
	fmt.Println("worker shutdown complete")
}
