package intelligence

import (
	"time"
)

type AudienceSegment struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Size        int                    `json:"size"`
	Confidence   float64               `json:"confidence"`
	Features    map[string]interface{} `json:"features"`
	Platforms   []string               `json:"platforms"`
	CreatedAt   time.Time              `json:"created_at"`
}

type ProductInsight struct {
	ProductID    uint64  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Margin       float64 `json:"margin"`
	Popularity   int     `json:"popularity"`
	Seasonality  string  `json:"seasonality"`
	BestPlatform string  `json:"best_platform"`
	BestAudience string  `json:"best_audience"`
	RecommendedBid float64 `json:"recommended_bid"`
}

type CampaignStrategy struct {
	Platform      string    `json:"platform"`
	Objective     string    `json:"objective"`
	BudgetSplit   float64   `json:"budget_split"`
	Audience      []string  `json:"audience"`
	BiddingType   string    `json:"bidding_type"`
	BidAmount     float64   `json:"bid_amount"`
	Creatives     []CreativeSuggestion `json:"creatives"`
	Schedule      []TimeSlot `json:"schedule"`
	Reasoning     string    `json:"reasoning"`
}

type CreativeSuggestion struct {
	Type        string `json:"type"`
	Headline    string `json:"headline"`
	Description string `json:"description"`
	ImagePrompt string `json:"image_prompt"`
	CTA         string `json:"cta"`
}

type GeneratedCreative struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Headline    string `json:"headline"`
	Description string `json:"description"`
	CTA         string `json:"cta"`
	ImageURL    string `json:"image_url,omitempty"`
	VideoURL    string `json:"video_url,omitempty"`
	Prompt      string `json:"prompt,omitempty"`
	Platform    string `json:"platform"`
	Score       float64 `json:"score"`
}

type TimeSlot struct {
	DayOfWeek int `json:"day_of_week"`
	HourStart int `json:"hour_start"`
	HourEnd   int `json:"hour_end"`
	BidAdjust float64 `json:"bid_adjust"`
}

type BiddingRecommendation struct {
	CampaignID  uint64    `json:"campaign_id"`
	CurrentBid  float64   `json:"current_bid"`
	Recommended float64   `json:"recommended"`
	Min         float64   `json:"min"`
	Max         float64   `json:"max"`
	Reason      string    `json:"reason"`
	Confidence  float64   `json:"confidence"`
	UpdatedAt   time.Time `json:"updated_at"`
}
