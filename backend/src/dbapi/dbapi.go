package dbapi

import (
	"database/sql"
	"fmt"
	"sync"
	"teirxserver/src/cfg"
	"teirxserver/src/security"
	"teirxserver/src/txlog"
)

var dbInitMutex sync.Mutex = sync.Mutex{}
var dbConnection *DbConnection

type DbConnection struct {
	Db *sql.DB
}

func GetDBConnection() *DbConnection {

	appConfig := cfg.GetAppConfig()

	dbInitMutex.Lock()
	defer dbInitMutex.Unlock()

	if dbConnection == nil {
		txlog.TxLogInfo("Openning database")
		dbInfo := appConfig.Database
		dbConnection = newDbConnection(dbInfo.User, dbInfo.Password, dbInfo.Ip, dbInfo.Port, dbInfo.DbName)
	}

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

func (db *DbConnection) Close() {
	db.Db.Close()
}

func (db *DbConnection) RetrievePasswordHash(usernameOrEmail string) (string, error) {
	query := "SELECT password_hash FROM users WHERE email = ? OR username = ? LIMIT 1"
	stmt, err := db.Db.Prepare(query)
	if err != nil {
        return "", err
	}
	var passwordHash string
	err = stmt.QueryRow(usernameOrEmail, usernameOrEmail).Scan(&passwordHash)
	if err != nil {
		return "", err
	}
    return passwordHash, nil
}

func (db *DbConnection) RegisterUser(username string, firstname string, lastname string, email string, password string) error {

	hash, err := security.EncodeArgon2Hash(password, security.DefaultArgon2Params())
	if err != nil {
		return err
	}

	query := "INSERT INTO users (username, firstname, lastname, email, password_hash) VALUES (?, ?, ?, ?, ?);"
	_, err = db.Db.Exec(query, username, firstname, lastname, email, hash)
	if err != nil {
		return err
	}

	return nil
}
