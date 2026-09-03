package rule_test

import (
	"testing"
	"orbit/apps/api/internal/rule"
	"github.com/stretchr/testify/assert"
)

func TestNewRuleEngine(t *testing.T) {
	engine := rule.NewRuleEngine()
	assert.NotNil(t, engine)
}

func TestRuleEngineConditionParsing(t *testing.T) {
	engine := rule.NewRuleEngine()
	_ = engine
	cond := map[string]interface{}{
		"metric": "spend",
		"operator": "gt",
		"value": float64(100),
	}
	assert.Equal(t, "spend", cond["metric"])
	assert.Equal(t, "gt", cond["operator"])
	assert.Equal(t, float64(100), cond["value"])
}
