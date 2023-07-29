package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	sqlite_insertQuery        = "INSERT INTO %v (id, data) VALUES($1, $2) ON CONFLICT (id) DO UPDATE SET data = $2"
	sqlite_tablesQuery        = "SELECT  `name` FROM sqlite_master WHERE `type`='table'  ORDER BY name"
	sqlite_getQuery           = "SELECT data FROM %v WHERE id = $1"
	sqlite_getAllQuery        = "SELECT id, data FROM %v ORDER BY id"
	sqlite_deleteQuery        = "DELETE FROM %v WHERE id = $1"
	sqlite_dropNamespaceQuery = "DROP TABLE %v"
	sqlite_createTableQuery   = "CREATE TABLE IF NOT EXISTS %v ( id string PRIMARY KEY, data string NOT NULL)"
)

type SQLiteDatabase struct {
	DirPath string
	db      *sql.DB
}

func (s *SQLiteDatabase) Init() {
	db, err := sql.Open("sqlite3", s.DirPath)
	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
	}
	s.db = db
}

func (s *SQLiteDatabase) Disconnect() {
	err := s.db.Close()
	if err != nil {
		panic(err)
	}
	log.Println("diconnected")
}

func (s SQLiteDatabase) Upsert(namespace string, key string, value []byte) *DbError {
	ctx, cancel := context.WithTimeout(context.Background(), pg_dbTimeout)
	defer cancel()
	err := s.ensureNamespace(namespace)

	if err != nil {
		return &DbError{
			ErrorCode: NAMESPACE_NOT_FOUND,
			Message:   fmt.Sprintf("namespace %v does not exist", namespace),
		}
	}
	_, dbErr := s.db.ExecContext(ctx, fmt.Sprintf(sqlite_insertQuery, namespace), key, string(value))
	if dbErr != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on Upsert: %v", dbErr),
		}
	}
	return nil
}

func (s SQLiteDatabase) Get(namespace string, key string) ([]byte, *DbError) {
	ctx, cancel := context.WithTimeout(context.Background(), pg_dbTimeout)
	defer cancel()
	rows, dbErr := s.db.QueryContext(ctx, fmt.Sprintf(sqlite_getQuery, namespace), key)
	if dbErr != nil {
		return nil, &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on Get: %v", dbErr),
		}
	}
	defer rows.Close()
	if rows.Next() {
		var data string
		scanErr := rows.Scan(&data)
		if scanErr != nil {
			return nil, &DbError{
				ErrorCode: INTERNAL_ERROR,
				Message:   fmt.Sprintf("scan %v", scanErr),
			}
		}
		return []byte(data), nil
	}
	return nil, &DbError{
		ErrorCode: ID_NOT_FOUND,
		Message:   fmt.Sprintf("value not found in namespace %v for key %v", namespace, key),
	}
}

func (s SQLiteDatabase) GetAll(namespace string) (map[string][]byte, *DbError) {
	ctx, cancel := context.WithTimeout(context.Background(), pg_dbTimeout)
	defer cancel()
	sqlStatement := fmt.Sprintf(sqlite_getAllQuery, namespace)
	rows, dbErr := s.db.QueryContext(ctx, sqlStatement)
	if dbErr != nil {
		return nil, &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on Get: %v", dbErr),
		}
	}
	defer rows.Close()

	ret := make(map[string][]byte)

	for rows.Next() {
		var id, data string
		scanErr := rows.Scan(&id, &data)
		if scanErr != nil {
			return nil, &DbError{
				ErrorCode: INTERNAL_ERROR,
				Message:   fmt.Sprintf("scan %v", scanErr),
			}
		}
		ret[id] = []byte(data)
	}
	return ret, nil
}

func (s SQLiteDatabase) Delete(namespace string, key string) *DbError {
	ctx, cancel := context.WithTimeout(context.Background(), pg_dbTimeout)
	defer cancel()
	_, err := s.db.ExecContext(ctx, fmt.Sprintf(sqlite_deleteQuery, namespace), key)
	if err != nil {
		message := fmt.Sprintf("error on Delete: %v", err)
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   message,
		}
	}
	return nil
}

func (s SQLiteDatabase) DeleteAll(namespace string) *DbError {
	ctx, cancel := context.WithTimeout(context.Background(), pg_dbTimeout)
	defer cancel()
	sqlStatement := fmt.Sprintf(sqlite_dropNamespaceQuery, namespace)
	_, err := s.db.ExecContext(ctx, sqlStatement)
	if err != nil {
		message := fmt.Sprintf("error on DeleteAll: %v", err)
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   message,
		}
	}
	return nil
}

func (p SQLiteDatabase) GetNamespaces() []string {
	ctx, cancel := context.WithTimeout(context.Background(), pg_dbTimeout)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, sqlite_tablesQuery)
	if err != nil {
		log.Printf("error on GetNamespaces: %v\n", err)
	}
	defer rows.Close()

	ret := make([]string, 0)
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			log.Printf("error on Scan: %v\n", err)
		}
		ret = append(ret, tableName)
	}
	return ret
}

func (p SQLiteDatabase) ensureNamespace(namespace string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), pg_dbTimeout)
	defer cancel()
	query := fmt.Sprintf(sqlite_createTableQuery, namespace)
	_, err = p.db.ExecContext(ctx, query)

	if err != nil {
		log.Printf("error creating table: %v\n", err)
	}

	return err
}
