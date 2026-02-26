package logic

import (
	"context"
	"fmt"

	"comment-system/rpc/comment"
	"comment-system/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCommentLogic) DeleteComment(in *comment.DeleteCommentRequest) (*comment.DeleteCommentResponse, error) {
	// 1. 先查询评论的基本信息（主要是 parent_id 和 root_id）
	var c struct {
		Id       int64 `db:"id"`
		ParentId int64 `db:"parent_id"`
		RootId   int64 `db:"root_id"`
		UserId   int64 `db:"user_id"`
	}
	queryFind := "SELECT id, parent_id, root_id, user_id FROM comment WHERE id = ?"
	err := l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &c, queryFind, in.CommentId)
	if err != nil {
		l.Errorf("查询评论失败: %v", err)
		return nil, err
	}

	// 权限校验：只能删除自己的
	if c.UserId != in.UserId {
		return nil, fmt.Errorf("没有权限删除此评论")
	}

	if c.ParentId == 0 {
		// 顶级评论：物理删除自己及其所有子回复
		queryDel := "DELETE FROM comment WHERE id = ? OR root_id = ?"
		_, err = l.svcCtx.SqlConn.ExecCtx(l.ctx, queryDel, in.CommentId, in.CommentId)
		if err != nil {
			l.Errorf("物理删除顶级评论失败: %v", err)
			return nil, err
		}
	} else {
		// 子评论：软删除
		queryUpdate := "UPDATE comment SET content = '该评论已删除', is_deleted = 1 WHERE id = ?"
		_, err = l.svcCtx.SqlConn.ExecCtx(l.ctx, queryUpdate, in.CommentId)
		if err != nil {
			l.Errorf("软删除子评论失败: %v", err)
			return nil, err
		}

		// 2. 减少父评论和根评论的回复数
		if c.ParentId > 0 {
			l.svcCtx.SqlConn.ExecCtx(l.ctx, "UPDATE comment SET reply_count = reply_count - 1 WHERE id = ? AND reply_count > 0", c.ParentId)
			if c.RootId > 0 && c.RootId != c.ParentId {
				l.svcCtx.SqlConn.ExecCtx(l.ctx, "UPDATE comment SET reply_count = reply_count - 1 WHERE id = ? AND reply_count > 0", c.RootId)
			}
		}
	}

	return &comment.DeleteCommentResponse{
		Success: true,
	}, nil
}
