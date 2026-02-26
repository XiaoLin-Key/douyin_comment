package model

type CommentTree struct {
	Comment *Comment       `json:"comment"`
	User    *User          `json:"user"`     // 评论作者信息
	IsLiked bool           `json:"is_liked"` // 当前用户是否已点赞
	Replies []*CommentTree `json:"replies,omitempty"`
}

type CommentList struct {
	Comments []*CommentTree `json:"comments"`
	Total    int64          `json:"total"`
}
