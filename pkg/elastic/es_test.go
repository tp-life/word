package elastic

import (
	"fmt"
	"testing"
	"time"
)

func TestBoolQuery(t *testing.T) {
	esService := NewSearchService("task_index", "task")
	cond := NewBoolQueryCtx().Should(NewTermQuery("top_kind", 2)).SetShouldMatchMinParam(1).
		Must(NewTermQuery("task_module", 2))
	source, err := cond.Source()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(source)
	res, err := esService.Query(cond).Do()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range res {
		fmt.Println(string(v.Source))
	}
}

func TestBoolQuery1(t *testing.T) {
	esService := NewSearchService("task_index", "task")
	res, err := esService.Query(NewBoolQueryCtx().Should(NewTermQuery("top_kind", 3), NewTermQuery("top_kind", 1)).
		MustNot(NewTermQuery("task_module", 1))).Do()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range res {
		fmt.Println(string(v.Source))
	}
}

func TestBoolQuery2(t *testing.T) {
	esService := NewSearchService("task_index", "task")
	res, err := esService.Query(NewBoolQueryCtx().Should(NewTermQuery("top_kind", 3), NewTermQuery("top_kind", 2))).Size(20).Do()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range res {
		fmt.Println(string(v.Source))
	}
}

func TestBoolQueryFilter(t *testing.T) {
	esService := NewSearchService("task_index", "task")
	res, err := esService.Query(NewBoolQueryCtx().Should(NewTermQuery("top_kind", 3), NewTermQuery("top_kind", 1)).SetShouldMatchMinParam(1).
		Filter(NewTermQuery("task_module", 2))).Do()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range res {
		fmt.Println(string(v.Source))
	}
}

func TestQueryRange(t *testing.T) {
	esService := NewSearchService("task_index", "task")
	//res, err := esService.Query(NewRangeQuery("top_kind", map[string]i{Gte: 2, Lte: 3})).Size(2).Do()
	res, err := esService.Query(NewRangeQuery("top_kind", map[string]interface{}{Gte: 2, Lte: 3})).Size(2).Do()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range res {
		fmt.Println(string(v.Source))
	}
}

func TestBoolQueryFilter1(t *testing.T) {
	esService := NewSearchService("task_index", "task")
	res, err := esService.Query(NewBoolQueryCtx().Filter(NewRangeQuery("top_kind", map[string]interface{}{Gte: 2, Lte: 3}))).Do()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range res {
		fmt.Println(string(v.Source))
	}
}

func TestQueryMultiMatch(t *testing.T) {
	esService := NewSearchService("task_index", "task")
	cond := NewTermsQuery("task_kind", 2, 3)
	source, err := cond.Source()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(source)
	res, err := esService.Query(cond).Do()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range res {
		fmt.Println(string(v.Source))
	}
}

func TestQueryTermsSet(t *testing.T) {
	esService := NewSearchService("task_index", "task")
	cond := NewTermsSetQuery("top_kind", "task_kind", 2, 3)
	source, err := cond.Source()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(source)
	res, err := esService.Query(cond).Do()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range res {
		fmt.Println(string(v.Source))
	}
}

func TestMultiMatchQuery(t *testing.T) {
	esServiceV1 := NewSearchService("task_index", "task")
	cond1 := NewBoolQueryCtx().Should(NewMatchQuery("task_name", "大陆")).SetShouldMatchMinParam(1).
		Must(NewTermsQuery("os_type", 1)).Filter(NewRangeQuery("end_time", map[string]interface{}{Gte: time.Now().Unix()})).
		Filter(NewRangeQuery("start_time", map[string]interface{}{Lte: time.Now().Unix()}))
	data, err := esServiceV1.Query(cond1).Do()
	fmt.Println(err)
	for _, v := range data {
		fmt.Println(string(v.Source))
	}
}
