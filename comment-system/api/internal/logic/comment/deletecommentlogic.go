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

type DeleteCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除评论
func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCommentLogic) DeleteComment(req *types.DeleteCommentRequest) (resp *types.DeleteCommentResponse, err error) {
	userIdNumber, ok := l.ctx.Value("userId").(json.Number)
	if !ok {
		userIdNumber = json.Number(fmt.Sprintf("%d", req.UserID))
	}
	userId, _ := userIdNumber.Int64()

	rpcResp, err := l.svcCtx.CommentRpc.DeleteComment(l.ctx, &comment.DeleteCommentRequest{
		CommentId: req.CommentID,
		UserId:    userId,
	})
	if err != nil {
		return nil, err
	}

	return &types.DeleteCommentResponse{
		Success: rpcResp.Success,
	}, nil
}
