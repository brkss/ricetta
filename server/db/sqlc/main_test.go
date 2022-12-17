package db

import (
	"os"
	"testing"
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:root@localhost:5432/vanillefraise?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M){
	
	con, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	testQueries = New(con)

	os.Exit(m.Run())


}
