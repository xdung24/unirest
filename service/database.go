package service

import "github.com/xdung24/unirest/database"

type Database interface {
	Init()
	Disconnect()
	CreateNameSpace(namespace string) *database.DbError
	GetNamespaces() []string
	DropNameSpace(namespace string) *database.DbError
	Upsert(namespace string, key string, value []byte, allowOverWrite bool) *database.DbError
	Get(namespace string, key string) ([]byte, *database.DbError)
	GetAll(namespace string) (map[string][]byte, *database.DbError)
	Delete(namespace string, key string) *database.DbError
	DeleteAll(namespace string) *database.DbError
}
