package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/namsral/flag"

	"github.com/rehacktive/caffeine/database"
	"github.com/rehacktive/caffeine/service"
)

const (
	projectName = `
	 ██████╗ █████╗ ███████╗███████╗███████╗██╗███╗   ██╗███████╗
	██╔════╝██╔══██╗██╔════╝██╔════╝██╔════╝██║████╗  ██║██╔════╝
	██║     ███████║█████╗  █████╗  █████╗  ██║██╔██╗ ██║█████╗  
	██║     ██╔══██║██╔══╝  ██╔══╝  ██╔══╝  ██║██║╚██╗██║██╔══╝  
	╚██████╗██║  ██║██║     ██║     ███████╗██║██║ ╚████║███████╗
	 ╚═════╝╚═╝  ╚═╝╚═╝     ╚═╝     ╚══════╝╚═╝╚═╝  ╚═══╝╚══════╝	
	`
	projectVersion = "1.0.0"

	MEMORY = "memory"
	FS     = "fs"
	SQLITE = "sqlite"
	PG     = "postgres"
	MYSQL  = "mysql"
	REDIS  = "redis"
	MONGO  = "mongo"

	// env
	envHostPort      = "IP_PORT"
	envDbType        = "DB_TYPE"
	envDbHost        = "DB_HOST"
	envDbName        = "DB_NAME"
	envDbUser        = "DB_USER"
	envDbPass        = "DB_PASS"
	envDbPath        = "DB_PATH"
	envAuthEnabled   = "AUTH_ENABLED"
	envRawSqlEnabled = "RAW_SQL_ENABLED"
)

func main() {
	var addr, dbType, dbHost, dbName, dbUser, dbPass, dbPath string
	var authEnabled bool

	flag.StringVar(&addr, envHostPort, ":8000", "ip:port to expose")
	flag.BoolVar(&authEnabled, envAuthEnabled, false, "enable JWT auth")
	flag.StringVar(&dbType, envDbType, MEMORY, "db type to use, options: memory | fs | sqlite| postgres | mysql | redis | mongo")
	flag.StringVar(&dbPath, envDbPath, "./data", "path of the file storage (for fs or sqlite)")
	flag.StringVar(&dbHost, envDbHost, "localhost", "database host (for postgres | mysql | redis | mongo)")
	flag.StringVar(&dbName, envDbName, "", "database name (for postgres or mysql | mongo)")
	flag.StringVar(&dbUser, envDbUser, "", "database user (for postgres or mysql | mongo)")
	flag.StringVar(&dbPass, envDbPass, "", "database password (for postgres or mysql | mongo)")
	flag.Parse()

	server := service.Server{
		Address:     addr,
		AuthEnabled: authEnabled,
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
	log.Println(projectName)
	log.Println("version: ", projectVersion)
	log.Println("db mode: ", dbType)
	log.Println("server started at: ", server.Address)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	db.Disconnect()

	log.Println("bye")
}
