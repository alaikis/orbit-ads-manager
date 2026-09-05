package connection

import (
	"context"
	"fmt"
	"sync"
)

// PlatformClient 统一接口 — 所有平台必须实现。
// 前端通过 PlatformSchema.Fields 动态渲染表单，后端通过 Clients 多态分发。
type PlatformClient interface {
	Test(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error)
	Sync(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, action string) (map[string]interface{}, error)
	RefreshToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error
}

var (
	clientsOnce sync.Once
	clients     = map[string]PlatformClient{}
	clientNames = []string{}
)

// Register 注册一个 platform client。通常在各平台 client 文件末尾的 init() 中调用。
// 线程安全，幂等（重复注册同一 platform key 会被忽略）。
func Register(platform string, c PlatformClient) {
	clientsOnce.Do(func() {})
	if c == nil {
		panic(fmt.Sprintf("connection: Register client is nil for %q", platform))
	}
	if _, ok := clients[platform]; ok {
		panic(fmt.Sprintf("connection: Register called twice for %q", platform))
	}
	clients[platform] = c
	clientNames = append(clientNames, platform)
}

// Clients 返回所有已注册的 platform client，按注册顺序。
// 用于前端 GET /platforms 和测试。
func GetClients() map[string]PlatformClient {
	return clients
}

// GetClientNames 返回已注册的 platform key 列表。
func GetClientNames() []string {
	return append([]string(nil), clientNames...)
}

// MustGetClient 返回指定 platform 的 client，未注册时 panic。
// 用于内部已知 platform 的场景（如 cron refresh）。
func MustGetClient(platform string) PlatformClient {
	c, ok := clients[platform]
	if !ok {
		panic(fmt.Sprintf("connection: client not registered for %q (registered: %v)", platform, clientNames))
	}
	return c
}

// GetClient 返回指定 platform 的 client，未注册时返回 (nil, false)。
func GetClient(platform string) (PlatformClient, bool) {
	c, ok := clients[platform]
	return c, ok
}
