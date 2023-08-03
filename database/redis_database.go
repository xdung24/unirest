package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
)

const (
	redis_namespace_prefix = "caffein_"
	redis_schema_suffix    = "_schema"
	redis_dbTimeout        = 10 * time.Second
)

type RedisDatabase struct {
	Host string

	db *redis.Client
}

func (r *RedisDatabase) Init() {
	ctx, cancel := context.WithTimeout(context.Background(), redis_dbTimeout)
	defer cancel()
	rdb := redis.NewClient(&redis.Options{
		Addr:     r.Host,
		Password: "", // no password set
		DB:       0,  // use default DB
		Protocol: 3,  // specify 2 for RESP 2 or 3 for RESP 3
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		log.Fatalf("error connecting to mysql: %v", err)
	}

	r.db = rdb
}

func (r *RedisDatabase) Disconnect() {
	err := r.db.Close()
	if err != nil {
		panic(err)
	}
	log.Println("diconnected")
}

func (r *RedisDatabase) Upsert(namespace string, key string, value []byte, allowOverWrite bool) *DbError {
	ctx, cancel := context.WithTimeout(context.Background(), redis_dbTimeout)
	defer cancel()
	_, err := r.db.HSet(ctx, redis_namespace_prefix+namespace, key, string(value)).Result()
	if err != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on Upsert: %v", err),
		}
	}
	return nil
}

func (r *RedisDatabase) Get(namespace string, key string) ([]byte, *DbError) {
	ctx, cancel := context.WithTimeout(context.Background(), redis_dbTimeout)
	defer cancel()
	val, err := r.db.HGet(ctx, redis_namespace_prefix+namespace, key).Result()
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

func (r *RedisDatabase) GetAll(namespace string) (map[string][]byte, *DbError) {
	ctx, cancel := context.WithTimeout(context.Background(), redis_dbTimeout)
	defer cancel()
	val, err := r.db.HGetAll(ctx, redis_namespace_prefix+namespace).Result()
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

func (r *RedisDatabase) Delete(namespace string, key string) *DbError {
	ctx, cancel := context.WithTimeout(context.Background(), redis_dbTimeout)
	defer cancel()
	_, err := r.db.HDel(ctx, redis_namespace_prefix+namespace, key).Result()
	if err != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on Delete: %v", err),
		}
	}
	return nil
}

func (r *RedisDatabase) DeleteAll(namespace string) *DbError {
	ctx, cancel := context.WithTimeout(context.Background(), redis_dbTimeout)
	defer cancel()
	val, err := r.db.HGetAll(ctx, redis_namespace_prefix+namespace).Result()
	if err != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on HGetAll: %v", err),
		}
	}
	for i := range val {
		_, err := r.db.HDel(ctx, redis_namespace_prefix+namespace, i).Result()
		if err != nil {
			return &DbError{
				ErrorCode: INTERNAL_ERROR,
				Message:   fmt.Sprintf("error on HDel: %v", err),
			}
		}
	}
	return nil
}

func (r *RedisDatabase) GetNamespaces() []string {
	ctx, cancel := context.WithTimeout(context.Background(), redis_dbTimeout)
	defer cancel()
	var ret = []string{}
	val, _, err := r.db.Scan(ctx, 0, redis_namespace_prefix+"*", 0).Result()
	if err != nil {
		return ret
	}
	for _, v := range val {
		if !strings.HasSuffix(v, redis_schema_suffix) {
			ret = append(ret, strings.Replace(v, redis_namespace_prefix, "", 1))
		}
	}
	return ret
}
