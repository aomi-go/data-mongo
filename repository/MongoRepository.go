package repository

import (
	"context"
	datacommon "github.com/aomi-go/data-common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

type MongoRepository struct {
	Datasource      *mongo.Database
	EntityType      reflect.Type
	entitySliceType reflect.Type

	collection     *mongo.Collection
	collectionName string
}

type NewOptions struct {
	Datasource     *mongo.Database
	EntityType     reflect.Type
	CollectionName string
}

func NewMongoRepository(opts *NewOptions) *MongoRepository {
	repo := &MongoRepository{
		Datasource:      opts.Datasource,
		EntityType:      opts.EntityType,
		entitySliceType: reflect.SliceOf(opts.EntityType),
	}

	var cName string
	if "" == opts.CollectionName {
		tmpName := opts.EntityType.Name()
		cName = strings.ToLower(string(tmpName[0])) + tmpName[1:]
	} else {
		cName = opts.CollectionName
	}
	repo.collectionName = cName
	return repo
}

func (d MongoRepository) GetCollection() *mongo.Collection {
	if nil != d.collection {
		return d.collection
	}

	if "" != d.collectionName {
		d.collection = d.Datasource.Collection(d.collectionName)
	}
	return d.collection
}

func (d MongoRepository) Save(entity interface{}) interface{} {
	_, err := d.GetCollection().InsertOne(context.TODO(), entity)
	if err != nil {
		panic(err.(any))
	}
	return entity
}

func (d MongoRepository) FindById(id interface{}) interface{} {
	result := d.GetCollection().FindOne(context.TODO(), bson.D{{
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

func (d MongoRepository) ExistsById(id *interface{}) bool {
	count, err := d.GetCollection().CountDocuments(context.TODO(), bson.D{{
		"_id", id,
	}})
	if err != nil {
		panic(err.(any))
	}
	return count > 0
}

func (d MongoRepository) DeleteById(id interface{}) bool {
	result, err := d.GetCollection().DeleteOne(context.TODO(), bson.D{{
		"_id", id,
	}})
	if err != nil {
		panic(err.(any))
	}
	return result.DeletedCount > 0
}

//FindAll 查询所有带分页排序
func (d MongoRepository) FindAll(filter interface{}, pageable datacommon.Pageable, opts ...*options.FindOptions) *datacommon.Page {

	finalOpts := opts
	var totalElements int64 = 0

	if nil != pageable {
		skip := pageable.GetOffset()
		limit := pageable.GetPageSize()

		finalOpts = append(opts, &options.FindOptions{
			Skip:  &skip,
			Limit: &limit,
		})

		t, err := d.GetCollection().CountDocuments(context.TODO(), filter)

		if nil != err {
			panic(err.(any))
		}
		totalElements = t
	}

	cursor, err := d.GetCollection().Find(context.TODO(), filter, finalOpts...)
	if nil != err {
		panic(err.(any))
	}

	var contentValue = reflect.MakeSlice(d.entitySliceType, 0, 0)

	for cursor.Next(context.TODO()) {
		item := reflect.New(d.EntityType)
		e := cursor.Decode(item.Interface())
		if e != nil {
			panic(e.(any))
		}
		contentValue = reflect.Append(contentValue, item.Elem())
	}
	content := contentValue.Interface()

	return datacommon.Of(content, pageable, totalElements)
}
