package main

import (
	"database/sql"
	"log"

	"github.com/Just-A-NoobieDev/bankapi-gin-sqlc/api"
	db "github.com/Just-A-NoobieDev/bankapi-gin-sqlc/db/sqlc"
	"github.com/Just-A-NoobieDev/bankapi-gin-sqlc/util"

	_ "github.com/lib/pq"
)

//	@title			Simple Bank API
//	@version		1.0
//	@description	A simple bank API using Go, gin-gonic framework, postgresql and sqlc

//	@host		localhost:8080
//	@BasePath	/api/v1
func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	
}