package service

import (
	"context"

	cv1 "github.com/huicod/reviewapis/consumer/v1"
	"review-c/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// ConsumerService 实现 cv1.ConsumerServer（HTTP/gRPC 入口的 4 个 C 端 RPC）。
// 职责：
//  1. 把 cv1.XxxRequest 拆成 biz.XxxParam（业务中立）；
//  2. 把 biz 返回的 *rv1.XxxReply 转成 cv1.XxxReply（Mode B 固有成本，集中在 converter.go）；
//  3. 错误直接透传下游 Kratos errors（gRPC code → HTTP status 由 Kratos 自动映射）。
// user_id 越权校验在 HTTP 中间件完成，service 层到达时 body/path 的 user_id 已经过校验。
type ConsumerService struct {
	cv1.UnimplementedConsumerServer

	uc  *biz.ConsumerUsecase
	log *log.Helper
}

// NewConsumerService 由 wire 注入。
func NewConsumerService(uc *biz.ConsumerUsecase, logger log.Logger) *ConsumerService {
	return &ConsumerService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/consumer")),
	}
}

// ---------------- 写路径 ----------------

func (s *ConsumerService) CreateReview(ctx context.Context, req *cv1.CreateReviewRequest) (*cv1.CreateReviewReply, error) {
	resp, err := s.uc.CreateReview(ctx, &biz.CreateReviewParam{
		UserID:            req.GetUserId(),
		OrderID:           req.GetOrderId(),
		StoreID:           req.GetStoreId(),
		SkuID:             req.GetSkuId(),
		SpuID:             req.GetSpuId(),
		Score:             req.GetScore(),
		ServiceScore:      req.GetServiceScore(),
		ExpressScore:      req.GetExpressScore(),
		Content:           req.GetContent(),
		PicURLs:           req.GetPicUrls(),
		VideoURLs:         req.GetVideoUrls(),
		Anonymous:         req.GetAnonymous(),
		TagsJSON:          req.GetTagsJson(),
		GoodsSnapshotJSON: req.GetGoodsSnapshotJson(),
	})
	if err != nil {
		return nil, err
	}
	return &cv1.CreateReviewReply{ReviewId: resp.GetReviewId()}, nil
}

// ---------------- 读路径 ----------------

func (s *ConsumerService) GetReview(ctx context.Context, req *cv1.GetReviewRequest) (*cv1.GetReviewReply, error) {
	resp, err := s.uc.GetReview(ctx, req.GetReviewId())
	if err != nil {
		return nil, err
	}
	return &cv1.GetReviewReply{Review: reviewDetailFromRv(resp.GetReview())}, nil
}

func (s *ConsumerService) ListReviewByUserID(ctx context.Context, req *cv1.ListReviewByUserIDRequest) (*cv1.ListReviewByUserIDReply, error) {
	resp, err := s.uc.ListReviewByUserID(ctx, &biz.ListReviewByUserIDParam{
		UserID:   req.GetUserId(),
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*cv1.ReviewSummary, 0, len(resp.GetList()))
	for _, rv := range resp.GetList() {
		list = append(list, reviewSummaryFromRv(rv))
	}
	return &cv1.ListReviewByUserIDReply{
		List: list,
		Page: paginationFromRv(resp.GetPage()),
	}, nil
}

func (s *ConsumerService) ListReviews(ctx context.Context, req *cv1.ListReviewsRequest) (*cv1.ListReviewsReply, error) {
	resp, err := s.uc.ListReviews(ctx, &biz.ListReviewsParam{
		StoreID:  req.GetStoreId(),
		SpuID:    req.GetSpuId(),
		SkuID:    req.GetSkuId(),
		Score:    req.GetScore(),
		HasMedia: req.GetHasMedia(),
		Sort:     int32(req.GetSort()),
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*cv1.ReviewSummary, 0, len(resp.GetList()))
	for _, rv := range resp.GetList() {
		list = append(list, reviewSummaryFromRv(rv))
	}
	return &cv1.ListReviewsReply{
		List: list,
		Page: paginationFromRv(resp.GetPage()),
	}, nil
}
