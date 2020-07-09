package entiy

// Project 专题
type Project struct {
	Base   `bson:"#expand"`
	Name   string   `bson:"name"`
	Image  string   `bson:"image"`  // 专题图片
	Desc   string   `bson:"desc"`   // 专题描述
	Status int8     `bson:"status"` // 0 默认启用 -1 禁用 1 推荐
	Tags   []string `bson:"tags"`   // 标签
}

func (Project) TableName() string { return "project" }
