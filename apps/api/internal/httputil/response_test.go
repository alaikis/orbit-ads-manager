package httputil_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"orbit/apps/api/internal/httputil"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	httputil.Success(c, map[string]string{"hello": "world"})
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, map[string]interface{}(map[string]interface{}{"hello": "world"}), resp["data"])
}

func TestErrorResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name       string
		status     int
		code       int
		message    string
		fn         func(*gin.Context, string)
	}{
		{"bad_request", http.StatusBadRequest, 1000, "bad", func(c *gin.Context, m string) { httputil.BadRequest(c, m, nil) }},
		{"unauthorized", http.StatusUnauthorized, 1001, "unauth", httputil.Unauthorized},
		{"forbidden", http.StatusForbidden, 1002, "forbid", httputil.Forbidden},
		{"not_found", http.StatusNotFound, 1003, "not found", httputil.NotFound},
		{"conflict", http.StatusConflict, 1004, "conflict", httputil.Conflict},
		{"too_many", http.StatusTooManyRequests, 1005, "rate limit", httputil.TooManyRequests},
		{"internal", http.StatusInternalServerError, 1006, "error", httputil.InternalError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/", nil)
			tt.fn(c, tt.message)
			assert.Equal(t, tt.status, w.Code)
			var resp map[string]interface{}
			json.NewDecoder(w.Body).Decode(&resp)
			assert.Equal(t, float64(tt.code), resp["code"])
		})
	}
}
