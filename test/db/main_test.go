package db

import (
	"database/sql"
	db "github.com/NeptuneYeh/simplerecommend/internal/infra/database/mysql/sqlc"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "mysql"
	dbSource = "user:password@tcp(localhost:3306)/simplerecommend?multiStatements=true&parseTime=true"
)

var testQueries *db.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = db.New(testDB)
	os.Exit(m.Run())
}
