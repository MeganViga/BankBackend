package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	_ "github.com/lib/pq"
)
var testQueries *Queries
var DBDriver = "postgres"
var DBSource = "postgresql://root:secret@localhost:5432/bankdb?sslmode=disable"
func TestMain(m *testing.M){
	conn, err := sql.Open(DBDriver, DBSource)
	if err != nil{
		log.Fatal(err)
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}