package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"
	"github.com/hibiken/asynq"
)

type SyncTaskPayload struct {
	TenantID uint64 `json:"tenant_id"`
	StoreID  uint64 `json:"store_id"`
	Type     string `json:"type"`
	Mode     string `json:"mode"`
}

func main() {
	config.Load()
	database.Connect()
	database.ConnectRedis()
	model.Migrate()

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.AppCfg.Redis.URL})
	defer client.Close()

	payloadBytes, _ := json.Marshal(SyncTaskPayload{TenantID: 1, StoreID: 1, Type: "products", Mode: "full"})
	task := asynq.NewTask("sync:products", payloadBytes)
	info, err := client.Enqueue(task, asynq.Queue("default"), asynq.MaxRetry(5))
	if err != nil {
		log.Fatalf("failed to enqueue task: %v", err)
	}
	fmt.Printf("enqueued task: %s\n", info.ID)

	<-time.After(5 * time.Second)
	var job model.SyncJob
	if err := database.DB.Where("id = ?", 1).First(&job).Error; err == nil {
		fmt.Printf("job status: %s\n", job.Status)
	}
}
