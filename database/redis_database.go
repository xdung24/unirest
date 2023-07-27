package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	redis "github.com/redis/go-redis/v9"
)

const (
	namespace_prefix = "caffein_"
	schema_suffix    = "_schema"
)

type RedisDatabase struct {
	Host string

	db *redis.Client
}

var ctx = context.Background()

func (p *RedisDatabase) Init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     p.Host,
		Password: "", // no password set
		DB:       0,  // use default DB
		Protocol: 3,  // specify 2 for RESP 2 or 3 for RESP 3
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		log.Fatalf("error connecting to mysql: %v", err)
	}

	p.db = rdb
}

func (p *RedisDatabase) Upsert(namespace string, key string, value []byte) *DbError {
	_, err := p.db.HSet(ctx, namespace_prefix+namespace, key, string(value)).Result()
	if err != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on Upsert: %v", err),
		}
	}
	return nil
}

func (p *RedisDatabase) Get(namespace string, key string) ([]byte, *DbError) {
	val, err := p.db.HGet(ctx, namespace_prefix+namespace, key).Result()
	if err == redis.Nil {
		return nil, &DbError{
			ErrorCode: ID_NOT_FOUND,
			Message:   fmt.Sprintf("value not found in namespace %v for key %v", namespace, key),
		}
	}
	if err != nil {
		return nil, &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on Get: %v", err),
		}
	}

	return []byte(val), nil
}

func (p *RedisDatabase) GetAll(namespace string) (map[string][]byte, *DbError) {
	val, err := p.db.HGetAll(ctx, namespace_prefix+namespace).Result()
	if err != nil {
		return nil, &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on GetAll: %v", err),
		}
	}

	ret := make(map[string][]byte)

	for i, v := range val {
		ret[i] = []byte(v)
	}

	return ret, nil
}

func (p *RedisDatabase) Delete(namespace string, key string) *DbError {
	_, err := p.db.HDel(ctx, namespace_prefix+namespace, key).Result()
	if err != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on Delete: %v", err),
		}
	}
	return nil
}

func (p *RedisDatabase) DeleteAll(namespace string) *DbError {
	val, err := p.db.HGetAll(ctx, namespace_prefix+namespace).Result()
	if err != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on HGetAll: %v", err),
		}
	}
	for i := range val {
		_, err := p.db.HDel(ctx, namespace_prefix+namespace, i).Result()
		if err != nil {
			return &DbError{
				ErrorCode: INTERNAL_ERROR,
				Message:   fmt.Sprintf("error on HDel: %v", err),
			}
		}
	}
	return nil
}

func (p *RedisDatabase) GetNamespaces() []string {
	var ret = []string{}
	val, _, err := p.db.Scan(ctx, 0, namespace_prefix+"*", 0).Result()
	if err != nil {
		return ret
	}
	for _, v := range val {
		if !strings.HasSuffix(v, schema_suffix) {
			ret = append(ret, strings.Replace(v, namespace_prefix, "", 1))
		}
	}
	return ret
}
