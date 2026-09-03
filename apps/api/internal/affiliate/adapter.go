package affiliate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"orbit/apps/api/config"
)

type AffiliateProduct struct {
	ID          string
	Name        string
	Description string
	Link        string
	ImageURL    string
	Price       float64
	Currency    string
	Category    string
	Brand       string
}

type AffiliateTransaction struct {
	ID           string
	Status       string
	Amount       float64
	Currency     string
	Commission   float64
	ProductName  string
	PurchaseDate string
	ClickDate    string
}

type AffiliateCommission struct {
	ID          string
	Status      string
	Amount      float64
	Currency    string
	Transaction string
	PaidDate    string
}

type AffiliateNetwork interface {
	Connect(ctx context.Context, cfg PlatformConfig) error
	TestConnection(ctx context.Context, cfg PlatformConfig) error
	ListProducts(ctx context.Context, cfg PlatformConfig, limit int) ([]AffiliateProduct, error)
	ListTransactions(ctx context.Context, cfg PlatformConfig, startDate, endDate string) ([]AffiliateTransaction, error)
	ListCommissions(ctx context.Context, cfg PlatformConfig, startDate, endDate string) ([]AffiliateCommission, error)
}

type PlatformConfig struct {
	AccessToken string
	AccountID   string
	WebsiteID   string
}

type CJAdapter struct {
	config *config.Config
}

func NewCJAdapter(cfg *config.Config) *CJAdapter {
	return &CJAdapter{config: cfg}
}

func (a *CJAdapter) Connect(ctx context.Context, cfg PlatformConfig) error {
	return a.TestConnection(ctx, cfg)
}

func (a *CJAdapter) TestConnection(ctx context.Context, cfg PlatformConfig) error {
	url := "https://api.cj.com/v3/commissions"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("cj API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (a *CJAdapter) ListProducts(ctx context.Context, cfg PlatformConfig, limit int) ([]AffiliateProduct, error) {
	url := fmt.Sprintf("https://api.cj.com/v3/products?website_id=%s&limit=%d", cfg.WebsiteID, limit)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cj API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Products []struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Link        string `json:"link"`
			ImageURL    string `json:"image_url"`
			Price       string `json:"price"`
			Currency    string `json:"currency"`
			Category    string `json:"category"`
			Brand       string `json:"brand"`
		} `json:"products"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	products := make([]AffiliateProduct, 0, len(result.Products))
	for _, p := range result.Products {
		price := 0.0
		_ = price
		products = append(products, AffiliateProduct{
			ID:          p.ID,
			Name:        p.Title,
			Description: p.Description,
			Link:        p.Link,
			ImageURL:    p.ImageURL,
			Price:       0,
			Currency:    p.Currency,
			Category:    p.Category,
			Brand:       p.Brand,
		})
	}

	return products, nil
}

func (a *CJAdapter) ListTransactions(ctx context.Context, cfg PlatformConfig, startDate, endDate string) ([]AffiliateTransaction, error) {
	url := fmt.Sprintf("https://api.cj.com/v3/transactions?website_id=%s&start_date=%s&end_date=%s", cfg.WebsiteID, startDate, endDate)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cj API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Transactions []struct {
			ID           string `json:"id"`
			Status       string `json:"status"`
			Amount       string `json:"amount"`
			Currency     string `json:"currency"`
			Commission   string `json:"commission"`
			ProductName  string `json:"product_name"`
			PurchaseDate string `json:"purchase_date"`
			ClickDate    string `json:"click_date"`
		} `json:"transactions"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	transactions := make([]AffiliateTransaction, 0, len(result.Transactions))
	for _, t := range result.Transactions {
		transactions = append(transactions, AffiliateTransaction{
			ID:           t.ID,
			Status:       t.Status,
			Amount:       0,
			Currency:     t.Currency,
			Commission:   0,
			ProductName:  t.ProductName,
			PurchaseDate: t.PurchaseDate,
			ClickDate:    t.ClickDate,
		})
	}

	return transactions, nil
}

func (a *CJAdapter) ListCommissions(ctx context.Context, cfg PlatformConfig, startDate, endDate string) ([]AffiliateCommission, error) {
	url := fmt.Sprintf("https://api.cj.com/v3/commissions?website_id=%s&start_date=%s&end_date=%s", cfg.WebsiteID, startDate, endDate)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cj API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Commissions []struct {
			ID          string `json:"id"`
			Status      string `json:"status"`
			Amount      string `json:"amount"`
			Currency    string `json:"currency"`
			Transaction string `json:"transaction_id"`
			PaidDate    string `json:"paid_date"`
		} `json:"commissions"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	commissions := make([]AffiliateCommission, 0, len(result.Commissions))
	for _, c := range result.Commissions {
		commissions = append(commissions, AffiliateCommission{
			ID:          c.ID,
			Status:      c.Status,
			Amount:      0,
			Currency:    c.Currency,
			Transaction: c.Transaction,
			PaidDate:    c.PaidDate,
		})
	}

	return commissions, nil
}

type AWINAdapter struct {
	config *config.Config
}

func NewAWINAdapter(cfg *config.Config) *AWINAdapter {
	return &AWINAdapter{config: cfg}
}

func (a *AWINAdapter) Connect(ctx context.Context, cfg PlatformConfig) error {
	return a.TestConnection(ctx, cfg)
}

func (a *AWINAdapter) TestConnection(ctx context.Context, cfg PlatformConfig) error {
	url := fmt.Sprintf("https://api.awin.com/publishers/%s", cfg.AccountID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("awin API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (a *AWINAdapter) ListProducts(ctx context.Context, cfg PlatformConfig, limit int) ([]AffiliateProduct, error) {
	url := fmt.Sprintf("https://api.awin.com/publishers/%s/lists/1/products?limit=%d", cfg.AccountID, limit)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("awin API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result []struct {
		ID          string  `json:"id"`
		Name        string  `json:"product_name"`
		Description string  `json:"description"`
		Link        string  `json:"url"`
		ImageURL    string  `json:"image_url"`
		Price       float64 `json:"price"`
		Currency    string  `json:"currency"`
		Category    string  `json:"category"`
		Brand       string  `json:"brand"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	products := make([]AffiliateProduct, 0, len(result))
	for _, p := range result {
		products = append(products, AffiliateProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Link:        p.Link,
			ImageURL:    p.ImageURL,
			Price:       p.Price,
			Currency:    p.Currency,
			Category:    p.Category,
			Brand:       p.Brand,
		})
	}

	return products, nil
}

func (a *AWINAdapter) ListTransactions(ctx context.Context, cfg PlatformConfig, startDate, endDate string) ([]AffiliateTransaction, error) {
	url := fmt.Sprintf("https://api.awin.com/publishers/%s/transactions?start_date=%s&end_date=%s", cfg.AccountID, startDate, endDate)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("awin API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result []struct {
		ID           string  `json:"id"`
		Status       string  `json:"status"`
		Amount       float64 `json:"amount"`
		Currency     string  `json:"currency"`
		Commission   float64 `json:"commission_amount"`
		ProductName  string  `json:"product_name"`
		PurchaseDate string  `json:"purchase_date"`
		ClickDate    string  `json:"click_date"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	transactions := make([]AffiliateTransaction, 0, len(result))
	for _, t := range result {
		transactions = append(transactions, AffiliateTransaction{
			ID:           t.ID,
			Status:       t.Status,
			Amount:       t.Amount,
			Currency:     t.Currency,
			Commission:   t.Commission,
			ProductName:  t.ProductName,
			PurchaseDate: t.PurchaseDate,
			ClickDate:    t.ClickDate,
		})
	}

	return transactions, nil
}

func (a *AWINAdapter) ListCommissions(ctx context.Context, cfg PlatformConfig, startDate, endDate string) ([]AffiliateCommission, error) {
	url := fmt.Sprintf("https://api.awin.com/publishers/%s/commissions?start_date=%s&end_date=%s", cfg.AccountID, startDate, endDate)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("awin API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result []struct {
		ID          string  `json:"id"`
		Status      string  `json:"status"`
		Amount      float64 `json:"amount"`
		Currency    string  `json:"currency"`
		Transaction string  `json:"transaction_id"`
		PaidDate    string  `json:"paid_date"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	commissions := make([]AffiliateCommission, 0, len(result))
	for _, c := range result {
		commissions = append(commissions, AffiliateCommission{
			ID:          c.ID,
			Status:      c.Status,
			Amount:      c.Amount,
			Currency:    c.Currency,
			Transaction: c.Transaction,
			PaidDate:    c.PaidDate,
		})
	}

	return commissions, nil
}
