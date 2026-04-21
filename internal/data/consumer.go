package data

import (
	"context"

	rv1 "github.com/huicod/reviewapis/review/v1"

	"review-c/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// consumerRepo 实现 biz.ConsumerRepo。
// 职责：
//  1. 把 biz-internal 的 Param 映射成下游 review-service 的 *rv1.XxxRequest；
//  2. 对读路径自动注入 caller_role=CALLER_ROLE_C（review-c 固定语义）；
//  3. ListReviews 在 data 层强制 Status=APPROVED —— C 端永远只看已审通过的评价。
type consumerRepo struct {
	data *Data
	log  *log.Helper
}

// NewConsumerRepo 由 wire 注入到 biz.ConsumerUsecase。
func NewConsumerRepo(d *Data, logger log.Logger) biz.ConsumerRepo {
	return &consumerRepo{
		data: d,
		log:  log.NewHelper(log.With(logger, "module", "data/consumer")),
	}
}

// ---------------- 写路径 ----------------

// CreateReview 注意：review-service.CreateReviewRequest 没有 caller_role 字段（创建时不需要 role 过滤）。
func (r *consumerRepo) CreateReview(ctx context.Context, p *biz.CreateReviewParam) (*rv1.CreateReviewReply, error) {
	return r.data.rc.CreateReview(ctx, &rv1.CreateReviewRequest{
		UserId:            p.UserID,
		OrderId:           p.OrderID,
		StoreId:           p.StoreID,
		SkuId:             p.SkuID,
		SpuId:             p.SpuID,
		Score:             p.Score,
		ServiceScore:      p.ServiceScore,
		ExpressScore:      p.ExpressScore,
		Content:           p.Content,
		PicUrls:           append([]string(nil), p.PicURLs...),
		VideoUrls:         append([]string(nil), p.VideoURLs...),
		Anonymous:         p.Anonymous,
		TagsJson:          p.TagsJSON,
		GoodsSnapshotJson: p.GoodsSnapshotJSON,
	})
}

// ---------------- 读路径 ----------------

func (r *consumerRepo) GetReview(ctx context.Context, reviewID int64) (*rv1.GetReviewReply, error) {
	return r.data.rc.GetReview(ctx, &rv1.GetReviewRequest{
		ReviewId:   reviewID,
		CallerRole: rv1.CallerRole_CALLER_ROLE_C,
	})
}

// ListReviewByUserID 下游 proto 没有 caller_role 字段（"我的评价"隐式等价 role=C）。
func (r *consumerRepo) ListReviewByUserID(ctx context.Context, p *biz.ListReviewByUserIDParam) (*rv1.ListReviewByUserIDReply, error) {
	return r.data.rc.ListReviewByUserID(ctx, &rv1.ListReviewByUserIDRequest{
		UserId:   p.UserID,
		Page:     p.Page,
		PageSize: p.PageSize,
	})
}

// ListReviews C 端固定 status=APPROVED（即使下游 role=C 会强制，这里显式传也不伤 —— 更早失败更好排查）。
func (r *consumerRepo) ListReviews(ctx context.Context, p *biz.ListReviewsParam) (*rv1.ListReviewsReply, error) {
	return r.data.rc.ListReviews(ctx, &rv1.ListReviewsRequest{
		StoreId:    p.StoreID,
		SpuId:      p.SpuID,
		SkuId:      p.SkuID,
		Status:     rv1.ReviewStatus_REVIEW_STATUS_APPROVED,
		Score:      p.Score,
		HasMedia:   p.HasMedia,
		Sort:       rv1.ReviewSort(p.Sort),
		CallerRole: rv1.CallerRole_CALLER_ROLE_C,
		Page:       p.Page,
		PageSize:   p.PageSize,
	})
}
