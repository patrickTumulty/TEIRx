package dbapi

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)


var dbConnection *DbConnection
var once sync.Once

type DbConnection struct {
	Db *sql.DB
}

func GetDBConnection() *DbConnection {
    once.Do(func() {

    })
}

func NewDbConnection(user string, password string, ip string, port string, dbname string) *DbConnection {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, ip, port, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	dbConn := DbConnection{
		db,
	}

	return &dbConn
}

func (db *DbConnection) CloseDbConnection() {
	db.Db.Close()
}


