// Package ctxkey 存放 server 层中间件注入、service 层读取的 context key。
// 单独抽子包是为了：
//  1. 避免 server 与 service 互相 import 触发循环依赖；
//  2. 把"可信头 ↔ key"映射集中在一处，新增字段修改面最小。
package ctxkey

import "context"

type userIDKey struct{}
type requestIDKey struct{}

// WithUserID 把 X-User-Id header 值放进 context，供 user-id-guard 中间件与 service 层读取。
func WithUserID(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, userIDKey{}, v)
}

// UserIDFromContext 返回中间件注入的 user_id；未设置时返回 0。
func UserIDFromContext(ctx context.Context) int64 {
	if v, ok := ctx.Value(userIDKey{}).(int64); ok {
		return v
	}
	return 0
}

// WithRequestID 把 X-Request-Id header 值放进 context，供日志与下游透传。
func WithRequestID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, v)
}

// RequestIDFromContext 返回 request id，未设置时返回空字符串。
func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey{}).(string); ok {
		return v
	}
	return ""
}
