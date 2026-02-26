// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"comment-system/api/internal/config"
	"comment-system/model"
	"comment-system/rpc/commentservice"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	CommentRpc commentservice.CommentService
	UserModel  model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:     c,
		CommentRpc: commentservice.NewCommentService(zrpc.MustNewClient(c.CommentRpc)),
		UserModel:  model.NewUserModel(conn),
	}
}
