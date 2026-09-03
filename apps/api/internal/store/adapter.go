package store

import (
	"context"
	"strconv"
)

type Product struct {
	ID          string
	Title       string
	Description string
	Link        string
	ImageURL    string
	Brand       string
	Gtin        string
	SKU         string
	PriceCents  int64
	StockQty    int
	Status      string
	Variants    []ProductVariant
}

type ProductVariant struct {
	ID           string
	SKU          string
	OptionString string
	PriceCents   int64
	InventoryQty int
}

type Order struct {
	ID            string
	OrderNumber   string
	Status        string
	TotalCents    int64
	Currency      string
	CustomerEmail string
	Items         []map[string]interface{}
	PlacedAt      string
}

type StoreInfo struct {
	ID           string
	Name         string
	Platform     string
	URL          string
	ConnectedAt  string
}

type StoreAdapter interface {
	Connect(ctx context.Context, config interface{}) (*StoreInfo, error)
	TestConnection(ctx context.Context) error
	GetStoreInfo(ctx context.Context) (*StoreInfo, error)
	SyncProducts(ctx context.Context, cursor interface{}, limit int) ([]Product, interface{}, error)
	SyncOrders(ctx context.Context, since string, cursor interface{}) ([]Order, interface{}, error)
	FetchInventory(ctx context.Context, productIDs []string) (map[string]int, error)
}

func parsePrice(s string) int64 {
	f, _ := strconv.ParseFloat(s, 64)
	return int64(f * 100)
}
