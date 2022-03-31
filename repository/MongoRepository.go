package repository

import (
	"context"
	"errors"
	datacommon "github.com/aomi-go/data-common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

type MongoRepository struct {
	// Database 数据库数据源信息
	Database *mongo.Database
	// EntityType 数据对象结构体类型
	EntityType      reflect.Type
	EntitySliceType reflect.Type
	// EntityType 对应的数据库集合名称,可选参数，默认使用结构体类型名称作为集合名称
	CollectionName string
	Collection     *mongo.Collection
}

// NewMongoRepo 创建一个MongoRepository 实例
func NewMongoRepo(database *mongo.Database, entityType reflect.Type, collectionName string) MongoRepository {
	if nil == database {
		panic(errors.New("mongo database 不能为 nil").(any))
	}
	if nil == entityType {
		panic(errors.New("实体类型不能为 nil").(any))
	}

	if "" == collectionName {
		collectionName = entityType.Name()
		collectionName = strings.ToLower(string(collectionName[0])) + collectionName[1:]
	}

	collection := database.Collection(collectionName)

	repo := MongoRepository{
		Database:        database,
		EntityType:      entityType,
		EntitySliceType: reflect.SliceOf(entityType),
		CollectionName:  collectionName,
		Collection:      collection,
	}

	return repo
}

func (d *MongoRepository) Save(entity interface{}) interface{} {
	_, err := d.Collection.InsertOne(context.TODO(), entity)
	if err != nil {
		panic(err.(any))
	}
	return entity
}

func (d *MongoRepository) FindById(id interface{}) interface{} {
	result := d.Collection.FindOne(context.TODO(), bson.D{{
		"_id", id,
	}})
	if nil != result.Err() {
		panic(result.Err().(any))
	}
	item := reflect.New(d.EntityType).Interface()
	err := result.Decode(item)
	if err != nil {
		panic(err.(any))
	}
	return &item
}

func (d *MongoRepository) ExistsById(id *interface{}) bool {
	count, err := d.Collection.CountDocuments(context.TODO(), bson.D{{
		"_id", id,
	}})
	if err != nil {
		panic(err.(any))
	}
	return count > 0
}

func (d *MongoRepository) DeleteById(id interface{}) bool {
	result, err := d.Collection.DeleteOne(context.TODO(), bson.D{{
		"_id", id,
	}})
	if err != nil {
		panic(err.(any))
	}
	return result.DeletedCount > 0
}

func (d *MongoRepository) FindOne(filter interface{}, opts ...*options.FindOneOptions) interface{} {
	item := reflect.New(d.EntityType)
	err := d.Collection.FindOne(context.TODO(), filter, opts...).Decode(item.Interface())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil
		}
		panic(err.(any))
	}
	return item
}

func (d *MongoRepository) Find(filter interface{}, opts ...*options.FindOptions) interface{} {
	cursor, err := d.Collection.Find(context.TODO(), filter, opts...)
	if nil != err {
		panic(err.(any))
	}
	return d.getContent(cursor, err)
}

//FindAll 查询所有带分页排序
func (d *MongoRepository) FindAll(filter interface{}, pageable datacommon.Pageable, opts ...*options.FindOptions) *datacommon.Page {

	finalOpts := opts
	var totalElements int64 = 0

	if nil != pageable {
		skip := pageable.GetOffset()
		limit := pageable.GetPageSize()

		finalOpts = append(opts, &options.FindOptions{
			Skip:  &skip,
			Limit: &limit,
		})

		t, err := d.Collection.CountDocuments(context.TODO(), filter)

		if nil != err {
			panic(err.(any))
		}
		totalElements = t
	}

	cursor, err := d.Collection.Find(context.TODO(), filter, finalOpts...)
	content := d.getContent(cursor, err)
	return datacommon.Of(content, pageable, totalElements)
}

func (repo *MongoRepository) getContent(cursor *mongo.Cursor, err error) interface{} {
	if nil != err {
		panic(err.(any))
	}

	var contentValue = reflect.MakeSlice(repo.EntitySliceType, 0, 0)

	for cursor.Next(context.TODO()) {
		item := reflect.New(repo.EntityType)
		e := cursor.Decode(item.Interface())
		if e != nil {
			panic(e.(any))
		}
		contentValue = reflect.Append(contentValue, item.Elem())
	}
	return contentValue.Interface()
}
