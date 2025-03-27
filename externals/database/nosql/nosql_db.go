package database_nosql

import "context"

type DBProvider interface {
	GetItem(ctx context.Context, keyname string, keyvalue interface{}) (interface{}, error)
	GetBatchItem(ctx context.Context, keyname string, keyvalue []string) (interface{}, error)
	PutItem(ctx context.Context, value interface{}) error
	BatchPutItem(ctx context.Context, values interface{}) error
	GetAllItemsWithGSI(ctx context.Context, partitionKey string, partitionValue []string) (interface{}, error)
	UpdateItem(ctx context.Context, keyMap map[string]interface{}, updateKeys interface{}) error
}
