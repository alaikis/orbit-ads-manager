package adplatform_test

import (
	"testing"
	"orbit/apps/api/internal/adplatform"
	"github.com/stretchr/testify/assert"
)

func TestNewGoogleAdapter(t *testing.T) {
	adapter := adplatform.NewGoogleAdapter(nil)
	assert.NotNil(t, adapter)
}

func TestNewMetaAdapter(t *testing.T) {
	adapter := adplatform.NewMetaAdapter(nil)
	assert.NotNil(t, adapter)
}

func TestGoogleListCampaigns(t *testing.T) {
	adapter := adplatform.NewGoogleAdapter(nil)
	_, err := adapter.ListCampaigns(nil, adplatform.PlatformConfig{})
	assert.Error(t, err)
}

func TestMetaListCampaigns(t *testing.T) {
	adapter := adplatform.NewMetaAdapter(nil)
	_, err := adapter.ListCampaigns(nil, adplatform.PlatformConfig{})
	assert.Error(t, err)
}
