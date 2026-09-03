package feed

import (
	"context"
	"fmt"
	"strings"
	"time"

	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
)

type FeedService struct{}

func NewFeedService() *FeedService {
	return &FeedService{}
}

func (s *FeedService) Generate(ctx context.Context, feed *model.Feed) (*model.FeedRun, error) {
	if feed == nil {
		return nil, fmt.Errorf("feed is nil")
	}

	run := &model.FeedRun{
		FeedID:    feed.ID,
		Status:    "running",
		StartedAt: *getNow(),
	}
	if err := database.DB.Create(run).Error; err != nil {
		return nil, err
	}

	var products []model.Product
	if err := database.DB.Where("store_id = ? AND status = 'active' AND deleted_at IS NULL", feed.StoreID).Find(&products).Error; err != nil {
		run.Status = "failed"
		run.FinishedAt = getNow()
		database.DB.Save(run)
		return run, err
	}

	var sb strings.Builder
	switch feed.Format {
	case "google_shopping_xml":
		sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?><rss version="2.0" xmlns:g="http://base.google.com/ns/1.0"><channel>`)
		for _, p := range products {
			sb.WriteString(fmt.Sprintf(`<item><g:id>%d</g:id><g:title>%s</g:title><g:link>%s</g:link><g:price>%.2f</g:price><g:availability>in stock</g:availability></item>`,
				p.ID, escapeXML(p.Title), p.Link, float64(p.PriceCents)/100))
		}
		sb.WriteString(`</channel></rss>`)
	case "meta_catalog_csv":
		sb.WriteString("id,title,description,availability,condition,price,link,image_link\n")
		for _, p := range products {
			sb.WriteString(fmt.Sprintf("%d,\"%s\",\"%s\",in stock,new,%.2f,%s,%s\n",
				p.ID, escapeCSV(p.Title), escapeCSV(p.Description), float64(p.PriceCents)/100, p.Link, p.ImageURL))
		}
	case "google_shopping_tsv":
		sb.WriteString("id\ttitle\tdescription\tlink\timage_link\tprice\tavailability\tcondition\n")
		for _, p := range products {
			sb.WriteString(fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%.2f\tin stock\tnew\n",
				p.ID, p.Title, p.Description, p.Link, p.ImageURL, float64(p.PriceCents)/100))
		}
	default:
		sb.WriteString("id\ttitle\tprice\n")
		for _, p := range products {
			sb.WriteString(fmt.Sprintf("%d\t%s\t%.2f\n", p.ID, p.Title, float64(p.PriceCents)/100))
		}
	}

	run.Status = "success"
	run.Rows = len(products)
	now := getNow()
	run.FinishedAt = now
	database.DB.Save(run)

	feed.LastGeneratedAt = now
	feed.GeneratedRows = len(products)
	feed.ErrorCount = 0
	database.DB.Save(feed)

	return run, nil
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&")
	s = strings.ReplaceAll(s, "<", "<")
	s = strings.ReplaceAll(s, ">", ">")
	s = strings.ReplaceAll(s, `"`, "\"")
	return s
}

func escapeCSV(s string) string {
	s = strings.ReplaceAll(s, `"`, `""`)
	if strings.Contains(s, ",") || strings.Contains(s, "\n") {
		return fmt.Sprintf(`"%s"`, s)
	}
	return s
}

func getNow() *time.Time {
	now := time.Now().UTC()
	return &now
}