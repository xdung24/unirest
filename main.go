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
	// db driver
	MEMORY = "memory"   // fastest (hashmap on RAM)
	FS     = "fs"       // depend on disk speed, second fastest on ssd nvme
	SQLITE = "sqlite"   // local disk relational database (sql)
	PG     = "postgres" // relational database (sql)
	MYSQL  = "mysql"    // relational database (sql)
	REDIS  = "redis"    // keypair value database (nosql)
	MONGO  = "mongo"    // bson document database (nosql)

	// env
	envHostPort       = "IP_PORT"
	envDbDriver       = "DB_DRIVER"
	envDbHost         = "DB_HOST"
	envDbName         = "DB_NAME"
	envDbUser         = "DB_USER"
	envDbPass         = "DB_PASS"
	envDbPath         = "DB_PATH"
	envBrokerEnabled  = "BROKER_ENABLED"
	envBrokerHostPort = "BROKER_IP_PORT"
	envSwaggerEnabled = "SWAGGER_ENABLED"
	envAuthEnabled    = "AUTH_ENABLED"
	envRawSqlEnabled  = "RAW_SQL_ENABLED"
)

func main() {
	printInfo()

	var addr, dbDriver, dbHost, dbName, dbUser, dbPass, dbPath, brokerHostPort string
	var swaggerEnabled, brokerEnabled, authEnabled, rawSqlEnabled bool

	flag.StringVar(&addr, envHostPort, "0.0.0.0:8000", "ip:port for rest api to expose")
	flag.StringVar(&brokerHostPort, envBrokerHostPort, "0.0.0.0:8001", "ip:port for broker to expose")

	flag.BoolVar(&swaggerEnabled, envSwaggerEnabled, false, "enable swagger")
	flag.BoolVar(&brokerEnabled, envBrokerEnabled, false, "enable broker")
	flag.BoolVar(&authEnabled, envAuthEnabled, false, "enable JWT auth")
	flag.BoolVar(&rawSqlEnabled, envRawSqlEnabled, false, "enable raw sql (postgres)")

	flag.StringVar(&dbDriver, envDbDriver, MEMORY, "db type to use (memory | fs | sqlite| postgres | mysql | redis | mongo)")
	flag.StringVar(&dbPath, envDbPath, "./data", "path of the file storage (for fs | sqlite)")
	flag.StringVar(&dbHost, envDbHost, "localhost", "database host (for postgres | mysql | redis | mongo)")
	flag.StringVar(&dbName, envDbName, "", "database name (for postgres | mysql | mongo)")
	flag.StringVar(&dbUser, envDbUser, "", "database user (for postgres | mysql | mongo)")
	flag.StringVar(&dbPass, envDbPass, "", "database password (for postgres | mysql | mongo)")

	flag.Parse()

	server := service.Server{
		Address:        addr,
		BrokerAddress:  brokerHostPort,
		SwaggerEnabled: swaggerEnabled,
		BrokerEnabled:  brokerEnabled,
		AuthEnabled:    authEnabled,
		RawSqlEnabled:  rawSqlEnabled,
	}

	var db service.Database
	switch dbDriver {
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
	log.Println("db type: ", dbDriver)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	db.Disconnect()

	log.Println("Good bye")
}
