package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	Host string
	Name string
	User string
	Pass string

	client *mongo.Client
	db     *mongo.Database
}

const (
	mongo_dbTimeout = 10 * time.Second
)

func (m *MongoDatabase) Init() {
	ctx, cancel := context.WithTimeout(context.Background(), mongo_dbTimeout)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%v/?directConnection=true&serverSelectionTimeoutMS=2000", m.Host)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalf("error connecting to mongo: %v", err)
	}

	m.client = client
	db := client.Database(m.Name)
	if db == nil {
		log.Fatalf("error connecting to database: %v", m.Name)
	}
	m.db = db
	log.Println("db connected")
}

func (m *MongoDatabase) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), mongo_dbTimeout)
	defer cancel()

	err := m.client.Disconnect(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("diconnected")
}

func (m *MongoDatabase) GetNamespaces() []string {
	ctx, cancel := context.WithTimeout(context.Background(), mongo_dbTimeout)
	defer cancel()

	filter := bson.D{{}}
	names, err := m.db.ListCollectionNames(ctx, filter)
	if err != nil {
		log.Panicf("error on GetNamespaces: %v", err.Error())
		return []string{}
	}
	return names
}

func (m *MongoDatabase) DropNameSpace(namespace string) *DbError {
	ctx, cancel := context.WithTimeout(context.Background(), mongo_dbTimeout)
	defer cancel()

	err := m.db.Collection(namespace).Drop(ctx)
	if err != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on DeleteAll: %v", err),
		}
	}
	return nil
}

func (m *MongoDatabase) Upsert(namespace string, key string, value []byte, allowOverWrite bool) *DbError {
	ctx, cancel := context.WithTimeout(context.Background(), mongo_dbTimeout)
	defer cancel()

	err := m.ensureNamespace(namespace)
	if err != nil {
		return &DbError{
			ErrorCode: NAMESPACE_NOT_FOUND,
			Message:   fmt.Sprintf("namespace %v does not exist", namespace),
		}
	}

	coll := m.db.Collection(namespace)

	var bdoc interface{}
	err = bson.UnmarshalExtJSON(value, true, &bdoc)
	if err != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   err.Error(),
		}
	}

	filter := bson.D{{Key: "id", Value: key}}
	var update = bson.D{{Key: "$set", Value: bdoc}}
	if !allowOverWrite {
		res := coll.FindOne(ctx, filter)
		if res != nil && res.Err() == nil { // document exists
			return &DbError{
				ErrorCode: ITEM_CONFLICT,
				Message:   "item already exists",
			}
		}
	}

	opts := options.Update().SetUpsert(true)
	_, err = coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   err.Error(),
		}
	}

	return nil
}

func (m *MongoDatabase) Get(namespace string, key string) ([]byte, *DbError) {
	ctx, cancel := context.WithTimeout(context.Background(), mongo_dbTimeout)
	defer cancel()

	coll := m.db.Collection(namespace)

	filter := bson.D{
		{Key: "id", Value: key},
	}

	var document bson.M
	err := coll.FindOne(ctx, filter).Decode(&document)
	if err != nil {
		return nil, &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   err.Error(),
		}
	}

	delete(document, "_id") // delete _id
	delete(document, "id")  // delete id

	res, err := json.Marshal(document)
	if err != nil {
		return nil, &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   err.Error(),
		}
	}
	return res, nil
}

func (m *MongoDatabase) GetAll(namespace string) (map[string][]byte, *DbError) {
	ctx, cancel := context.WithTimeout(context.Background(), mongo_dbTimeout)
	defer cancel()

	coll := m.db.Collection(namespace)

	filter := bson.M{}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   err.Error(),
		}
	}

	ret := make(map[string][]byte)
	for cur.Next(ctx) {
		var result map[string]interface{}
		err := bson.Unmarshal(cur.Current, &result)
		if err != nil {
			return nil, &DbError{
				ErrorCode: INTERNAL_ERROR,
				Message:   err.Error(),
			}
		}

		delete(result, "_id") // delete _id
		id := fmt.Sprintf("%v", result["id"])
		delete(result, "id") // delete id

		data, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		ret[id] = data
	}

	return ret, nil
}

func (m *MongoDatabase) Delete(namespace string, key string) *DbError {
	ctx, cancel := context.WithTimeout(context.Background(), mongo_dbTimeout)
	defer cancel()

	filter := bson.D{{Key: "id", Value: key}}
	opts := options.Delete().SetHint(bson.D{{Key: "id", Value: 1}})
	_, err := m.db.Collection(namespace).DeleteOne(ctx, filter, opts)
	if err != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on Delete: %v", err),
		}
	}
	return nil
}

func (m *MongoDatabase) DeleteAll(namespace string) *DbError {
	ctx, cancel := context.WithTimeout(context.Background(), mongo_dbTimeout)
	defer cancel()

	_, err := m.db.Collection(namespace).DeleteMany(ctx, bson.D{{}})
	if err != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on DeleteAll: %v", err),
		}
	}
	return nil
}

func (m *MongoDatabase) ensureNamespace(namespace string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), mysql_dbTimeout)
	defer cancel()

	filter := bson.D{{Key: "name", Value: namespace}}
	names, err := m.db.ListCollectionNames(ctx, filter)
	if err != nil {
		log.Panicf("error on GetNamespaces: %v", err.Error())
		return err
	}
	if len(names) == 0 {
		err := m.db.CreateCollection(ctx, namespace)
		if err != nil {
			log.Printf("error creating collection: %v\n", err)
			return err
		}

		indexModel := mongo.IndexModel{
			Keys: bson.D{{Key: "id", Value: 1}},
		}
		name, err := m.db.Collection(namespace).Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			log.Printf("error creating index: %v\n", err)
			return err
		}
		log.Printf("Name of Index Created: %s\n", name)
	}

	return nil
}
