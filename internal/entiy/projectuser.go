package entiy

// ProjectUser 专题管理员
type ProjectUser struct {
	Base      `bson:"#expand"`
	ProjectID string `bson:"project_id"`
	UserID    string `bson:"user_id"`
	Status    int8   `bson:"status"` // -1 禁用 1 正常  0 待确认
}

func (ProjectUser) TableName() string { return "project_user" }
