package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/xdung24/unirest/database"
	"github.com/xdung24/unirest/service"
)

func main() {
	printInfo()

	config := getConfig()

	// create db driver
	var db service.Database
	switch config.DbDriver {
	case MEMORY:
		db = &database.MemDatabase{}
	case FS:
		db = &database.StorageDatabase{
			RootDirPath: config.DbPath,
		}
	case SQLITE:
		db = &database.SQLiteDatabase{
			DirPath: config.DbPath,
		}
	case PG:
		db = &database.PGDatabase{
			Host: config.DbHost,
			Name: config.DbName,
			User: config.DbUser,
			Pass: config.DbPass,
		}
	case MYSQL:
		db = &database.MySqlDatabase{
			Host: config.DbHost,
			Name: config.DbName,
			User: config.DbUser,
			Pass: config.DbPass,
		}
	case REDIS:
		db = &database.RedisDatabase{
			Host: config.DbHost,
		}
	case MONGO:
		db = &database.MongoDatabase{
			Host: config.DbHost,
			Name: config.DbName,
			User: config.DbUser,
			Pass: config.DbPass,
		}
	default:
		panic("invalid db type")
	}

	log.Println("db type: ", config.DbDriver)

	// create web server
	server := service.Server{
		Address:        config.Addr,
		BrokerAddress:  config.BrokerHostPort,
		SwaggerEnabled: config.SwaggerEnabled,
		BrokerEnabled:  config.BrokerEnabled,
		AuthEnabled:    config.AuthEnabled,
		RawSqlEnabled:  config.RawSqlEnabled,
	}
	go server.Init(db)

	log.Println("server started at: ", server.Address)

	// wait for interrupt signal to stop the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	db.Disconnect()

	log.Println("Good bye")
}
