package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysql_insertQuery        = "INSERT INTO %v (id, data) VALUES(?, ?) ON DUPLICATE KEY UPDATE data = ?"
	mysql_tablesQuery        = "SELECT table_name FROM information_schema.tables WHERE table_schema = '%v'"
	mysql_getQuery           = "SELECT data FROM %v WHERE id = ?"
	mysql_getAllQuery        = "SELECT id, data FROM %v ORDER BY id"
	mysql_deleteQuery        = "DELETE FROM %v WHERE id = ?"
	mysql_dropNamespaceQuery = "DROP TABLE %v"
	mysql_createTableQuery   = "CREATE TABLE IF NOT EXISTS %v (id VARCHAR(14) NOT NULL, data json NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB;"
)

type MySqlDatabase struct {
	Host string
	User string
	Pass string
	Name string

	db *sql.DB
}

func (p *MySqlDatabase) Init() {
	connInfo := fmt.Sprintf("%v:%v@tcp(%v)/%v", p.User, p.Pass, p.Host, p.Name)
	db, err := sql.Open("mysql", connInfo)

	if err != nil {
		log.Println(connInfo)
		log.Fatalf("error connecting to mysql: %v", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	p.db = db
}

func (p MySqlDatabase) Upsert(namespace string, key string, value []byte) *DbError {
	err := p.ensureNamespace(namespace)

	if err != nil {
		return &DbError{
			ErrorCode: NAMESPACE_NOT_FOUND,
			Message:   fmt.Sprintf("namespace %v does not exist", namespace),
		}
	}
	_, dbErr := p.db.Exec(fmt.Sprintf(mysql_insertQuery, namespace), key, string(value), string(value))
	if dbErr != nil {
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   fmt.Sprintf("error on Upsert: %v", dbErr),
		}
	}
	return nil
}

func (p MySqlDatabase) Get(namespace string, key string) ([]byte, *DbError) {
	rows, dbErr := p.db.Query(fmt.Sprintf(mysql_getQuery, namespace), key)
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

func (p MySqlDatabase) GetAll(namespace string) (map[string][]byte, *DbError) {
	sqlStatement := fmt.Sprintf(mysql_getAllQuery, namespace)
	rows, dbErr := p.db.Query(sqlStatement)
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

func (p MySqlDatabase) Delete(namespace string, key string) *DbError {
	sqlStatement := fmt.Sprintf(mysql_deleteQuery, namespace)
	_, err := p.db.Exec(sqlStatement, key)
	if err != nil {
		log.Println(sqlStatement)
		message := fmt.Sprintf("error on Delete: %v", err)
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   message,
		}
	}
	return nil
}

func (p MySqlDatabase) DeleteAll(namespace string) *DbError {
	sqlStatement := fmt.Sprintf(mysql_dropNamespaceQuery, namespace)
	_, err := p.db.Exec(sqlStatement)
	if err != nil {
		log.Println(sqlStatement)
		message := fmt.Sprintf("error on DeleteAll: %v", err)
		return &DbError{
			ErrorCode: INTERNAL_ERROR,
			Message:   message,
		}
	}
	return nil
}

func (p MySqlDatabase) GetNamespaces() []string {
	sqlStatement := fmt.Sprintf(mysql_tablesQuery, p.Name)
	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		log.Printf("error on GetNamespaces: %v\n", err)
	}
	defer rows.Close()

	ret := make([]string, 0)
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			log.Println(sqlStatement)
			log.Printf("error on Scan: %v\n", err)
		}
		ret = append(ret, tableName)
	}
	return ret
}

func (p MySqlDatabase) ensureNamespace(namespace string) (err error) {
	query := fmt.Sprintf(mysql_createTableQuery, namespace)
	_, err = p.db.Exec(query)

	if err != nil {
		log.Println(query)
		log.Printf("error creating table: %v\n", err)
	}

	return err
}
