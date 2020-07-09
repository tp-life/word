package contracts

type Manager struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// ProjectList 专题列表
type ProjectList struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Image   string    `json:"image"`
	Desc    string    `json:"desc"`
	Tags    []string  `json:"tags"`
	Manager []Manager `json:"manager"`  // 管理员
	PostNum int       `json:"post_num"` // 文章数量
}

// CreateProject 添加专题
type CreateProject struct {
	Name  string   `json:"name" form:"name" binding:"required"`
	Image string   `json:"image" form:"image"`
	Desc  string   `json:"desc" form:"desc"`
	Tags  []string `json:"tags" form:"tags"`
}
