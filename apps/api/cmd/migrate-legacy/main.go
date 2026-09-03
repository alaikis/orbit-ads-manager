package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"orbit/apps/api/config"
	"orbit/apps/api/internal/connection"
	"orbit/apps/api/pkg/crypto"
	"orbit/apps/api/pkg/database"
)

func main() {
	loadEnv()
	config.Load()
	database.Connect()

	type legacy struct {
		ID          uint64
		TenantID    *uint64 `gorm:"column:tenant_id"`
		WorkspaceID *uint64
		Type        string
		Name        string
		Status      string
		Config      string `gorm:"type:jsonb"`
	}
	var rows []legacy
	if err := database.DB.
		Table("providers").
		Select("id, tenant_id, workspace_id, type, name, status, config::text as config").
		Where("type IN ?", []string{"llm", "smtp"}).
		Scan(&rows).Error; err != nil {
		log.Fatalf("query providers: %v", err)
	}
	log.Printf("found %d legacy providers (llm/smtp)", len(rows))

	migrated, skipped := 0, 0
	for _, p := range rows {
		var existing int64
		if err := database.DB.Model(&connection.Connection{}).
			Where("tenant_id = ? AND name = ? AND platform = ?", derefU64(p.TenantID), p.Name, p.Type).
			Count(&existing).Error; err != nil {
			log.Printf("check existing: %v", err)
			continue
		}
		if existing > 0 {
			skipped++
			fmt.Printf("skip: connection already exists for provider %d\n", p.ID)
			continue
		}

		conn := connection.Connection{
			TenantID: derefU64(p.TenantID),
			Type:     "api_key",
			Platform: p.Type,
			Name:     p.Name,
			Status:   p.Status,
		}
		if p.WorkspaceID != nil {
			conn.WorkspaceID = p.WorkspaceID
		}
		if err := database.DB.Create(&conn).Error; err != nil {
			log.Printf("create connection for provider %d: %v", p.ID, err)
			continue
		}

		var fields map[string]string
		_ = json.Unmarshal([]byte(p.Config), &fields)
		for k, v := range fields {
			if v == "" {
				continue
			}
			ct, nonce, err := crypto.Encrypt([]byte(v))
			if err != nil {
				log.Printf("encrypt %s: %v", k, err)
				continue
			}
			cred := connection.ConnectionCredential{
				ConnID:         conn.ID,
				Key:            k,
				EncryptedValue: ct,
				Nonce:          nonce,
				Last4:          crypto.Last4(v),
			}
			if err := database.DB.Create(&cred).Error; err != nil {
				log.Printf("create credential %s: %v", k, err)
			}
		}

		migrated++
		fmt.Printf("migrated provider %d (%s/%s) → connection %d\n", p.ID, p.Type, p.Name, conn.ID)
	}

	fmt.Printf("\nDone: %d migrated, %d skipped\n", migrated, skipped)
}

func loadEnv() {
	for _, path := range []string{".env", "/apphub/ads.alakis.com/.env", "/apphub/ads.alakis.com/.env.local"} {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			eq := strings.IndexByte(line, '=')
			if eq <= 0 {
				continue
			}
			k := strings.TrimSpace(line[:eq])
			v := strings.TrimSpace(line[eq+1:])
			os.Setenv(k, v)
		}
		fmt.Printf("loaded env from %s\n", path)
		return
	}
}

func derefU64(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}
