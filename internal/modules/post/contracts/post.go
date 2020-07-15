package contracts

// CreatePost 发布文章
type CreatePost struct {
	Tags    []string `json:"tags" form:"tags"`     // 标签
	Status  int8     `json:"status" form:"status"` // 发布状态 1 发布
	Content string   `json:"content" form:"content" binding:"required"`
	Title   string   `json:"title" form:"title" binding:"required"`
}

// PostStatus 更改文章类型
type PostStatus struct {
	Status int8 `json:"status" form:"status" binding:"required"` // 发布状态 1 发布 0 草稿箱 -1 禁用
}

// PostDetail 文章类型
type PostDetail struct {
	PostList
}

// PostSearch 搜索条件
type PostSearch struct {
	Title  string   `json:"title" form:"title"`   // 标题
	Tags   []string `json:"tags" form:"tags"`     // 标签
	Status string   `json:"status" form:"status"` // 状态
	ID     []string `json:"id" form:"id"`
	Project string `json:"project" form:"project"`
}

// PostList 文章列表
type PostList struct {
	ID       string   `json:"id"`
	Projects []string `json:"projects"` // 被哪些专题收录
	Title    string   `json:"title" form:"title" binding:"required"`
	Tags     []string `json:"tags" form:"tags"`     // 标签
	Status   int8     `json:"status" form:"status"` // 发布状态 1 发布 0 草稿箱 -1 禁用
	Likes    int      `json:"likes"`                // 喜欢数
	Fav      int      `json:"fav"`                  // 收藏数
	Comment  int      `json:"comment"`              // 评论数
	Content string `json:"content"`
}

// JoinProject 加入专题
type JoinProject struct {
	Project string `json:"project" form:"project" binding:"required"`
}