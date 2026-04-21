package server

import (
	"context"
	"strconv"
	"strings"

	cv1 "review-c/api/consumer/v1"
	"review-c/internal/conf"
	"review-c/internal/server/ctxkey"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// Header 常量（约定与 API-Gateway 对齐）。
const (
	HeaderUserID    = "X-User-Id"
	HeaderRequestID = "X-Request-Id"
)

// headerExtract 从 HTTP 请求头读 X-User-Id / X-Request-Id，写入 context。
//   - 非 HTTP 入口（比如 gRPC）直接透传；
//   - X-User-Id 非数字或空均视为"未提供"（后续 userIDGuard 才决定是否报错）；
//   - 本中间件只写 context，不做拒绝判断（鉴权由 API-Gateway 负责）。
func headerExtract() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if hrt, ok := tr.(*khttp.Transport); ok {
					h := hrt.Request().Header
					if v := h.Get(HeaderUserID); v != "" {
						if uid, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil && uid > 0 {
							ctx = ctxkey.WithUserID(ctx, uid)
						}
					}
					if v := h.Get(HeaderRequestID); v != "" {
						ctx = ctxkey.WithRequestID(ctx, v)
					}
				}
			}
			return next(ctx, req)
		}
	}
}

// userIDGuard 越权防护：对涉及 user_id 的 RPC，强制 body/path 的 user_id == context.user_id（来自 X-User-Id）。
// 配置 security.enforce_user_id_consistency=false 时整条中间件短路（留给可信网关环境）。
//
// 适用 RPC：
//   - CreateReview  (body.user_id)
//   - ListReviewByUserID  (path.user_id，Kratos HTTP 会把 path 参数塞进 request.UserId)
// 其他 RPC 不包含 user_id，跳过检查。
func userIDGuard(sec *conf.Security) middleware.Middleware {
	enforce := true
	if sec != nil {
		enforce = sec.EnforceUserIdConsistency
	}
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if !enforce {
				return next(ctx, req)
			}
			bodyUID, ok := extractUserID(req)
			if !ok {
				// 请求不涉及 user_id，放行
				return next(ctx, req)
			}
			ctxUID := ctxkey.UserIDFromContext(ctx)
			if ctxUID == 0 {
				return nil, kerrors.Unauthorized("USER_ID_MISSING", "missing or invalid X-User-Id header")
			}
			if bodyUID != ctxUID {
				return nil, kerrors.Forbidden("USER_ID_MISMATCH", "user_id in request does not match X-User-Id header")
			}
			return next(ctx, req)
		}
	}
}

// extractUserID 对涉及 user_id 的 request 类型做 type switch；
// 其他 request 返回 (0, false) 表示本 RPC 无需守卫。
// 新增含 user_id 的 RPC 时，只需在此处追加一个 case。
func extractUserID(req interface{}) (int64, bool) {
	switch r := req.(type) {
	case *cv1.CreateReviewRequest:
		return r.GetUserId(), true
	case *cv1.ListReviewByUserIDRequest:
		return r.GetUserId(), true
	default:
		return 0, false
	}
}
