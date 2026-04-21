package service

// dto.go —— 【读路径】rv1 (review-service proto) → cv1 (review-c proto) 的 Reply 转换。
//
// 设计说明（与 review-o 的 dto.go 同构）：
//  - 写路径（CreateReview）Reply 字段少，直接在 consumer.go 里 inline 构造；
//  - 读路径的 ReviewSummary / ReviewDetail / ReplyInfo 结构复杂，集中转换。
//  - C 端不暴露 ReviewAuditInfo（运营审核元数据对消费者无意义），所以转换函数比 review-o 少一个。
//  - 所有 FromRv 函数都 nil-safe：上游返回 nil 时返回 nil，不 panic。

import (
	cv1 "review-c/api/consumer/v1"

	rv1 "review-service/api/review/v1"
)

// ---------------- Pagination ----------------

func paginationFromRv(p *rv1.PaginationReply) *cv1.PaginationReply {
	if p == nil {
		return nil
	}
	return &cv1.PaginationReply{
		Total:    p.GetTotal(),
		Page:     p.GetPage(),
		PageSize: p.GetPageSize(),
	}
}

// ---------------- Summary / Detail / Reply ----------------

func reviewSummaryFromRv(in *rv1.ReviewSummary) *cv1.ReviewSummary {
	if in == nil {
		return nil
	}
	return &cv1.ReviewSummary{
		ReviewId:     in.GetReviewId(),
		OrderId:      in.GetOrderId(),
		UserId:       in.GetUserId(), // 匿名脱敏已在下游 review-service 的 service 层完成
		StoreId:      in.GetStoreId(),
		SkuId:        in.GetSkuId(),
		SpuId:        in.GetSpuId(),
		Score:        in.GetScore(),
		ServiceScore: in.GetServiceScore(),
		ExpressScore: in.GetExpressScore(),
		Content:      in.GetContent(),
		Anonymous:    in.GetAnonymous(),
		HasMedia:     in.GetHasMedia(),
		HasReply:     in.GetHasReply(),
		IsDefault:    in.GetIsDefault(),
		Status:       cv1.ReviewStatus(in.GetStatus()), // enum 数值一致，直接强转
		CreateTime:   in.GetCreateTime(),
		UpdateTime:   in.GetUpdateTime(),
	}
}

func replyInfoFromRv(in *rv1.ReplyInfo) *cv1.ReplyInfo {
	if in == nil {
		return nil
	}
	return &cv1.ReplyInfo{
		ReplyId:    in.GetReplyId(),
		ReviewId:   in.GetReviewId(),
		StoreId:    in.GetStoreId(),
		Content:    in.GetContent(),
		PicUrls:    append([]string(nil), in.GetPicUrls()...),
		VideoUrls:  append([]string(nil), in.GetVideoUrls()...),
		CreateTime: in.GetCreateTime(),
		UpdateTime: in.GetUpdateTime(),
		Version:    in.GetVersion(),
	}
}

// reviewDetailFromRv C 端 ReviewDetail 不含 Audit 字段 —— 运营审核元数据对消费者无意义。
func reviewDetailFromRv(in *rv1.ReviewDetail) *cv1.ReviewDetail {
	if in == nil {
		return nil
	}
	out := &cv1.ReviewDetail{
		Summary:           reviewSummaryFromRv(in.GetSummary()),
		PicUrls:           append([]string(nil), in.GetPicUrls()...),
		VideoUrls:         append([]string(nil), in.GetVideoUrls()...),
		TagsJson:          in.GetTagsJson(),
		GoodsSnapshotJson: in.GetGoodsSnapshotJson(),
	}
	if replies := in.GetReplies(); len(replies) > 0 {
		out.Replies = make([]*cv1.ReplyInfo, 0, len(replies))
		for _, r := range replies {
			out.Replies = append(out.Replies, replyInfoFromRv(r))
		}
	}
	return out
}
