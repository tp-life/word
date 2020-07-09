// Package mongo 三方 mgo 太长时间不维护了
//  官方 mongo 驱动很不友好
//  所以这里稍微对常用方法做了处理,可以直接调用这里的方法进行一些常规操作
//  复杂的操作,调用这里的 Collection 之后可获取里边的 Database 属性 和 Table 属性操作
//  这里的添加和修改操作将会自动补全 created_at updated_at 和 _id
package mongo

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"
	"strings"
	"time"
	"word/pkg/app"
	"word/pkg/database"
	"word/pkg/unique"
)

var (
	client  *mongo.Client
	conf    config
	indexes = make(chan func(), 256)
)

type (
	// CollectionInfo 集合包含的连接信息和查询等操作信息
	CollectionInfo struct {
		Database   *mongo.Database
		Table      *mongo.Collection
		Context    context.Context
		filter     bson.M
		limit      int64
		skip       int64
		sort       bson.M
		fields     bson.M
		SafeDelete bool
	}

	config struct {
		URL             string `mapstructure:"url" toml:"url" env:"MONGO_URL"`
		Database        string `mapstructure:"database" toml:"database" env:"MONGO_DB"`
		MaxConnIdleTime int    `mapstructure:"max_conn_idle_time" toml:"max_conn_idle_time" env:"MONGO_MAX_CONN_TIME"`
		MaxPoolSize     int    `mapstructure:"max_pool_size" toml:"max_pool_size" env:"MONGO_MAX_POOL_SIZE"`
		Username        string `mapstructure:"username" toml:"username" env:"MONGO_USER"`
		Password        string `mapstructure:"password" toml:"password" env:"MONGO_PASS"`
	}

	// Transaction 事务
	Transaction struct {
		Session mongo.SessionContext
	}
)

// Start 启动 mongo
func Start() {
	var err error
	err = app.InitConfig("mongo", &conf)
	if err != nil {
		app.Logger().Fatalln("unable to decode mongo config", err)
	}
	mongoOptions := options.Client()
	mongoOptions.SetMaxConnIdleTime(time.Duration(conf.MaxConnIdleTime) * time.Second)
	mongoOptions.SetMaxPoolSize(uint64(conf.MaxPoolSize))
	if conf.Username != "" && conf.Password != "" {
		mongoOptions.SetAuth(options.Credential{Username: conf.Username, Password: conf.Password})
	}

	client, err = mongo.NewClient(mongoOptions.ApplyURI(conf.URL))
	if err != nil {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error(err)
	}

	go createIndex()
}

func createIndex() {
	for {
		select {
		case index := <-indexes:
			index()
		}
	}
}

// CreateIndex 创建索引
func CreateIndex(index func()) {
	indexes <- index
}

// Collection 得到一个mongo操作对象
func Collection(table database.Table) *CollectionInfo {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	db := client.Database(conf.Database)
	return &CollectionInfo{
		Database:   db,
		Table:      db.Collection(table.TableName()),
		Context:    ctx,
		filter:     make(bson.M),
		SafeDelete: true,
	}
}

// CollectionWithSession 得到一个mongo操作对象, 接收传入的context
func CollectionWithSession(table database.Table, ctx context.Context) *CollectionInfo {
	dataBase := client.Database(conf.Database)
	return &CollectionInfo{
		Database: dataBase,
		Table:    dataBase.Collection(table.TableName()),
		filter:   make(bson.M),
		Context:  ctx,
	}
}

// Database 获取数据库连接
func Database(name ...string) *CollectionInfo {
	var db *mongo.Database
	if len(name) == 1 {
		db = client.Database(name[0])
	}
	db = client.Database(conf.Database)
	collection := &CollectionInfo{
		Database: db,
		filter:   make(bson.M),
	}
	return collection
}

// SetSessionCxt 设置事务上下文
func (collection *CollectionInfo) SetSessionCxt(ctx mongo.SessionContext) *CollectionInfo {
	if ctx != nil {
		collection.Context = ctx
	}
	return collection
}

// SetTable 设置集合名称
func (collection *CollectionInfo) SetTable(name string) *CollectionInfo {
	collection.Table = collection.Database.Collection(name)
	return collection
}

func (collection *CollectionInfo) SetSafe(safe bool) *CollectionInfo {
	collection.SafeDelete = safe
	return collection
}

// Where 条件查询, bson.M{"field": "value"}
func (collection *CollectionInfo) Where(m bson.M) *CollectionInfo {
	if collection.SafeDelete {
		m["deleted_at"] = 0
	}
	collection.filter = m
	return collection
}

// Limit 限制条数
func (collection *CollectionInfo) Limit(n int64) *CollectionInfo {
	collection.limit = n
	return collection
}

// Skip 跳过条数
func (collection *CollectionInfo) Skip(n int64) *CollectionInfo {
	collection.skip = n
	return collection
}

// Sort 排序 bson.M{"created_at":-1}
func (collection *CollectionInfo) Sort(sorts bson.M) *CollectionInfo {
	collection.sort = sorts
	return collection
}

// Fields 指定查询字段
func (collection *CollectionInfo) Fields(fields bson.M) *CollectionInfo {
	collection.fields = fields
	return collection
}

// InsertOne 写入单条数据
func (collection *CollectionInfo) InsertOne(document interface{}) *mongo.InsertOneResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.InsertOne(ctx, BeforeCreate(document))
	if err != nil {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error(err)
	}
	return result
}

// InsertMany 写入多条数据
func (collection *CollectionInfo) InsertMany(documents interface{}) *mongo.InsertManyResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var data []interface{}
	data = BeforeCreate(documents).([]interface{})
	result, err := collection.Table.InsertMany(ctx, data)
	if err != nil {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error(err)
	}
	return result
}

// UpdateOrInsert 存在更新,不存在写入, documents 里边的文档需要有 _id 的存在
func (collection *CollectionInfo) UpdateOrInsert(documents interface{}) *mongo.UpdateResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var upsert = true
	result, err := collection.Table.UpdateMany(ctx, collection.filter, documents, &options.UpdateOptions{Upsert: &upsert})
	if err != nil {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error(err)
	}
	return result
}

// UpdateOne 更新一条
func (collection *CollectionInfo) UpdateOne(document interface{}) *mongo.UpdateResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.UpdateOne(ctx, collection.filter, bson.M{"$set": BeforeUpdate(document)})
	if err != nil {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error(err)
	}
	return result
}

// UpdateMany 更新多条
func (collection *CollectionInfo) UpdateMany(document interface{}) *mongo.UpdateResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.UpdateMany(ctx, collection.filter, bson.M{"$set": BeforeUpdate(document)})
	if err != nil {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error(err)
	}
	return result
}

// FindOne 查询一条数据
func (collection *CollectionInfo) FindOne(document interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result := collection.Table.FindOne(ctx, collection.filter, &options.FindOneOptions{
		Skip:       &collection.skip,
		Sort:       collection.sort,
		Projection: collection.fields,
	})
	err := result.Decode(document)
	if err != nil {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error(err)
		return err
	}
	return nil
}

// FindMany 查询多条数据
func (collection *CollectionInfo) FindMany(documents interface{}) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.Find(ctx, collection.filter, &options.FindOptions{
		Skip:       &collection.skip,
		Limit:      &collection.limit,
		Sort:       collection.sort,
		Projection: collection.fields,
	})
	if err != nil {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error(err)
	}
	defer result.Close(ctx)
	err = result.All(context.Background(), documents)
	if err != nil {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error(err)
	}
}

// Delete 删除数据,并返回删除成功的数量
func (collection *CollectionInfo) Delete() int64 {
	if collection.filter == nil || len(collection.filter) == 0 {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error("you can't delete all documents, it's very dangerous")
		return 0
	}
	if collection.SafeDelete {
		rs := collection.UpdateMany(bson.M{"deleted_at": time.Now().Unix()})
		return rs.ModifiedCount
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.DeleteMany(ctx, collection.filter)
	if err != nil {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error(err)
	}
	return result.DeletedCount
}

// Count 根据指定条件获取总条数
func (collection *CollectionInfo) Count() int64 {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.CountDocuments(ctx, collection.filter)
	if err != nil {
		app.Logger().WithField("log_type", "pkg.mongo.mongo").Error(err)
		return 0
	}
	return result
}

// BeforeCreate 创建数据前置操作
func BeforeCreate(document interface{}) interface{} {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)

	switch typ.Kind() {
	case reflect.Ptr:
		return BeforeCreate(val.Elem().Interface())

	case reflect.Array, reflect.Slice:
		var sliceData = make([]interface{}, val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			sliceData[i] = BeforeCreate(val.Index(i).Interface()).(bson.M)
		}
		return sliceData

	case reflect.Struct:
		var data = make(bson.M)
		for i := 0; i < typ.NumField(); i++ {
			tag := typ.Field(i).Tag.Get("bson")
			if val.Field(i).Kind() == reflect.Struct && tag == "#expand" {
				b := BeforeCreate(val.Field(i).Interface())
				for key, v := range b.(bson.M) {
					data[key] = v
				}
			}
			if tag == "" || tag == "-" || tag == "#expand" {
				continue
			}
			data[tag] = val.Field(i).Interface()
		}
		if val.FieldByName("ID").Type() == reflect.TypeOf(primitive.ObjectID{}) {
			data["_id"] = primitive.NewObjectID()
		}

		if val.FieldByName("ID").Kind() == reflect.String && val.FieldByName("ID").Interface() == "" {
			data["_id"] = primitive.NewObjectID().Hex()
		}

		if IsIntn(val.FieldByName("ID").Kind()) && val.FieldByName("ID").Interface() == 0 {
			data["_id"] = unique.ID()
		}
		now := time.Now().Unix()
		if v, ok := data["created_at"]; !ok || v.(int64) == 0 {
			data["created_at"] = now
		}
		data["updated_at"] = now
		data["deleted_at"] = 0
		return data

	default:
		if val.Type() == reflect.TypeOf(bson.M{}) {
			if !val.MapIndex(reflect.ValueOf("_id")).IsValid() {
				val.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(primitive.NewObjectID()))
			}
			val.SetMapIndex(reflect.ValueOf("created_at"), reflect.ValueOf(time.Now().Unix()))
			val.SetMapIndex(reflect.ValueOf("updated_at"), reflect.ValueOf(time.Now().Unix()))
		}
		return val.Interface()
	}
}

// BeforeUpdate 更新数据前置操作
func BeforeUpdate(document interface{}) interface{} {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)
	switch typ.Kind() {
	case reflect.Ptr:
		return BeforeUpdate(val.Elem().Interface())

	case reflect.Array, reflect.Slice:
		var sliceData = make([]interface{}, val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			sliceData[i] = BeforeUpdate(val.Index(i).Interface()).(bson.M)
		}
		return sliceData

	case reflect.Struct:
		var data = make(bson.M)
		for i := 0; i < typ.NumField(); i++ {
			tag := strings.Split(typ.Field(i).Tag.Get("bson"), ",")[0]
			if val.Field(i).Kind() == reflect.Struct && tag == "#expand" {
				b := BeforeUpdate(val.Field(i).Interface())
				for key, v := range b.(bson.M) {
					data[key] = v
				}
			}
			if tag == "" || tag == "-" || tag == "#expand" {
				continue
			}
			if !isZero(val.Field(i)) {
				if tag != "_id" {
					data[tag] = val.Field(i).Interface()
				}
			}
		}
		data["updated_at"] = time.Now().Unix()
		return data

	default:
		if val.Type() == reflect.TypeOf(bson.M{}) {
			val.SetMapIndex(reflect.ValueOf("updated_at"), reflect.ValueOf(time.Now().Unix()))
		}
		return val.Interface()
	}
}

// IsIntn 是否为整数
func IsIntn(p reflect.Kind) bool {
	return p == reflect.Int || p == reflect.Int64 || p == reflect.Uint64 || p == reflect.Uint32
}

func isZero(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

// CreateManyIndex 多字段创建index索引
func (collection *CollectionInfo) CreateManyIndex(keys map[string]interface{}) (err error) {
	ctx := context.Background()
	indexView := collection.Table.Indexes()
	indexModels := make([]mongo.IndexModel, len(keys))
	j := 0
	for i, v := range keys {
		key := map[string]interface{}{i: v}
		indexModels[j] = mongo.IndexModel{
			Keys: key,
		}
		j++
	}
	_, err = indexView.CreateMany(ctx, indexModels)
	if err != nil {
		return
	}
	return
}

// CreateIndexes 创建非unique索引
func (collection *CollectionInfo) CreateIndexes(keys []map[string]interface{}) (err error) {
	ctx := context.Background()
	indexView := collection.Table.Indexes()
	indexModels := make([]mongo.IndexModel, 0)
	for _, v := range keys {
		option := options.Index()
		item := mongo.IndexModel{
			Keys:    v,
			Options: option,
		}
		indexModels = append(indexModels, item)
	}
	fmt.Println(indexModels)
	_, err = indexView.CreateMany(ctx, indexModels)
	return
}

// CreateUniqueIndex 创建唯一索引
func (collection *CollectionInfo) CreateUniqueIndex(keys map[string]interface{}) error {
	ctx := context.Background()
	unique := true
	indexView := collection.Table.Indexes()
	option := options.Index()
	option.Unique = &unique
	indexModel := mongo.IndexModel{Keys: keys, Options: option}
	_, err := indexView.CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}
	return nil
}

// CreateUniqueIndexes 批量创建唯一索引
func (collection *CollectionInfo) CreateUniqueIndexes(keys []map[string]interface{}) error {
	ctx := context.Background()
	unique := true
	indexView := collection.Table.Indexes()
	indexModels := make([]mongo.IndexModel, 0)
	for _, v := range keys {
		option := options.Index()
		option.Unique = &unique
		item := mongo.IndexModel{Keys: v, Options: option}
		indexModels = append(indexModels, item)
	}
	_, err := indexView.CreateMany(ctx, indexModels)
	if err != nil {
		return err
	}
	return nil
}

// NewTx 新建事务
func NewTx() *Transaction {
	tx := Transaction{}
	return &tx
}

// Begin 开启事务, mongo 如果要支持事务需要部署为集群
func (t *Transaction) Begin(fn func(tx *Transaction) (err error)) error {
	err := client.UseSession(context.Background(), func(sessionCtx mongo.SessionContext) (err error) {
		t.Session = sessionCtx
		err = sessionCtx.StartTransaction()
		if err != nil {
			return err
		}
		return fn(t)
	})
	return err
}

// Commit 提交事务
func (t *Transaction) Commit() error {
	return t.Session.CommitTransaction(t.Session)
}

// Rollback 回滚事务
func (t *Transaction) Rollback() error {
	return t.Session.AbortTransaction(t.Session)
}

// Collection 得到一个mongo操作对象
func (t *Transaction) Collection(table string) *CollectionInfo {
	db := client.Database(conf.Database)
	return &CollectionInfo{
		Database: db,
		Table:    db.Collection(table),
		Context:  t.Session,
		filter:   make(bson.M),
	}
}

// InsertOneWithError 支持事务
func (collection *CollectionInfo) InsertOneWithError(document interface{}) (result *mongo.InsertOneResult, err error) {
	result, err = collection.Table.InsertOne(collection.Context, BeforeCreate(document))
	if err != nil {
		return
	}
	return
}

//InsertManyWithError 支持事务
func (collection *CollectionInfo) InsertManyWithError(documents interface{}) (rsp *mongo.InsertManyResult, err error) {
	var data []interface{}
	data = BeforeCreate(documents).([]interface{})
	rsp, err = collection.Table.InsertMany(collection.Context, data)
	if err != nil {
		return
	}
	return
}

// AggregateWithError 聚合操作, 返回错误信息
func (collection *CollectionInfo) AggregateWithError(pipeline interface{}, result interface{}) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := collection.Table.Aggregate(ctx, pipeline)
	if err != nil {
		return
	}
	err = cursor.All(ctx, result)
	if err != nil {
		return
	}
	return
}

// UpdateOneWithError 支持事务
func (collection *CollectionInfo) UpdateOneWithError(document interface{}) (result *mongo.UpdateResult, err error) {
	result, err = collection.Table.UpdateOne(collection.Context, collection.filter, bson.M{"$set": BeforeUpdate(document)})
	if err != nil {
		return
	}
	return
}

// UpdateOneRawWithError 原生update
// 支持事务
func (collection *CollectionInfo) UpdateOneRawWithError(document interface{}, opt ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	result, err := collection.Table.UpdateOne(collection.Context, collection.filter, document, opt...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateManyWithError 更新多条数据, 并返回错误信息
func (collection *CollectionInfo) UpdateManyWithError(document interface{}) (result *mongo.UpdateResult, err error) {
	result, err = collection.Table.UpdateMany(collection.Context, collection.filter, bson.M{"$set": BeforeUpdate(document)})
	if err != nil {
		return
	}
	return
}

// FindManyWithError 查询多条数据，将错误外抛出
func (collection *CollectionInfo) FindManyWithError(documents interface{}) (err error) {
	result, err := collection.Table.Find(collection.Context, collection.filter, &options.FindOptions{
		Skip:       &collection.skip,
		Limit:      &collection.limit,
		Sort:       collection.sort,
		Projection: collection.fields,
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer result.Close(collection.Context)

	val := reflect.ValueOf(documents)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		err = errors.New("result argument must be a slice address")
		return
	}

	slice := reflect.MakeSlice(val.Elem().Type(), 0, 0)
	itemTyp := val.Elem().Type().Elem()
	for result.Next(collection.Context) {
		item := reflect.New(itemTyp)
		err := result.Decode(item.Interface())
		if err != nil {
			err = errors.New("result argument must be a slice address")
			return err
		}
		slice = reflect.Append(slice, reflect.Indirect(item))
	}
	val.Elem().Set(slice)
	return
}

// DeleteWithError 支持事务
func (collection *CollectionInfo) DeleteWithError() (count int64, err error) {
	if collection.filter == nil || len(collection.filter) == 0 {
		err = errors.New("you can't delete all documents, it's very dangerous")
		return
	}
	if collection.SafeDelete {
		rs, err := collection.UpdateManyWithError(bson.M{"deleted_at": time.Now().Unix()})
		return rs.ModifiedCount, err
	}
	result, err := collection.Table.DeleteMany(collection.Context, collection.filter)
	if err != nil {
		return
	}
	count = result.DeletedCount
	return
}

// CountWithError 计算查询条数, 并同时返回错误信息
func (collection *CollectionInfo) CountWithError() (result int64, err error) {
	result, err = collection.Table.CountDocuments(collection.Context, collection.filter)
	if err != nil {
		return
	}
	return
}

// StFindMany 用于结构体嵌套查询
func (collection *CollectionInfo) StFindMany(doc interface{}) error {
	var d []interface{}
	err := collection.FindManyWithError(&d)
	if err != nil {
		return err
	}
	err = find(d, doc)
	return err
}

// StFindOne 用于结构体嵌套查询
func (collection *CollectionInfo) StFindOne(doc interface{}) error {
	var d interface{}
	err := collection.FindOne(&d)
	if err != nil {
		return err
	}
	err = find(d, doc)
	return err
}

func find(document, setDoc interface{}) error {
	val := reflect.ValueOf(document)
	dtyp := reflect.TypeOf(document)
	typ := reflect.TypeOf(setDoc)
	sval := reflect.ValueOf(setDoc).Elem()
	if typ.Kind() != reflect.Ptr {
		return errors.New("setDoc must ptr")
	}
	if typ.Elem().Kind() != reflect.Array && typ.Elem().Kind() != reflect.Slice {
		all(document, setDoc)
		return nil
	}
	switch dtyp.Kind() {
	case reflect.Array, reflect.Slice:
		x := reflect.New(typ.Elem().Elem())
		v := x.Interface()
		for i := 0; i < val.Len(); i++ {
			st := all(val.Index(i).Interface(), v)
			n := reflect.Append(sval, reflect.ValueOf(st))
			sval.Set(n)
		}
	}
	return nil
}

func all(document, setDoc interface{}) interface{} {
	val := reflect.ValueOf(setDoc)
	typ := reflect.TypeOf(setDoc)
	dtyp := reflect.TypeOf(document)
	dval := reflect.ValueOf(document)
	if typ.Kind() != reflect.Ptr {
		return setDoc
	}
	var temp = make(map[string]reflect.Value)
	if dtyp.Kind() == reflect.Slice || dtyp.Kind() == reflect.Array {
		for i := 0; i < dval.Len(); i++ {
			str := dval.Index(i)
			temp[str.Field(0).String()] = str.Field(1)
		}
	}
	ty := val.Elem().Type()
	switch ty.Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int64, reflect.Int32, reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		val.Elem().Set(reflect.ValueOf(document))
	case reflect.Struct:
		for i := 0; i < ty.NumField(); i++ {
			tag := ty.Field(i).Tag.Get("bson")
			y := val.Elem().FieldByName(ty.Field(i).Name)
			if tag == "#expand" {
				st := all(document, reflect.New(y.Type()).Interface())
				y.Set(reflect.ValueOf(st).Convert(y.Type()))
			}
			if tag == "" || tag == "-" {
				continue
			}

			if v, ok := temp[tag]; ok {
				if !v.IsNil() {
					y.Set(handleD(v.Interface(), y.Type()))
				}
			}

		}
	}
	return val.Elem().Interface()
}

// 转化为指定类型的数据格式
func handleD(doc interface{}, kind reflect.Type) reflect.Value {
	var res reflect.Value
	if doc == nil {
		return res
	}
	nvl := reflect.ValueOf(doc)
	sk := reflect.TypeOf(doc).Kind()
	switch {
	case kind.Kind() == reflect.Slice && sk == reflect.Slice:
		x := reflect.New(kind.Elem())
		vs := x.Interface()
		dr := make([]reflect.Value, 0)
		ab := reflect.MakeSlice(kind, 0, 0)
		for i := 0; i < nvl.Len(); i++ {
			st := all(nvl.Index(i).Elem().Interface(), vs)
			dr = append(dr, reflect.ValueOf(st))
		}
		v1 := reflect.Append(ab, dr...)
		res = v1
	default:
		res = nvl.Convert(kind)
	}
	return res
}
