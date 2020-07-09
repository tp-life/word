package entiy

const (
	DISABLE = iota - 1
	DEFAULT
	ENABLE
)

// ProjectPost 专题于文章对应
type ProjectPost struct {
	Base      `bson:"#expand"`
	PostID    string `bson:"post_id"`
	ProjectID string `bson:"project_id"`
	Status    int8   `bson:"status"` // -1 驳回 0 待审核 1 已收录
}

func (ProjectPost) TableName() string { return "project_post" }
