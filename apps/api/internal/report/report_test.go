package report_test

import (
	"context"
	"testing"
	"orbit/apps/api/internal/report"
	"github.com/stretchr/testify/assert"
)

func TestNewReportService(t *testing.T) {
	svc := report.NewReportService()
	assert.NotNil(t, svc)
}

func TestGetSummary(t *testing.T) {
	svc := report.NewReportService()
	_, err := svc.GetSummary(context.Background(), 1, "tenant", nil, "7d")
	assert.Error(t, err)
}

func TestExportCSV(t *testing.T) {
	svc := report.NewReportService()
	_, err := svc.ExportCSV(context.Background(), 1, "tenant", nil)
	assert.Error(t, err)
}
