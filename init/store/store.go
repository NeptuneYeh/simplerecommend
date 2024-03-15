package store

import (
	"database/sql"
	"github.com/NeptuneYeh/simplerecommend/init/config"
	db "github.com/NeptuneYeh/simplerecommend/internal/infra/database/mysql/sqlc"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var MyStore *db.Store

type Module struct {
	Store db.Store
}

func NewModule() *Module {
	// init Store
	conn, err := sql.Open(config.MyConfig.DBDriver, config.MyConfig.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	store := db.NewStore(conn)
	MyStore = &store

	storeModule := &Module{
		Store: store,
	}

	return storeModule
}
