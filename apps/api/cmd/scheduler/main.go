package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"
	"github.com/robfig/cron/v3"
)

func main() {
	config.Load()
	database.Connect()
	database.ConnectRedis()
	model.Migrate()

	c := cron.New()

	// Every 3 hours: ad meta sync
	c.AddFunc("0 */3 * * *", func() { fmt.Println("scheduled: ad_meta sync") })

	// Every 30 minutes: product/order sync
	c.AddFunc("*/30 * * * *", func() { fmt.Println("scheduled: product/order sync") })

	// Every 4 hours: feed regenerate
	c.AddFunc("0 */4 * * *", func() { fmt.Println("scheduled: feed regenerate") })

	// Every 10 minutes: token refresh scan
	c.AddFunc("*/10 * * * *", func() { fmt.Println("scheduled: token refresh scan") })

	// Daily 09:00: rule evaluation
	c.AddFunc("0 9 * * *", func() { fmt.Println("scheduled: daily rule evaluation") })

	// Daily 08:00: ad metrics sync
	c.AddFunc("0 8 * * *", func() { fmt.Println("scheduled: daily ad metrics sync") })

	c.Start()
	fmt.Println("scheduler started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	c.Stop()
	fmt.Println("scheduler shutdown complete")
}
