package server

import (
	cv1 "github.com/huicod/reviewapis/consumer/v1"
	"review-c/internal/conf"
	"review-c/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer 构造 review-c 的 HTTP 服务器。
// 中间件链（Kratos 链式语义：上游到下游依序执行）：
//   recovery    —— 兜底 panic，避免服务崩溃
//   tracing     —— OTel Server span
//   logging     —— 结构化请求日志（不打 body，避免 PII 泄露）
//   validate    —— 执行 protoc-gen-validate 规则；非法入参直接返 400 INVALID_ARGUMENT
//   headerExtract —— 从 X-User-Id / X-Request-Id 注入 context
//   userIDGuard —— C 端越权防护（body/path.user_id == X-User-Id，可配置关闭）
func NewHTTPServer(c *conf.Server, sec *conf.Security, consumer *service.ConsumerService, logger log.Logger) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			validate.Validator(),
			headerExtract(),
			userIDGuard(sec),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	cv1.RegisterConsumerHTTPServer(srv, consumer)
	return srv
}
