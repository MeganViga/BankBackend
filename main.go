package main

import (
	"database/sql"
	"log"

	db "github.com/MeganViga/BankBackend/db/sqlc"
	"github.com/MeganViga/BankBackend/api"
	_ "github.com/lib/pq"
)
var DBDriver = "postgres"
var DBSource = "postgresql://root:secret@localhost:5432/bankdb?sslmode=disable"
var serverAddress = "0.0.0.0:8080"
func main(){
	conn, err := sql.Open(DBDriver, DBSource)
	if err != nil{
		log.Fatal(err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.StartServer(serverAddress)
	if err != nil{
		log.Fatal("cannot start server", err)
	}
}