package logic

import (
	"context"
	"time"

	"comment-system/model"
	"comment-system/rpc/comment"
	"comment-system/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCommentLogic) CreateComment(in *comment.CreateCommentRequest) (*comment.CreateCommentResponse, error) {
	rootId := int64(0)

	if in.ParentId > 0 {
		parent, err := l.svcCtx.CommentModel.FindOne(l.ctx, in.ParentId)
		if err != nil {
			return nil, err
		}
		if parent.RootId == 0 {
			rootId = in.ParentId
		} else {
			rootId = parent.RootId
		}
	}

	res, err := l.svcCtx.CommentModel.Insert(l.ctx, &model.Comment{
		VideoId:    in.VideoId,
		UserId:     in.UserId,
		Content:    in.Content,
		ParentId:   in.ParentId,
		RootId:     rootId,
		IpLocation: "未知", // 实际可以从 ctx 获取
		CreatedAt:  time.Now(),
	})
	if err != nil {
		return nil, err
	}

	// 增加回复数逻辑
	if in.ParentId > 0 {
		// 1. 增加直接父评论的回复数
		l.svcCtx.SqlConn.ExecCtx(l.ctx, "UPDATE comment SET reply_count = reply_count + 1 WHERE id = ?", in.ParentId)
		// 2. 如果父评论不是根评论，也要增加根评论的回复数
		if rootId > 0 && rootId != in.ParentId {
			l.svcCtx.SqlConn.ExecCtx(l.ctx, "UPDATE comment SET reply_count = reply_count + 1 WHERE id = ?", rootId)
		}
	}

	newId, _ := res.LastInsertId()

	return &comment.CreateCommentResponse{
		Id:        newId,
		CreatedAt: "刚刚",
	}, nil
}
