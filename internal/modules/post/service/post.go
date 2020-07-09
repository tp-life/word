package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"word/internal/entiy"
	"word/internal/modules/post/contracts"
	service2 "word/internal/modules/project/service"
	"word/internal/service"
	"word/pkg/database/mongo"
)

type Post struct {
	ID     string
	userID string
	base   entiy.Post
}

func NewPost(c *gin.Context, ID string) *Post {
	var uid string
	if c != nil {
		user := service.GetUser(c)
		uid = user.ID.Hex()
	}

	return &Post{
		ID:     ID,
		userID: uid,
	}
}

func (p *Post) Info() {
	objID, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		return
	}
	where := bson.M{"_id": objID}
	if p.userID != "" {
		where["user_id"] = p.userID
	}
	_ = mongo.Collection(p.base).Where(where).StFindOne(&p.base)
}

// CreatePost 创建文章
func (p Post) CreatePost(data contracts.CreatePost) error {
	var post = entiy.Post{
		Title:   data.Title,
		Content: data.Content,
		Tags:    data.Tags,
		Status:  data.Status,
		UserID:  p.userID,
	}
	_, err := mongo.Collection(post).InsertOneWithError(post)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePost 更新文章
func (p Post) UpdatePost(data contracts.CreatePost) error {
	p.Info()
	if p.base.ID.IsZero() {
		return errors.New("not found")
	}
	var post = entiy.Post{
		Title:   data.Title,
		Content: data.Content,
		Tags:    data.Tags,
	}
	_, err := mongo.Collection(post).Where(bson.M{"_id": p.base.ID}).UpdateOneWithError(post)
	if err != nil {
		return err
	}
	return nil
}

// UpdateStatus 修改状态
func (p Post) UpdateStatus(d contracts.PostStatus) error {
	if d.Status != entiy.Disable && d.Status != entiy.Push {
		return errors.New("param error")
	}
	p.Info()
	if p.base.ID.IsZero() {
		return errors.New("not found")
	}
	post := p.base
	post.Status = d.Status
	_, err := mongo.Collection(post).Where(bson.M{"_id": p.base.ID}).UpdateOneWithError(post)
	if err != nil {
		return err
	}
	return nil
}

// Lists 文章列表
func (p Post) Lists(q contracts.PostSearch) (rs []contracts.PostList, count int64) {
	cond := p.cond(q)
	var (
		post    entiy.Post
		result  []entiy.Post
		pids    []string
		project service2.Project
	)
	client := mongo.Collection(post).Where(cond)
	count = client.Count()
	_ = client.StFindMany(&result)
	for _, v := range result {
		pids = append(pids, v.ID.Hex())
	}
	pro := project.GetProjectByPostID(pids)
	for _, v := range result {
		var jects = make([]string, 0)
		for _, vv := range pro {
			if v.ID.Hex() == vv.PostID {
				for _, vvv := range vv.Project {
					jects = append(jects, vvv.Name)
				}
				break
			}
		}
		rs = append(rs, contracts.PostList{
			ID:       v.ID.Hex(),
			Projects: jects,
			Title:    v.Title,
			Tags:     v.Tags,
			Status:   v.Status,
			Likes:    v.Likes,
			Fav:      v.Fav,
			Comment:  v.Comment,
		})
	}
	return
}

// JoinProject 加入专题
func (p Post) JoinProject(postID, projectID string) error {
	p.Info()
	if p.base.ID.IsZero() {
		return errors.New("article not found")
	}
	project := service2.ProjectPostM{}
	err := project.ApplyPost(postID, projectID)
	return err
}

// Detail 详情
func (p Post) Detail() (rs contracts.PostDetail, err error) {
	p.Info()
	if p.base.ID.IsZero() {
		err = errors.New("not found")
		return
	}
	var (
		project service2.Project
		pr      []string
	)
	pro := project.GetProjectByPostID([]string{p.ID})
	for _, vv := range pro {
		if p.ID == vv.PostID {
			for _, vvv := range vv.Project {
				pr = append(pr, vvv.Name)
			}
			break
		}
	}
	rs = contracts.PostDetail{
		PostList: contracts.PostList{
			ID:       p.ID,
			Projects: pr,
			Title:    p.base.Title,
			Tags:     p.base.Tags,
			Status:   p.base.Status,
			Likes:    p.base.Likes,
			Fav:      p.base.Fav,
			Comment:  p.base.Comment,
		},
		Content: p.base.Content,
	}
	return
}

// Delete 删除
func (p Post) Delete() (err error) {
	p.Info()
	if p.base.ID.IsZero() {
		err = errors.New("not found")
		return
	}
	_, err = mongo.Collection(p.base).Where(bson.M{"_id": p.base.ID}).DeleteWithError()
	return
}

// cond 查询条件
func (p Post) cond(q contracts.PostSearch) bson.M {
	where := bson.M{}
	if p.userID != "" {
		where["user_id"] = p.userID
	}
	if q.Status != "" {
		s, err := strconv.Atoi(q.Status)
		if err == nil {
			where["status"] = s
		}
	}
	if q.Title != "" {
		where["title"] = bson.M{"$regex": q.Title}
	}
	if len(q.Tags) > 0 {
		where["tags"] = bson.M{"$in": q.Tags}
	}
	if len(q.ID) > 0 {
		var pids []primitive.ObjectID
		for _, v := range q.ID {
			if obj, err := primitive.ObjectIDFromHex(v); err == nil {
				pids = append(pids, obj)
			}
		}
		if len(pids) > 0 {
			where["_id"] = bson.M{"$in": pids}
		}
	}
	if q.Project != "" {
		var (
			projects []entiy.ProjectPost
			pids     []primitive.ObjectID
		)
		mongo.Collection(entiy.ProjectPost{}).Where(bson.M{"project_id": q.Project}).FindMany(&projects)
		if len(projects) > 0 {
			for _, v := range projects {
				if obj, err := primitive.ObjectIDFromHex(v.PostID); err == nil {
					pids = append(pids, obj)
				}
			}
			if len(pids) > 0 {
				where["_id"] = bson.M{"$in": pids}
			}
		}
	}
	return where
}
