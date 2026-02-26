package logic

import (
	"context"

	"comment-system/model"
	"comment-system/rpc/comment"
	"comment-system/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeCommentLogic {
	return &LikeCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeCommentLogic) LikeComment(in *comment.LikeCommentRequest) (*comment.LikeCommentResponse, error) {
	if in.Action == 1 {
		_, err := l.svcCtx.CommentLikeModel.Insert(l.ctx, &model.CommentLike{
			CommentId: in.CommentId,
			UserId:    in.UserId,
		})
		if err != nil {
			return nil, err
		}
		query := "UPDATE comment SET like_count = like_count + 1 WHERE id = ?"
		_, err = l.svcCtx.SqlConn.ExecCtx(l.ctx, query, in.CommentId)
		if err != nil {
			return nil, err
		}
	} else {
		queryDel := "DELETE FROM comment_like WHERE comment_id = ? AND user_id = ?"
		_, err := l.svcCtx.SqlConn.ExecCtx(l.ctx, queryDel, in.CommentId, in.UserId)
		if err != nil {
			return nil, err
		}
		queryUpd := "UPDATE comment SET like_count = GREATEST(like_count - 1, 0) WHERE id = ?"
		_, err = l.svcCtx.SqlConn.ExecCtx(l.ctx, queryUpd, in.CommentId)
		if err != nil {
			return nil, err
		}
	}

	comm, err := l.svcCtx.CommentModel.FindOne(l.ctx, in.CommentId)
	if err != nil {
		return nil, err
	}

	return &comment.LikeCommentResponse{
		LikeCount: comm.LikeCount,
	}, nil
}
