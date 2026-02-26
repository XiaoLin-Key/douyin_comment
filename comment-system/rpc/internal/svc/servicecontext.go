package svc

import (
	"comment-system/model"
	"comment-system/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	SqlConn          sqlx.SqlConn
	CommentModel     model.CommentModel
	UserModel        model.UserModel
	CommentLikeModel model.CommentLikeModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:           c,
		SqlConn:          conn,
		CommentModel:     model.NewCommentModel(conn),
		UserModel:        model.NewUserModel(conn),
		CommentLikeModel: model.NewCommentLikeModel(conn),
	}
}
