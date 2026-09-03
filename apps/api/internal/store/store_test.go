package store_test

import (
	"testing"
	"orbit/apps/api/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestNewWooCommerceAdapter(t *testing.T) {
	adapter := store.NewWooCommerceAdapter(nil)
	assert.NotNil(t, adapter)
}

func TestNewShopifyAdapter(t *testing.T) {
	adapter := store.NewShopifyAdapter(nil)
	assert.NotNil(t, adapter)
}

func TestWooCommerceConnect(t *testing.T) {
	adapter := store.NewWooCommerceAdapter(nil)
	_, err := adapter.Connect(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

func TestShopifyConnect(t *testing.T) {
	adapter := store.NewShopifyAdapter(nil)
	_, err := adapter.Connect(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}
