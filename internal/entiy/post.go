package entiy

const (
	Push    = 1
	Disable = -1
)

// Post 文章
type Post struct {
	Base    `bson:"#expand"`
	Title   string   `bson:"title"`
	Content string   `bson:"content"`
	Tags    []string `bson:"tags"`    // 标签
	Status  int8     `bson:"status"`  // 0 草稿箱  1 发布  -1 下架
	UserID  string   `bson:"user_id"` // 所属用户
	Likes   int      `bson:"likes"`   // 喜欢数
	Fav     int      `bson:"fav"`     // 收藏数
	Comment int      `bson:"comment"` // 评论数
}

func (Post) TableName() string { return "post" }
