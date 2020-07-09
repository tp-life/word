package service

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"word/internal/entiy"
	"word/internal/modules/project/contracts"
	"word/pkg/database/mongo"
)

type Project struct {
	pid primitive.ObjectID
	pro entiy.Project
}

type ProjectPost struct {
	PostID  string
	Project []entiy.Project
}

// GetProjectByPostID 通过postID 获取专题
func (Project) GetProjectByPostID(postIDs []string) (rs []ProjectPost) {
	if len(postIDs) == 0 {
		return
	}
	var (
		pp  entiy.ProjectPost
		prs []entiy.Project
		pid []primitive.ObjectID
		pro []struct {
			PostID   string   `bson:"post_id"`
			Projects []string `bson:"projects"`
		}
		pmap = make(map[string]entiy.Project)
	)
	pipe := []bson.M{
		{"$match": bson.M{"post_id": postIDs}},
		{"$group": bson.M{"_id": "$post_id", "projects": bson.M{"$push": "$project_id"}}},
		{"$project": bson.M{"post_id": "$_id", "projects": 1, "_id": -1}},
	}
	_ = mongo.Collection(pp).AggregateWithError(pipe, &pro)
	for _, v := range pro {
		for _, vv := range v.Projects {
			id, _ := primitive.ObjectIDFromHex(vv)
			pid = append(pid, id)
		}
	}
	if len(pid) == 0 {
		return
	}
	_ = mongo.Collection(entiy.Project{}).Where(bson.M{"_id": bson.M{"$in": pid}}).StFindMany(&prs)
	for _, v := range prs {
		pmap[v.ID.Hex()] = v
	}
	for _, v := range pro {
		var temp []entiy.Project
		for _, vv := range v.Projects {
			if tem, ok := pmap[vv]; ok {
				temp = append(temp, tem)
			}
		}
		rs = append(rs, ProjectPost{
			PostID:  v.PostID,
			Project: temp,
		})
	}
	return
}

// Projects 获取专题列表
func (p Project) Projects() (result []contracts.ProjectList) {
	var (
		rs        []entiy.Project
		projectID []string
	)
	_ = mongo.Collection(entiy.Project{}).Where(bson.M{"status": bson.M{"$gte": 0}}).Sort(bson.M{"status": -1}).StFindMany(&rs)
	for _, v := range rs {
		projectID = append(projectID, v.ID.Hex())
	}
	projectUser := p.ProjectUser(projectID)
	projectPostNo := p.ProjectPostNo(projectID)
	for _, v := range rs {
		pl := contracts.ProjectList{
			ID:      v.ID.Hex(),
			Name:    v.Name,
			Image:   v.Image,
			Desc:    v.Desc,
			Tags:    v.Tags,
			Manager: nil,
			PostNum: 0,
		}
		if no, ok := projectPostNo[v.ID.Hex()]; ok {
			pl.PostNum = no
		}
		if us, ok := projectUser[v.ID.Hex()]; ok {
			var temp = make([]contracts.Manager, 0)
			for _, u := range us {
				temp = append(temp, contracts.Manager{
					ID:     u.ID.Hex(),
					Name:   u.Name,
					Avatar: u.Avatar,
				})
			}
			pl.Manager = temp
		}
		result = append(result, pl)
	}
	return
}

// Create 添加
func (Project) Create(q contracts.CreateProject) (err error) {
	insertData := entiy.Project{
		Name:  q.Name,
		Image: q.Image,
		Desc:  q.Desc,
		Tags:  q.Tags,
	}
	_, err = mongo.Collection(insertData).InsertOneWithError(insertData)
	return
}

func (p *Project) Info() error {
	if p.pid.IsZero() {
		return errors.New("not found project")
	}
	err := mongo.Collection(entiy.Project{}).Where(bson.M{"_id": p.pid}).StFindOne(&p.pro)
	return err
}

// SetID 设置PID
func (p *Project) SetID(pid string) {
	o, _ := primitive.ObjectIDFromHex(pid)
	p.pid = o
}

// Update 修改
func (p Project) Update(q contracts.CreateProject, pid string) error {
	err := p.checkInfo(pid)
	if err != nil {
		return err
	}
	insertData := entiy.Project{
		Name:  q.Name,
		Image: q.Image,
		Desc:  q.Desc,
		Tags:  q.Tags,
	}
	_, err = mongo.Collection(p.pro).Where(bson.M{"_id": p.pid}).UpdateOneWithError(insertData)
	return nil
}

// Delete 删除
func (p Project) Delete(pid string) error {
	err := p.checkInfo(pid)
	if err != nil {
		return err
	}
	_, err = mongo.Collection(p.pro).Where(bson.M{"_id": p.pid}).DeleteWithError()
	return err
}

// ChangeState 修改状态
func (p Project) ChangeState(pid string, state int8) error {
	err := p.checkInfo(pid)
	if err != nil {
		return err
	}
	p.pro.Status = state
	_, err = mongo.Collection(p.pro).Where(bson.M{"_id": p.pro.ID}).UpdateOneWithError(p.pro)
	return err
}

func (p Project) checkInfo(pid string) error {
	p.SetID(pid)
	err := p.Info()
	return err
}

// ProjectUser 获取专题管理人员
func (Project) ProjectUser(pids []string) (rs map[string][]entiy.User) {
	if len(pids) == 0 {
		return
	}
	var (
		pu    []entiy.ProjectUser
		users []entiy.User
		uids  []struct {
			user primitive.ObjectID
			pid  primitive.ObjectID
		}
	)
	rs = make(map[string][]entiy.User)
	mongo.Collection(entiy.ProjectUser{}).Where(bson.M{"project_id": bson.M{"$in": pids}}).FindMany(&pu)
	for _, v := range pu {
		o, e := primitive.ObjectIDFromHex(v.UserID)
		if e != nil {
			continue
		}
		uids = append(uids, struct {
			user primitive.ObjectID
			pid  primitive.ObjectID
		}{user: o, pid: v.ID})
	}
	_ = mongo.Collection(entiy.User{}).StFindMany(&users)
	for _, v := range uids {
		for _, user := range users {
			if v.user == user.ID {
				temp := []entiy.User{user}
				if us, ok := rs[v.pid.Hex()]; ok {
					temp = append(us, user)
				}
				rs[v.pid.Hex()] = temp
				break
			}
		}
	}
	return
}

// ProjectPostNo 获取专题下文章数量
func (Project) ProjectPostNo(pids []string) (result map[string]int) {
	if len(pids) == 0 {
		return
	}
	result = make(map[string]int)
	var (
		where = bson.M{"project_id": bson.M{"$in": pids}}
		rs    []struct {
			Project string `bson:"project_id"`
			Count   int    `bson:"count"`
		}
	)
	pipe := []bson.M{
		{"$match": where},
		{"$group": bson.M{"_id": "$project_id", "count": bson.M{"$sum": 1}}},
		{"$project": bson.M{"project": "$_id", "count": 1, "_id": 0}},
	}
	_ = mongo.Collection(entiy.ProjectPost{}).AggregateWithError(pipe, &rs)
	for _, v := range rs {
		result[v.Project] = v.Count
	}
	return
}
