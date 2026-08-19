// Package webhookhub 是 Webhook 投递中心的库面。
//
// 调用方注册 Endpoint（URL、HMAC 密钥、事件过滤），随后通过 Dispatch /
// DispatchContext 按事件名扇出 HTTP POST。每次投递带 X-Hub-Signature-256
// 签名头，失败按策略指数退避重试，结果写入内存投递日志（可选落盘）。
//
// Close 之后拒绝新的订阅与投递，返回 typed ErrClosed，不会解引用已清空的客户端。
package webhookhub
