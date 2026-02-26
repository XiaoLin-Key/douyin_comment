// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package comment

import (
	"context"

	"comment-system/api/internal/svc"
	"comment-system/api/internal/types"
	"comment-system/rpc/comment"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取评论列表
func NewGetCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentsLogic {
	return &GetCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentsLogic) GetComments(req *types.GetCommentsRequest) (resp *types.GetCommentsResponse, err error) {
	rpcResp, err := l.svcCtx.CommentRpc.GetComments(l.ctx, &comment.GetCommentsRequest{
		VideoId:  req.VideoID,
		RootId:   req.RootID,
		UserId:   req.UserID,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	comments := make([]*types.CommentTree, 0)
	for _, c := range rpcResp.Comments {
		comments = append(comments, &types.CommentTree{
			Comment: types.Comment{
				ID:         c.Comment.Id,
				VideoID:    c.Comment.VideoId,
				UserID:     c.Comment.UserId,
				Nickname:      c.Comment.Nickname,
				Avatar:        c.Comment.Avatar,
				IPLocation:    c.Comment.IpLocation,
				IsAuthor:      c.Comment.IsAuthor,
				Content:       c.Comment.Content,
				ParentID:      c.Comment.ParentId,
				ReplyToID:     c.Comment.ReplyToId,
				ReplyToName:   c.Comment.ReplyToName,
				RootID:        c.Comment.RootId,
				LikeCount:     c.Comment.LikeCount,
				ReplyCount:    c.Comment.ReplyCount,
				IsLiked:       c.Comment.IsLiked,
				IsAuthorLiked: c.Comment.IsAuthorLiked,
				CreatedAt:     c.Comment.CreatedAt,
			},
			HasMore:    c.HasMore,
			ReplyTotal: c.ReplyTotal,
			ReplyPage:  c.ReplyPage,
			// 初始Replies根据RPC逻辑可能是空的，等待前端点击“查看更多”
			Replies: make([]*types.CommentTree, 0),
		})
	}

	return &types.GetCommentsResponse{
		Comments: comments,
		Total:    rpcResp.Total,
	}, nil
}
