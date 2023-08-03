package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/namsral/flag"

	"github.com/xdung24/universal-rest/database"
	"github.com/xdung24/universal-rest/service"
)

const (
	MEMORY = "memory"
	FS     = "fs"
	SQLITE = "sqlite"
	PG     = "postgres"
	MYSQL  = "mysql"
	REDIS  = "redis"
	MONGO  = "mongo"

	// env
	envHostPort       = "IP_PORT"
	envDbType         = "DB_TYPE"
	envDbHost         = "DB_HOST"
	envDbName         = "DB_NAME"
	envDbUser         = "DB_USER"
	envDbPass         = "DB_PASS"
	envDbPath         = "DB_PATH"
	envSwaggerEnabled = "SWAGGER_ENABLED"
	envBrokerEnabled  = "BROKER_ENABLED"
	envAuthEnabled    = "AUTH_ENABLED"
	envRawSqlEnabled  = "RAW_SQL_ENABLED"
)

func main() {
	var addr, dbType, dbHost, dbName, dbUser, dbPass, dbPath string
	var swaggerEnabled, brokerEnabled, authEnabled, rawSqlEnabled bool

	flag.StringVar(&addr, envHostPort, "0.0.0.0:8000", "ip:port to expose")

	flag.BoolVar(&swaggerEnabled, envSwaggerEnabled, false, "enable swagger")
	flag.BoolVar(&brokerEnabled, envBrokerEnabled, false, "enable broker")
	flag.BoolVar(&authEnabled, envAuthEnabled, false, "enable JWT auth")
	flag.BoolVar(&rawSqlEnabled, envRawSqlEnabled, false, "enable raw sql (postgres)")

	flag.StringVar(&dbType, envDbType, MEMORY, "db type to use, options: memory | fs | sqlite| postgres | mysql | redis | mongo")
	flag.StringVar(&dbPath, envDbPath, "./data", "path of the file storage (for fs | sqlite)")
	flag.StringVar(&dbHost, envDbHost, "localhost", "database host (for postgres | mysql | redis | mongo)")
	flag.StringVar(&dbName, envDbName, "", "database name (for postgres | mysql | mongo)")
	flag.StringVar(&dbUser, envDbUser, "", "database user (for postgres | mysql | mongo)")
	flag.StringVar(&dbPass, envDbPass, "", "database password (for postgres | mysql | mongo)")

	flag.Parse()

	server := service.Server{
		Address:        addr,
		SwaggerEnabled: swaggerEnabled,
		BrokerEnabled:  brokerEnabled,
		AuthEnabled:    authEnabled,
		RawSqlEnabled:  rawSqlEnabled,
	}

	var db service.Database
	switch dbType {
	case MEMORY:
		db = &database.MemDatabase{}
	case FS:
		db = &database.StorageDatabase{
			RootDirPath: dbPath,
		}
	case SQLITE:
		db = &database.SQLiteDatabase{
			DirPath: dbPath,
		}
	case PG:
		db = &database.PGDatabase{
			Host: dbHost,
			Name: dbName,
			User: dbUser,
			Pass: dbPass,
		}
	case MYSQL:
		db = &database.MySqlDatabase{
			Host: dbHost,
			Name: dbName,
			User: dbUser,
			Pass: dbPass,
		}
	case REDIS:
		db = &database.RedisDatabase{
			Host: dbHost,
		}
	case MONGO:
		db = &database.MongoDatabase{
			Host: dbHost,
			Name: dbName,
			User: dbUser,
			Pass: dbPass,
		}
	default:
		panic("invalid db type")
	}

	go server.Init(db)

	log.Println("server started at: ", server.Address)
	log.Println("db type: ", dbType)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	db.Disconnect()

	log.Println("bye")
}
