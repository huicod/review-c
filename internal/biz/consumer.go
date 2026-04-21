package biz

import (
	"context"

	rv1 "github.com/huicod/reviewapis/review/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// biz 层约定（与 review-o 对齐）：
//  - 入参使用业务中立的 Param struct，不依赖任何 proto 类型；
//  - 读路径出于「BFF 不做数据加工」的现实考虑，Reply 类型借用下游 proto（*rv1.XxxReply）；
//  - 写路径入参严格走 Param，避免下游 proto 改字段连锁改 biz/service。

// ---------- Params（业务中立） ----------

// CreateReviewParam C 端创建评价的业务入参。user_id 由上层 HTTP 中间件保证已与 X-User-Id 一致。
type CreateReviewParam struct {
	UserID            int64
	OrderID           int64
	StoreID           int64
	SkuID             int64
	SpuID             int64
	Score             int32
	ServiceScore      int32
	ExpressScore      int32
	Content           string
	PicURLs           []string
	VideoURLs         []string
	Anonymous         bool
	TagsJSON          string
	GoodsSnapshotJSON string
}

// ListReviewsParam C 端条件列表参数（仅暴露 C 端允许的筛选项）。
// Sort/Status 用 int32 承载，避免 biz 依赖 proto enum 常量；status 由 data 层固定传 APPROVED。
type ListReviewsParam struct {
	StoreID  int64
	SpuID    int64
	SkuID    int64
	Score    int32
	HasMedia bool
	Sort     int32
	Page     int32
	PageSize int32
}

// ListReviewByUserIDParam C 端"我的评价"参数。
type ListReviewByUserIDParam struct {
	UserID   int64
	Page     int32
	PageSize int32
}

// ---------- Repo 接口（由 data 层实现） ----------

// ConsumerRepo 是 biz 对 data 层的契约。
// 所有读路径返回下游 review-service 的 proto（*rv1.XxxReply）；
// 写路径同样返回下游 proto（带 review_id 等回显字段）。
type ConsumerRepo interface {
	CreateReview(ctx context.Context, p *CreateReviewParam) (*rv1.CreateReviewReply, error)
	GetReview(ctx context.Context, reviewID int64) (*rv1.GetReviewReply, error)
	ListReviewByUserID(ctx context.Context, p *ListReviewByUserIDParam) (*rv1.ListReviewByUserIDReply, error)
	ListReviews(ctx context.Context, p *ListReviewsParam) (*rv1.ListReviewsReply, error)
}

// ---------- Usecase ----------

// ConsumerUsecase 是 review-c 的业务编排层。当前 4 个方法均为透传；
// 保留 biz 层是为后续加入 C 端特有策略（如发评频控、草稿合并）时 service 层无需改动。
type ConsumerUsecase struct {
	repo ConsumerRepo
	log  *log.Helper
}

func NewConsumerUsecase(repo ConsumerRepo, logger log.Logger) *ConsumerUsecase {
	return &ConsumerUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "biz/consumer")),
	}
}

func (uc *ConsumerUsecase) CreateReview(ctx context.Context, p *CreateReviewParam) (*rv1.CreateReviewReply, error) {
	uc.log.WithContext(ctx).Infof("CreateReview param: user_id=%d order_id=%d score=%d", p.UserID, p.OrderID, p.Score)
	return uc.repo.CreateReview(ctx, p)
}

func (uc *ConsumerUsecase) GetReview(ctx context.Context, reviewID int64) (*rv1.GetReviewReply, error) {
	return uc.repo.GetReview(ctx, reviewID)
}

func (uc *ConsumerUsecase) ListReviewByUserID(ctx context.Context, p *ListReviewByUserIDParam) (*rv1.ListReviewByUserIDReply, error) {
	return uc.repo.ListReviewByUserID(ctx, p)
}

func (uc *ConsumerUsecase) ListReviews(ctx context.Context, p *ListReviewsParam) (*rv1.ListReviewsReply, error) {
	return uc.repo.ListReviews(ctx, p)
}
