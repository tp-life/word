package elastic

import (
	"context"
	"word/pkg/app"
	"word/pkg/unique"
	jsoniter "github.com/json-iterator/go"
	"github.com/olivere/elastic/v7"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"
)

var (
	mapping = `
	{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"properties":{
				"timestamp": {
                    "type": "date"
                },
				"location":{
					"type":"geo_point"
				}
			}
		}
	}
`

	ipData   *geoip2.Reader
	loadOnce sync.Once
)

// Body 更新时可使用的文档结构
type Body map[string]interface{}

type iInit interface {
	init()
}

// Location 如果包含此结构体
type Location struct {
	*BaseItem
	IP       string            `json:"ip"`
	Location *elastic.GeoPoint `json:"location"`
	Country  string            `json:"country"`
	City     string            `json:"city"`
}

// BaseItem 数据
type BaseItem struct {
	Timestamp string `json:"timestamp"`
}

func (baseItem *BaseItem) init() {
	if baseItem.Timestamp == "" {
		baseItem.Timestamp = time.Now().Format(time.RFC3339)
	}
}

func (location *Location) init() {
	location.BaseItem.init()
	loadOnce.Do(func() {
		var err error
		ipData, err = geoip2.Open(app.Root() + "/assets/geo/ipdb.mmdb")
		if err != nil {
			app.Logger().Printf("pkg.elastic.elastic, err: %v\n", err)
		}
	})
	record, _ := ipData.City(net.ParseIP(location.IP))
	location.Country = record.Country.Names["zh-CN"]
	location.City = record.City.Names["zh-CN"]
	location.Location = &elastic.GeoPoint{
		Lat: record.Location.Latitude,
		Lon: record.Location.Longitude,
	}
}

// NewBaseItem 填充
func NewBaseItem() *BaseItem {
	return new(BaseItem)
}

func InitLocation(location *Location) *Location {
	location.init()
	return location
}

// NewLocation 填充Location
func NewLocation(ip string) *Location {
	return &Location{
		IP:       ip,
		BaseItem: NewBaseItem(),
	}
}

// Connection 连接对象
type Connection struct {
	client *elastic.Client
	list   chan *Data
}

// Data 数据
type Data struct {
	Index string
	ID    string
	Body  iInit
}

var (
	connection *Connection
	single     sync.Once
)

// GetConnection 获取连接对象
func GetConnection() *Connection {
	single.Do(func() {
		connection = new(Connection)
		connection.list = make(chan *Data, 4096)
		connection.start()
	})
	return connection
}

func id() string {
	id := unique.ID()
	return strconv.Itoa(int(id))
}

func (connection *Connection) handle() {
	for {
		select {
		case item := <-connection.list:
			connection.createIndex(item.Index)
			_, err := connection.client.Index().Index(item.Index).Id(item.ID).BodyJson(item.Body).Do(context.Background())
			if err != nil {
				app.Logger().WithField("log_type", "pkg.elastic.elastic").Warn("insert body error: ", err)
			}
		}
	}
}

// start 连接 es
func (connection *Connection) start() {
	var err error
	connection.client, err = elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false))
	if err != nil {
		log.Println("connect elastic search error: ", err)
	}
	go connection.handle()
}

func (connection *Connection) createIndex(index string) {
	exists, err := connection.client.IndexExists(index).Do(context.Background())
	if err != nil {
		log.Println("search index err: ", err)
	}

	if exists {
		return
	}
	_, err = connection.client.CreateIndex(index).Body(mapping).Do(context.Background())
	if err != nil {
		log.Println("create index error: ", err)
		return
	}
}

// Save 保存文档, 返回id
func (connection *Connection) Save(index, uniqueID string, body iInit) string {
	body.init()
	if uniqueID == "" {
		uniqueID = id()
	}
	connection.list <- &Data{
		ID:    uniqueID,
		Index: index,
		Body:  body,
	}
	return uniqueID
}

// Update 更新文档
//  es.GetConnection().Update("index", "id", elastic.NewScript("ctx._source.your_field += num").Param("num", 10), es.Body{"field_1": "hello", "field_2": "world"})
//  如果 script 为 nil, 将把 body 作为 doc upsert 操作执行请求, 否则已脚本为主, body 作为辅助 upsert 数据, 只有在 elasticsearch 中不存在该字段时才会更新
func (connection *Connection) Update(index string, id string, script *elastic.Script, body Body) {
	var (
		response *elastic.UpdateResponse
		err      error
	)
	if script == nil {
		response, err = connection.client.Update().Index(index).Id(id).Doc(body).DocAsUpsert(true).Do(context.Background())
	} else {
		response, err = connection.client.Update().Index(index).Id(id).Script(script).Upsert(body).Do(context.Background())
	}
	if err != nil {
		log.Println("update document error: ", err)
		return
	}
	log.Println(response.Id)
}

// Find 查找指定文档
func (connection *Connection) Find(index, id string, structure interface{}) {
	response, err := connection.client.Get().Index(index).Id(id).Do(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	_ = jsoniter.Unmarshal(response.Source, structure)
}

// Client 连接客户端
func (connection *Connection) Client() *elastic.Client {
	return connection.client
}

// Search 搜索文档
func (connection *Connection) Search(index string) *Search {
	return &Search{Index: index, connection: connection}
}

// Search 搜索
type Search struct {
	Index      string
	query      *elastic.TermsQuery
	sort       map[string]bool // 排序, key 为字段, 值为是否是升序, true 升序排序, false 降序
	connection *Connection
	From       int // 翻页参数 , 从多少条开始
	Size       int // 翻页参数，取多少条
}

// Query 指定搜索参数
func (search *Search) Query(query *elastic.TermsQuery) *Search {
	search.query = query
	return search
}

// Sort 排序
func (search *Search) Sort(key string, ascending bool) *Search {
	if search.sort == nil {
		search.sort = make(map[string]bool)
	}
	search.sort[key] = ascending
	return search
}

// Page 分页参数
func (search *Search) Page(start, size int) *Search {
	search.From = start
	search.Size = size
	return search
}

// Do 搜索数据, Do 和 Each 是两种操作, Each 包含Do, 但是Do 不包含 Each
func (search *Search) Do(structure interface{}) (total int64, items []interface{}) {
	var ctx = context.Background()
	s := connection.client.Search().Index(search.Index).Query(search.query)
	if search.sort != nil {
		for field, ascending := range search.sort {
			s.Sort(field, ascending)
		}
	}
	if search.Size == 0 {
		search.Size = 12
	}

	response, err := s.From(search.From).Size(search.Size).Do(ctx)
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	return response.TotalHits(), response.Each(reflect.TypeOf(structure))
}

// Each 遍历数据
func (search *Search) Each(structure interface{}, handle func(item interface{}) (continued bool)) {
	total, items := search.Do(structure)
	if total == 0 {
		return
	}

	for _, item := range items {
		if handle(item) == false {
			return
		}
	}

	nextPage := search.From + search.Size
	if int64(nextPage) > total {
		return
	}

	search.From = nextPage
	search.Each(structure, handle)
}
