package entiy

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Base 基础数据模型 方便做事务
type Base struct {
	Ctx       mongo.SessionContext // 事务session
	ID        primitive.ObjectID   `bson:"_id"`
	CreatedAt int64                `bson:"created_at"`
	UpdatedAt int64                `bson:"updated_at"`
}