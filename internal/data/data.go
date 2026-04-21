package data

import (
	"context"
	"time"

	"review-c/internal/conf"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	tracingmw "github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	consulapi "github.com/hashicorp/consul/api"

	rv1 "review-service/api/review/v1"
)

// ProviderSet 暴露 data 层的 wire providers。
// 顺序：Discovery → ReviewServiceClient → Data → ConsumerRepo。
var ProviderSet = wire.NewSet(NewDiscovery, NewReviewServiceClient, NewData, NewConsumerRepo)

// Data 在 BFF 中退化为「下游服务客户端集合」。
type Data struct {
	rc  rv1.ReviewClient
	log *log.Helper
}

// NewData 组装下游客户端容器。连接 cleanup 由 NewReviewServiceClient 负责。
func NewData(rc rv1.ReviewClient, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "data"))
	cleanup := func() {
		helper.Info("closing the data resources")
	}
	return &Data{rc: rc, log: helper}, cleanup, nil
}

// NewDiscovery 基于 Consul 构造服务发现；未配置 Consul 时返回 nil（允许直连 endpoint）。
func NewDiscovery(c *conf.Registry, logger log.Logger) registry.Discovery {
	l := log.NewHelper(log.With(logger, "module", "registry"))
	if c == nil || c.Consul == nil || c.Consul.Address == "" {
		l.Warn("no consul registry configured; service discovery disabled")
		return nil
	}
	cfg := consulapi.DefaultConfig()
	cfg.Address = c.Consul.Address
	if c.Consul.Scheme != "" {
		cfg.Scheme = c.Consul.Scheme
	}
	cli, err := consulapi.NewClient(cfg)
	if err != nil {
		l.Errorf("new consul client failed: %v", err)
		return nil
	}
	return consul.New(cli)
}

// NewReviewServiceClient 构造下游 review-service 的 gRPC client。
// endpoint 支持：
//   - "discovery:///review-service" —— 需 dis != nil
//   - "127.0.0.1:9000"              —— 直连
func NewReviewServiceClient(c *conf.Data, dis registry.Discovery, logger log.Logger) (rv1.ReviewClient, func(), error) {
	timeout := 800 * time.Millisecond
	endpoint := ""
	if c != nil && c.ReviewService != nil {
		endpoint = c.ReviewService.Endpoint
		if c.ReviewService.Timeout != nil {
			if d := c.ReviewService.Timeout.AsDuration(); d > 0 {
				timeout = d
			}
		}
	}

	opts := []transgrpc.ClientOption{
		transgrpc.WithEndpoint(endpoint),
		transgrpc.WithTimeout(timeout),
		transgrpc.WithMiddleware(
			recovery.Recovery(),
			tracingmw.Client(),
		),
	}
	if dis != nil {
		opts = append(opts, transgrpc.WithDiscovery(dis))
	}

	conn, err := transgrpc.DialInsecure(context.Background(), opts...)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		_ = conn.Close()
		log.NewHelper(logger).Info("closed review-service gRPC connection")
	}
	return rv1.NewReviewClient(conn), cleanup, nil
}
