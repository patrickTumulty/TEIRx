package dbapi

import (
	"database/sql"
	"fmt"
	"sync"
	"teirxserver/src/cfg"
	"teirxserver/src/txlog"
)


var dbConnection *DbConnection
var once sync.Once

type DbConnection struct {
	Db *sql.DB
}

func GetDBConnection() *DbConnection {
    appConfig := cfg.GetAppConfig()

    once.Do(func() {
        dbInfo := appConfig.Database
        dbConnection = newDbConnection(dbInfo.User, dbInfo.Password, dbInfo.Ip, dbInfo.Port, dbInfo.DbName)
    })

    return dbConnection
}

func newDbConnection(user string, password string, ip string, port string, dbname string) *DbConnection {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, ip, port, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
        txlog.TxLogError("Unable to open MySQL connection: %s", err)
        return nil
	}

	dbConn := DbConnection{
		db,
	}

	return &dbConn
}


	// query := "INSERT INTO users (firstname, lastname, date_of_birth, zipcode, email) VALUES (?, ?, ?, ?, ?);"
	// _, err = dbConn.Db.Exec(query, user.FirstName, user.LastName, user.DateOfBirth, user.ZipCode, user.Email)
	// if err != nil {
	// 	panic(err.Error())
	// }
	//

func (db *DbConnection) CloseDbConnection() {
	db.Db.Close()
}


