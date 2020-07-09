package service

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"word/internal/entiy"
	"word/pkg/database/mongo"
)

type ProjectPostM struct {
}

// AuditPost 申请将文章添加入专题
func (ProjectPostM) ApplyPost(postID string, projectID string) error {
	project := Project{}
	err := project.checkInfo(projectID)
	if err != nil {
		return err
	}
	projectPost := entiy.ProjectPost{
		PostID:    postID,
		ProjectID: projectID,
	}
	_, err = mongo.Collection(projectPost).InsertOneWithError(projectPost)
	return err
}

// AuditPost 审核专题
func (p ProjectPostM) AuditPost(projectPostID string, state int8) error {
	o, e := primitive.ObjectIDFromHex(projectPostID)
	if e != nil {
		return e
	}
	_, err := mongo.Collection(entiy.ProjectPost{}).Where(bson.M{"_id": o}).UpdateOneWithError(bson.M{"status": state})
	return err
}
