package dbapi

import (
	"database/sql"
	"fmt"
	"sync"
	"teirxserver/src/cfg"
	"teirxserver/src/security"
	"teirxserver/src/txlog"
)

var dbConnection *DbConnection
var once sync.Once

type DbConnection struct {
	Db *sql.DB
}

func GetDBConnection() *DbConnection {
	appConfig := cfg.GetAppConfig()

    // TODO make this blocking
	once.Do(func() {
        txlog.TxLogInfo("Openning database")
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

func (db *DbConnection) CloseDbConnection() {
	db.Db.Close()
}

func RegisterUser(username string, firstname string, lastname string, email string, password string) error {

	hash, err := security.EncodeArgon2Hash(password, security.DefaultArgon2Params())
	if err != nil {
		return err
	}
    
    txlog.TxLogInfo(hash)

	query := "INSERT INTO users (username, firstname, lastname, email, password_hash) VALUES (?, ?, ?, ?, ?);"
    con := GetDBConnection()
    txlog.TxLogInfo("%s", con == nil)
	_, err = GetDBConnection().Db.Exec(query, username, firstname, lastname, email, hash)
	if err != nil {
		return err
	}

	return nil
}

