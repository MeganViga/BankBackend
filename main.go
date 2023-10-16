package main

import (
	"database/sql"
	"log"

	"github.com/MeganViga/BankBackend/api"
	db "github.com/MeganViga/BankBackend/db/sqlc"
	"github.com/MeganViga/BankBackend/util"
	_ "github.com/lib/pq"
)

func main(){
	config, err := util.LoadConfig(".")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil{
		log.Fatal("Cannot load configurations", err)
	}
	if err != nil{
		log.Fatal(err)
	}
	store := db.NewStore(conn)
	server , err := api.NewServer(config, store)
	if err != nil{
		log.Fatal("cannot start server", err)
	}
	err = server.StartServer(config.ServerAddress)
	if err != nil{
		log.Fatal("cannot start server", err)
	}
}