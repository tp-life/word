package elastic

import (
	"context"
	"word/pkg/common/log"
	"github.com/olivere/elastic/v7"
)

const (
	Gt  = "gt"
	Gte = "gte"
	Lt  = "lt"
	Lte = "lte"
)

// 创建
func Create(index string, _type string, id string, body interface{}) bool {
	connection := GetConnection()
	rs, err := connection.client.Index().Index(index).Type(_type).Id(id).BodyJson(body).Do(context.Background())
	if err != nil {
		log.Error("es create errors: ", err)
	} else {
		if rs.Status == 0 {
			return true
		}
	}
	return false
}

// Delete 删除
func Delete(index string, _type string, id string) bool {
	connection := GetConnection()
	rs, err := connection.client.Delete().Index(index).Type(_type).Id(id).Do(context.Background())
	if err != nil {
		log.Error("es delete errors: ", err)
	}
	if rs.Result == "deleted" {
		return true
	}
	return false
}

// 更新
func Update(index string, _type string, id string, body interface{}) bool {
	connection := GetConnection()
	rs, err := connection.client.Update().Index(index).Type(_type).Id(id).Doc(body).Do(context.Background())
	if err != nil {
		log.Error("es update errors: ", err)
	} else {
		if rs.Status == 201 {
			return true
		}
	}
	return false
}

type QueryCtx interface {
	Source() (interface{}, error)
	ctx() elastic.Query
}

type SearchService struct {
	Es *elastic.SearchService
}

func NewSearchService(index string, _type string) *SearchService {
	connection := GetConnection()
	es := &SearchService{Es: connection.client.Search().Index(index).Type(_type)}
	return es
}

func (e *SearchService) From(from int) *SearchService {
	e.Es = e.Es.From(from)
	return e
}

func (e *SearchService) Size(size int) *SearchService {
	e.Es = e.Es.Size(size)
	return e
}

//条件精确匹配
func (e *SearchService) Query(cond QueryCtx) *SearchService {
	e.Es = e.Es.Query(cond.ctx())
	return e
}

func (e *SearchService) Sort(sort string, ascending bool) *SearchService {
	e.Es = e.Es.Sort(sort, ascending)
	return e
}

func (e *SearchService) Do() (rsp []*elastic.SearchHit, err error) {
	esRsp, err := e.Es.Do(context.Background())
	if err != nil {
		return
	}
	rsp = esRsp.Hits.Hits
	return
}

type NormalQueryCtx struct {
	Query elastic.Query
}

//条件精确匹配
func NewTermQuery(field string, value interface{}) *NormalQueryCtx {
	ctx := NormalQueryCtx{}
	ctx.Query = elastic.NewTermQuery(field, value)
	return &ctx
}

func (ctx *NormalQueryCtx) Source() (interface{}, error) {
	return ctx.Query.Source()
}

func (ctx *NormalQueryCtx) ctx() elastic.Query {
	return ctx.Query
}

//多值查询
func NewTermsQuery(field string, value ...interface{}) *NormalQueryCtx {
	ctx := NormalQueryCtx{}
	ctx.Query = elastic.NewTermsQuery(field, value...)
	return &ctx
}

//多值查询
func NewTermsSetQuery(thresholdField, field string, value ...interface{}) *NormalQueryCtx {
	ctx := NormalQueryCtx{}
	ctx.Query = elastic.NewTermsSetQuery(field, value...).MinimumShouldMatchField(thresholdField)
	return &ctx
}

//分词搜索
func NewMatchQuery(field string, value interface{}) *NormalQueryCtx {
	ctx := NormalQueryCtx{}
	ctx.Query = elastic.NewMatchQuery(field, value)
	return &ctx
}

// 多字段分词搜索
func NewMultiMatchQuery(field []string, value string) *NormalQueryCtx {
	ctx := NormalQueryCtx{}
	ctx.Query = elastic.NewMultiMatchQuery(value, field...)
	return &ctx
}

/*
范围搜索
value参数 eg：map[string]int{
               "gte":10,
               "lt":100
             }
*/
func NewRangeQuery(field string, value map[string]interface{}) *NormalQueryCtx {
	ctx := NormalQueryCtx{}
	condition := elastic.NewRangeQuery(field)
	for k, v := range value {
		switch k {
		case Gt:
			condition.Gt(v)
		case Gte:
			condition.Gte(v)
		case Lt:
			condition.Lt(v)
		case Lte:
			condition.Lte(v)
		}
	}
	ctx.Query = condition
	return &ctx
}

//bool查询上下文
type BoolQueryCtx struct {
	Query *elastic.BoolQuery
}

//bool查询上下文
func NewBoolQueryCtx() *BoolQueryCtx {
	ctx := BoolQueryCtx{}
	ctx.Query = elastic.NewBoolQuery()
	return &ctx
}

func (ctx *BoolQueryCtx) Source() (interface{}, error) {
	return ctx.Query.Source()
}

func (ctx *BoolQueryCtx) ctx() elastic.Query {
	return ctx.Query
}

func (ctx *BoolQueryCtx) Filter(query ...*NormalQueryCtx) *BoolQueryCtx {
	filters := make([]elastic.Query, 0)
	for _, v := range query {
		filters = append(filters, v.Query)
	}
	ctx.Query = ctx.Query.Filter(filters...)
	return ctx
}

func (ctx *BoolQueryCtx) Must(query ...*NormalQueryCtx) *BoolQueryCtx {
	conds := make([]elastic.Query, 0)
	for _, v := range query {
		conds = append(conds, v.Query)
	}
	ctx.Query = ctx.Query.Must(conds...)
	return ctx
}

func (ctx *BoolQueryCtx) MustNot(query ...*NormalQueryCtx) *BoolQueryCtx {
	conds := make([]elastic.Query, 0)
	for _, v := range query {
		conds = append(conds, v.Query)
	}
	ctx.Query = ctx.Query.MustNot(conds...)
	return ctx
}

func (ctx *BoolQueryCtx) Should(query ...*NormalQueryCtx) *BoolQueryCtx {
	conds := make([]elastic.Query, 0)
	for _, v := range query {
		conds = append(conds, v.Query)
	}
	ctx.Query = ctx.Query.Should(conds...)
	return ctx
}

func (ctx *BoolQueryCtx) SetShouldMatchMinParam(num int) *BoolQueryCtx {
	ctx.Query = ctx.Query.MinimumNumberShouldMatch(num)
	return ctx
}
