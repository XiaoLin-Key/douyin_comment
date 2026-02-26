package logic

import (
	"context"
	"fmt"
	"time"

	"comment-system/model"
	"comment-system/rpc/comment"
	"comment-system/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentsLogic {
	return &GetCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 格式化时间显示规则
func formatRelativeTime(t time.Time) string {
	duration := time.Since(t)
	if duration < 0 { // 处理时钟回拨或微小误差
		return "刚刚"
	}
	if duration.Minutes() < 5 {
		return "刚刚"
	} else if duration.Hours() < 1 {
		return fmt.Sprintf("%d分钟之前", int(duration.Minutes()))
	} else if duration.Hours() < 24 {
		return fmt.Sprintf("%d小时之前", int(duration.Hours()))
	} else {
		return t.Format("01-02") // 超过一天显示月日
	}
}

func (l *GetCommentsLogic) GetComments(in *comment.GetCommentsRequest) (*comment.GetCommentsResponse, error) {
	if in.RootId == 0 {
		return l.getTopLevelComments(in)
	}
	return l.getSubComments(in)
}

func (l *GetCommentsLogic) getTopLevelComments(in *comment.GetCommentsRequest) (*comment.GetCommentsResponse, error) {
	// 1. 查询顶级评论
	// 实际生产中应使用带权重的排序，这里先实现基础的分页查询
	query := "SELECT id, video_id, user_id, content, parent_id, root_id, like_count, reply_count, ip_location, is_deleted, created_at " +
		"FROM comment WHERE video_id = ? AND root_id = 0 AND is_deleted = 0 ORDER BY created_at DESC LIMIT ? OFFSET ?"

	var comments []*model.Comment
	err := l.svcCtx.SqlConn.QueryRowsCtx(l.ctx, &comments, query, in.VideoId, in.PageSize, (in.Page-1)*in.PageSize)
	if err != nil {
		return nil, err
	}

	var respComments []*comment.CommentTree
	for _, c := range comments {
		// 查询用户信息
		u, _ := l.svcCtx.UserModel.FindOne(l.ctx, c.UserId)
		nickname, avatar := "匿名用户", ""
		if u != nil {
			nickname = u.Nickname
			avatar = u.Avatar
		}

		// 判断当前用户是否点赞 (如果 in.UserId > 0)
		isLiked := false
		if in.UserId > 0 {
			checkQuery := "SELECT count(*) FROM comment_like WHERE comment_id = ? AND user_id = ?"
			var count int64
			_ = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &count, checkQuery, c.Id, in.UserId)
			isLiked = count > 0
		}

		// 3. 动态计算该顶级评论下所有未删除的回复总数 (ReplyTotal)
		var actualReplyCount int64
		countQuery := "SELECT count(*) FROM comment WHERE root_id = ? AND is_deleted = 0"
		_ = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &actualReplyCount, countQuery, c.Id)

		respComments = append(respComments, &comment.CommentTree{
			Comment: &comment.Comment{
				Id:         c.Id,
				VideoId:    c.VideoId,
				UserId:     c.UserId,
				Nickname:   nickname,
				Avatar:     avatar,
				Content:    c.Content,
				ParentId:   c.ParentId,
				RootId:     c.RootId,
				LikeCount:  c.LikeCount,
				ReplyCount: actualReplyCount, // 使用实时计算的回复数
				IpLocation: c.IpLocation,
				IsLiked:    isLiked,
				CreatedAt:  formatRelativeTime(c.CreatedAt),
			},
			HasMore:    actualReplyCount > 0,
			ReplyTotal: actualReplyCount,
		})
	}

	// 2. 查询总评论数 (包括顶级和子评论，但不包括被软删除的)
	var total int64
	totalQuery := "SELECT count(*) FROM comment WHERE video_id = ? AND is_deleted = 0"
	_ = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &total, totalQuery, in.VideoId)

	return &comment.GetCommentsResponse{
		Comments: respComments,
		Total:    total,
	}, nil
}

func (l *GetCommentsLogic) getSubComments(in *comment.GetCommentsRequest) (*comment.GetCommentsResponse, error) {
	query := "SELECT id, video_id, user_id, content, parent_id, root_id, like_count, reply_count, ip_location, is_deleted, created_at " +
		"FROM comment WHERE root_id = ? AND is_deleted = 0 ORDER BY created_at ASC LIMIT ? OFFSET ?"

	var comments []*model.Comment
	err := l.svcCtx.SqlConn.QueryRowsCtx(l.ctx, &comments, query, in.RootId, in.PageSize, (in.Page-1)*in.PageSize)
	if err != nil {
		return nil, err
	}

	var respComments []*comment.CommentTree
	for _, c := range comments {
		u, _ := l.svcCtx.UserModel.FindOne(l.ctx, c.UserId)
		nickname, avatar := "匿名用户", ""
		if u != nil {
			nickname = u.Nickname
			avatar = u.Avatar
		}

		// 查询被回复者的昵称 (如果 parent_id > 0)
		replyToName := ""
		if c.ParentId > 0 && c.ParentId != c.RootId {
			parent, _ := l.svcCtx.CommentModel.FindOne(l.ctx, c.ParentId)
			if parent != nil {
				pu, _ := l.svcCtx.UserModel.FindOne(l.ctx, parent.UserId)
				if pu != nil {
					replyToName = pu.Nickname
				}
			}
		}

		respComments = append(respComments, &comment.CommentTree{
			Comment: &comment.Comment{
				Id:          c.Id,
				VideoId:     c.VideoId,
				UserId:      c.UserId,
				Nickname:    nickname,
				Avatar:      avatar,
				Content:     c.Content,
				ParentId:    c.ParentId,
				RootId:      c.RootId,
				ReplyToName: replyToName,
				CreatedAt:   formatRelativeTime(c.CreatedAt),
			},
		})
	}

	return &comment.GetCommentsResponse{
		Comments: respComments,
	}, nil
}
