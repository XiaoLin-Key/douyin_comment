// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package comment

import (
	"context"
	"encoding/json"
	"fmt"

	"comment-system/api/internal/svc"
	"comment-system/api/internal/types"
	"comment-system/rpc/comment"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发布评论
func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCommentLogic) CreateComment(req *types.CommentRequest) (resp *types.CommentResponse, err error) {
	// 从 JWT Context 中获取用户 ID
	userIdNumber, ok := l.ctx.Value("userId").(json.Number)
	if !ok {
		// 容错处理：如果没拿到（理论上 jwt 开启后不会发生），回退到请求参数
		userIdNumber = json.Number(fmt.Sprintf("%d", req.UserID))
	}
	userId, _ := userIdNumber.Int64()

	rpcResp, err := l.svcCtx.CommentRpc.CreateComment(l.ctx, &comment.CreateCommentRequest{
		VideoId:  req.VideoID,
		UserId:   userId,
		Content:  req.Content,
		ParentId: req.ParentID,
	})
	if err != nil {
		return nil, err
	}

	return &types.CommentResponse{
		ID:        rpcResp.Id,
		CreatedAt: rpcResp.CreatedAt,
	}, nil
}
