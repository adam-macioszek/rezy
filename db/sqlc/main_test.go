package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/adam-macioszek/rezy/config"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var store Store

func TestMain(m *testing.M) {
	//consider using os.Getenv("DATABASE_URL")
	config, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Println(config.DBSource)
		log.Fatal("cannot open database connection")
	}
	testQueries = New(conn)
	store = NewStore(conn)
	defer func() {
		_ = conn.Close()
		fmt.Println("closing database connection")
	}()
	os.Exit(m.Run())
}
