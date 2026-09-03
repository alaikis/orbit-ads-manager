package agent

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	RiskLevel   RiskLevel              `json:"risk_level"`
	Handler     string                 `json:"handler"`
}

type ToolRegistry struct {
	tools map[string]ToolDefinition
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: map[string]ToolDefinition{
			"get_metrics": {
				Name:        "get_metrics",
				Description: "Get campaign performance metrics",
				Parameters:  map[string]interface{}{"type": "object", "properties": map[string]interface{}{"scope": map[string]interface{}{"type": "string"}}},
				RiskLevel:   RiskLow,
				Handler:     "get_metrics",
			},
			"get_campaigns": {
				Name:        "get_campaigns",
				Description: "Get list of campaigns",
				Parameters:  map[string]interface{}{"type": "object", "properties": map[string]interface{}{"scope": map[string]interface{}{"type": "string"}}},
				RiskLevel:   RiskLow,
				Handler:     "get_campaigns",
			},
			"create_campaign_suggestion": {
				Name:        "create_campaign_suggestion",
				Description: "Get AI-powered campaign suggestions",
				Parameters:  map[string]interface{}{"type": "object", "properties": map[string]interface{}{"scope": map[string]interface{}{"type": "string"}}},
				RiskLevel:   RiskMedium,
				Handler:     "create_campaign_suggestion",
			},
			"generate_creative": {
				Name:        "generate_creative",
				Description: "Generate ad creatives with AI",
				Parameters:  map[string]interface{}{"type": "object", "properties": map[string]interface{}{"scope": map[string]interface{}{"type": "string"}}},
				RiskLevel:   RiskLow,
				Handler:     "generate_creative",
			},
			"optimize_bidding": {
				Name:        "optimize_bidding",
				Description: "Optimize bidding strategy",
				Parameters:  map[string]interface{}{"type": "object", "properties": map[string]interface{}{"scope": map[string]interface{}{"type": "string"}}},
				RiskLevel:   RiskMedium,
				Handler:     "optimize_bidding",
			},
		},
	}
}

func (r *ToolRegistry) Get(name string) (ToolDefinition, bool) {
	t, ok := r.tools[name]
	return t, ok
}

func (r *ToolRegistry) List() []ToolDefinition {
	var out []ToolDefinition
	for _, t := range r.tools {
		out = append(out, t)
	}
	return out
}

func (r *ToolRegistry) Register(def ToolDefinition) {
	r.tools[def.Name] = def
}
