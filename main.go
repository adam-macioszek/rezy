package main

import (
	"database/sql"
	"log"

	"github.com/adam-macioszek/rezy/api"
	"github.com/adam-macioszek/rezy/config"
	db "github.com/adam-macioszek/rezy/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	startApiServer()
}
func startApiServer() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Println("cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannort create server: ", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
