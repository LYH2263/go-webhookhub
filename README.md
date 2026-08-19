# go-webhookhub

Go 实现的 Webhook 投递中心：订阅 HTTP 端点、按事件名过滤、HMAC-SHA256 签名 POST、失败指数退避重试、内存投递日志（可选落盘）。配套 `webhookd` 管理服务与静态页。

## 运行

```bash
go test ./... -count=1
go run ./cmd/webhookd -addr :8090 -web web
```

浏览器打开 `http://localhost:8090/`：注册端点、触发 Dispatch、查看最近投递。

## 库面

```go
h := webhookhub.New()
id, err := h.Subscribe(webhookhub.Endpoint{
    URL:     "https://example.com/hook",
    Secret:  []byte("s3cret"),
    Events:  []string{"order.*"},
    Enabled: true,
})
results, err := h.Dispatch("order.created", []byte(`{"id":1}`))
_ = h.RecentDeliveries(20)
_ = h.Close()
```

投递请求头包含 `X-Hub-Signature-256`（`sha256=` + HMAC-SHA256 hex）、`X-Hub-Event`、`X-Hub-Delivery`。
