package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *MongoDB

// MongoDB structure that holds the MongoDB client and collection name
type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// GetMongoDB initializes and returns a MongoDB client and collection
func GetMongoDB(client *mongo.Client, databaseName, collectionName string) *MongoDB {
	collection := client.Database(databaseName).Collection(collectionName)
	db = &MongoDB{client, collection}
	return db
}

// GetItem retrieves a single item from MongoDB
func (m *MongoDB) GetItem(ctx context.Context, keyname string, keyvalue interface{}) (interface{}, error) {
	filter := bson.M{keyname: keyvalue}
	var result bson.M
	err := m.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetBatchItem retrieves multiple items from MongoDB based on a list of key values
func (m *MongoDB) GetBatchItem(ctx context.Context, keyname string, keyvalues []string) (interface{}, error) {
	filter := bson.M{keyname: bson.M{"$in": keyvalues}}
	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	for cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// PutItem inserts a single item into MongoDB
func (m *MongoDB) PutItem(ctx context.Context, value interface{}) error {
	_, err := m.collection.InsertOne(ctx, value)
	if err != nil {
		return err
	}
	return nil
}

// BatchPutItem inserts multiple items into MongoDB
func (m *MongoDB) BatchPutItem(ctx context.Context, values []interface{}) error {
	_, err := m.collection.InsertMany(ctx, values)
	if err != nil {
		return err
	}
	return nil
}

// GetAllItemsWithGSI retrieves all items based on a Global Secondary Index (GSI) from MongoDB
// Note: MongoDB doesn't directly use GSIs like DynamoDB. This is a placeholder method.
func (m *MongoDB) GetAllItemsWithGSI(ctx context.Context, partitionKey string, partitionValue []string) (interface{}, error) {
	// In MongoDB, you can use an index to achieve similar functionality to DynamoDB's GSI.
	filter := bson.M{partitionKey: bson.M{"$in": partitionValue}}
	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	for cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// UpdateItem updates an existing item in MongoDB
func (m *MongoDB) UpdateItem(ctx context.Context, keyMap map[string]interface{}, updateKeys interface{}) error {
	filter := bson.M{}
	for key, value := range keyMap {
		filter[key] = value
	}

	update := bson.M{
		"$set": updateKeys,
	}

	_, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func Get() *MongoDB {
	return db
}
