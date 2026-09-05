# Go 规则集 v1.1（增量规则）

> 适用于 `*.go` 文件审查。本规则集**增量**于 v1.0-jvm 之上，对 Go 项目启用。审查范围：本次变更涉及的所有 `.go` 文件（排除 `_test.go`，由 test 规则处理；排除 vendor/）。
>
> 优先级映射见 `common-review-guide.md`：P0→high, P1→medium, P2→low。

## 1. 架构与逻辑（arch）

### arch-go:1.1.1 错误处理路径（P0）
- **必须** 检查 `if err != nil { ... return ... }` 模式
- **禁止** 静默吞错（`if err != nil { /* ignore */ }`）
- **禁止** 仅 `log.Print` 后继续执行关键路径
- **必须** 错误向上传递时**包装**（`fmt.Errorf("ctx: %w", err)` 或 `errors.Join`）

```go
// BAD
if err != nil {
    return err  // 无上下文，调用方不知道哪一步失败
}
if err != nil {
    log.Print(err)
    // 继续执行，关键路径无错误检查
}

// GOOD
if err != nil {
    return fmt.Errorf("connection.Create tenant=%d: %w", tenantID, err)
}
```

### arch-go:1.1.2 Context 传递（P0）
- **必须** 第一个参数为 `ctx context.Context`（handler/service/dao 层）
- **禁止** 使用 `context.Background()` 在请求处理路径中（worker/cron 例外）
- **禁止** 将 ctx 存入 struct 字段
- HTTP handler **必须** 接收 ctx 并下传；DB query **必须** 用 ctx（`db.QueryContext`）

### arch-go:1.1.3 资源泄漏（P0）
- `rows.Close()` 必须在 `defer` 中或显式调用
- `resp.Body.Close()` 必须 defer
- `tx.Rollback()` / `tx.Commit()` 必须处理（`defer tx.Rollback()` + 显式 Commit）
- 文件句柄 / 网络连接同上

### arch-go:1.1.4 Goroutine 泄漏 / 并发安全（P0）
- `go func()` 启动 goroutine 必须有退出信号（ctx.Done() / channel close）
- 共享变量必须用 mutex / atomic / channel 保护
- 禁止 `time.Sleep` 在生产路径循环

### arch-go:1.2.1 接口设计（P1）
- 接口定义在消费方（"accept interfaces, return structs"）
- 避免接口污染（小接口、单方法）
- 避免循环依赖（包级别）

### arch-go:1.2.2 包组织（P1）
- `internal/` 用于私有包
- `pkg/` 用于可导出公共包
- 禁止跨层反向依赖（handler 不能 import dao 直接类型）

## 2. 安全（sec）

### sec-go:2.1.1 SQL 注入（P0）
- **禁止** 字符串拼接构造 SQL
- **必须** 用 `?` 占位符 + 参数化查询（`db.Query("... WHERE id = ?", id)`）
- **必须** 表/列名用白名单（不能用占位符的），**禁止** 用户输入拼接

```go
// BAD
db.Query("SELECT * FROM users WHERE id = " + userID)
db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))

// GOOD
db.Query("SELECT * FROM users WHERE id = ?", userID)
// 表名用 map 白名单
table, ok := tableWhitelist[name]
if !ok { return errors.New("invalid table") }
```

### sec-go:2.1.2 敏感数据泄露（P0）
- **禁止** 日志输出明文密码、token、API key、加密密钥
- **禁止** 返回错误时包含敏感信息（DB dsn 包含密码时需用 `errors.Is` 包装或剥离）
- 错误响应（HTTP body）**禁止** 包含堆栈或内部路径

### sec-go:2.1.3 加密实现（P0）
- **禁止** 使用 `crypto/md5` / `crypto/sha1` 用于安全目的
- **必须** `crypto/rand` 用于随机数（token、nonce、key）
- AES-GCM nonce **必须** 12 字节且每次唯一
- 密钥长度：AES-256 = 32 字节；RSA ≥ 2048

### sec-go:2.1.4 HTTP 安全（P0）
- TLS 配置**必须** 现代密码套件（禁用 TLS 1.0/1.1）
- Cookie **必须** `HttpOnly` + `Secure` + `SameSite=Lax/Strict`
- CORS **禁止** `Access-Control-Allow-Origin: *` 配合 credentials
- 重定向 URL **必须** 白名单校验（防 open redirect）

### sec-go:2.2.1 越权 / 租户隔离（P0）
- 所有查询**必须** 包含 `tenant_id` 过滤（多租户系统）
- 路径参数 `id` **必须** 验证属于当前 tenant
- 中间件 **必须** 注入 tenant 上下文

### sec-go:2.2.2 输入校验（P1）
- HTTP 请求参数用 validator/struct tag 校验
- 长度、范围、枚举必须显式
- 禁止信任 client 提供的 `user_id` / `tenant_id`

### sec-go:2.3.1 SSRF / URL 注入（P1）
- HTTP 客户端 URL **必须** 校验 scheme/host 白名单
- 禁止从用户输入直接拼 URL
- 内部 metadata endpoint（169.254.169.254）必须 deny

## 3. 性能（perf）

### perf-go:3.1.1 N+1 查询（P0）
- 禁止循环内单条查询（`for _, x := range items { db.Query(...) }`）
- **必须** 用 `WHERE id IN (?)` + JOIN / batch
- Eager loading 替代 lazy loading

```go
// BAD
for _, id := range ids {
    row := db.QueryRow("SELECT * FROM items WHERE id = ?", id)
}

// GOOD
db.Query("SELECT * FROM items WHERE id IN (?)", strings.Join(ids, ","))
// 或用 sqlx In() 扩展
```

### perf-go:3.1.2 缺少索引 / 全表扫描（P0）
- WHERE 条件列必须建索引
- `ORDER BY` 列无索引 → filesort（P0 for >10K rows）
- 复合索引顺序：等值在前，范围在后

### perf-go:3.1.3 内存分配（P1）
- 热点循环内避免 `[]byte` → string 转换
- 大对象 sync.Pool
- 字符串拼接在循环内用 `strings.Builder`

### perf-go:3.2.1 锁粒度（P1）
- `sync.Mutex` 保护范围**必须** 最小化
- 读多用 `sync.RWMutex`
- 跨函数持锁必须明确注释

### perf-go:3.2.2 连接池（P1）
- `db.SetMaxOpenConns` / `SetMaxIdleConns` 必须显式设置
- HTTP client `Transport.MaxIdleConnsPerHost` 必须设置
- HTTP client 必须复用（禁止每次 new）

## 4. 可维护性（maint）

### maint-go:4.1.1 命名规范（P1）
- 包名：小写、单词、避免下划线/驼峰
- 导出标识符：导出用大写开头，**必须** 有意义（禁止 `data1` `temp` `info`）
- 接口：方法名 + "er" 后缀（`Reader`、`Writer`）
- 接收者名称：1-2 字母一致（`func (s *Service)`）

### maint-go:4.1.2 函数长度（P1）
- 单函数 > 80 行需拆分
- 嵌套深度 > 4 需 early return

### maint-go:4.1.3 注释规范（P2）
- 导出标识符**必须** 有 doc 注释（`// FuncName does X`）
- TODO 格式：`// TODO(username): description`
- 禁止无意义注释（`// increment i`）

### maint-go:4.2.1 错误消息（P2）
- 错误消息小写开头（`errors.New("connection refused")` 而非 `"Connection refused"`）
- 不以标点结尾
- 包含足够上下文（哪个操作/哪个 ID）

## 5. 测试（test）

### test-go:5.1.1 测试覆盖（P0）
- 关键路径（CRUD、auth、加密、token refresh）**必须** 有 `_test.go`
- handler / service / dao 层**至少** 70% 覆盖
- 错误路径**必须** 覆盖（err 返回、边界值、空值）

### test-go:5.1.2 测试隔离（P1）
- 每个 test 用独立 DB schema 或 transaction rollback
- 禁止依赖外部服务（用 mock / interface）
- 并发测试用 `t.Parallel()` 显式标注

### test-go:5.2.1 表驱动测试（P2）
- 多 case 测试**优先** 表驱动（`tests := []struct{...}{...}`）

## 6. 风格（style）

### style-go:6.1.1 gofmt / goimports（P0）
- **必须** `gofmt -s` 通过
- **必须** `goimports` 通过
- CI 应 fail on lint

### style-go:6.1.2 显式导入分组（P2）
- 标准库 / 第三方 / 内部三组，空行分隔
- 禁止 `import .` 和 `import _` 滥用

### style-go:6.2.1 日志规范（P1）
- 用结构化日志（`slog` 或 `zap`），禁止 `fmt.Println` 在生产
- 关键事件（startup, shutdown, error）**必须** 记录
- 包含 request_id / trace_id 用于关联

### style-go:6.2.2 常量与配置（P2）
- 魔法数字提取为 const
- 配置项（timeout、retry count）从配置文件读，**禁止** 硬编码

---

## 启用方式

- bundle.json `file_extensions_seen` 包含 `.go` 时启用本规则集
- 与 v1.0-jvm 规则并存：JVM 文件用 v1.0，Go 文件用 v1.1
- 自动判定逻辑：在 `gate` 决策时按文件类型分别调用

## 已知误报抑制

| 模式 | 误报原因 | 抑制方法 |
|------|----------|----------|
| `context.Background()` | cron/worker 启动期 | 检查调用栈：若在 main 或 cron init 路径，放行 |
| `_ = err` | 显式忽略已记录日志 | 上下文包含 `log.` 调用 |
| 全局 mutex | 简单单例 | 单文件 + < 100 行 |
