package model

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentLikeModel = (*customCommentLikeModel)(nil)

type (
	// CommentLikeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentLikeModel.
	CommentLikeModel interface {
		commentLikeModel
	}

	customCommentLikeModel struct {
		*defaultCommentLikeModel
	}
)

// NewCommentLikeModel returns a model for the database table.
func NewCommentLikeModel(conn sqlx.SqlConn) CommentLikeModel {
	return &customCommentLikeModel{
		defaultCommentLikeModel: newCommentLikeModel(conn),
	}
}
