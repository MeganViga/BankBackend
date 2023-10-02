package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/MeganViga/BankBackend/util"
	_ "github.com/lib/pq"
)
var testQueries *Queries
var testStore *Store
// var DBDriver = "postgres"
// var DBSource = "postgresql://root:secret@localhost:5432/bankdb?sslmode=disable"
func TestMain(m *testing.M){
	config, err := util.LoadConfig("../..")
	if err != nil{
		log.Fatal("Cannot load config", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil{
		log.Fatal(err)
	}
	testQueries = New(conn)
	testStore = NewStore(conn)
	os.Exit(m.Run())
}