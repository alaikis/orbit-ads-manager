package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Details gin.H       `json:"details,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Code: 0, Message: "created", Data: data})
}

func Accepted(c *gin.Context, data interface{}) {
	c.JSON(http.StatusAccepted, Response{Code: 0, Message: "accepted", Data: data})
}

func Error(c *gin.Context, status int, code int, message string, details gin.H) {
	if details == nil {
		details = gin.H{}
	}
	c.JSON(status, Response{Code: code, Message: message, Data: nil, Details: details})
}

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Details gin.H       `json:"details,omitempty"`
}

func BadRequest(c *gin.Context, message string, details gin.H) {
	Error(c, http.StatusBadRequest, 1000, message, details)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, 1001, message, nil)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, 1002, message, nil)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, 1003, message, nil)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, 1004, message, nil)
}

func TooManyRequests(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, 1005, message, nil)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, 1006, message, nil)
}
