package feed_test

import (
	"context"
	"testing"
	"orbit/apps/api/internal/feed"
	"github.com/stretchr/testify/assert"
)

func TestNewFeedService(t *testing.T) {
	svc := feed.NewFeedService()
	assert.NotNil(t, svc)
}

func TestGenerateMissingStore(t *testing.T) {
	svc := feed.NewFeedService()
	run, err := svc.Generate(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, run)
}
