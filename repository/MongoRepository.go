package repository

import (
	"context"
	datacommon "github.com/aomi-go/data-common"
	"github.com/aomi-go/data-common/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

type MongoRepository struct {
	repository.CrudRepository

	Datasource *mongo.Database
	EntityType reflect.Type

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
		Datasource: opts.Datasource,
		EntityType: opts.EntityType,
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

	var content []interface{}

	for cursor.Next(context.TODO()) {
		item := reflect.New(d.EntityType).Interface()
		e := cursor.Decode(item)
		if e != nil {
			panic(e.(any))
		}
		content = append(content, item)
	}

	return datacommon.Of(&content, pageable, totalElements)
}
