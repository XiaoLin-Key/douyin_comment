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

type LikeCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 点赞/取消点赞评论
func NewLikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeCommentLogic {
	return &LikeCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeCommentLogic) LikeComment(req *types.LikeCommentRequest) (resp *types.LikeCommentResponse, err error) {
	userIdNumber, ok := l.ctx.Value("userId").(json.Number)
	if !ok {
		userIdNumber = json.Number(fmt.Sprintf("%d", req.UserID))
	}
	userId, _ := userIdNumber.Int64()

	rpcResp, err := l.svcCtx.CommentRpc.LikeComment(l.ctx, &comment.LikeCommentRequest{
		CommentId: req.CommentID,
		UserId:    userId,
		Action:    req.Action,
	})
	if err != nil {
		return nil, err
	}

	return &types.LikeCommentResponse{
		LikeCount: rpcResp.LikeCount,
	}, nil
}
