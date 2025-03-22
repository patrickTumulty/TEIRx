package main

import (
	"database/sql"
	"fmt"
	"log"
	"teirxserver/src/cfg"

	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	Name string
	Age  int
}

type User struct {
	FirstName   string
	LastName    string
	DateOfBirth string // yyyy-MM-dd
	ZipCode     string
	Email       string
}

func (usr *User) ToSQLString() string {
	return fmt.Sprintf("'%s', '%s', '%s', '%s', '%s'", usr.FirstName, usr.LastName, usr.DateOfBirth, usr.ZipCode, usr.Email)
}

type DbConnection struct {
	Db *sql.DB
}

func NewDbConnection(dbInfo *cfg.DbInfo) *DbConnection {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbInfo.User, dbInfo.Password, dbInfo.Ip, dbInfo.Port, dbInfo.DbName)

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

func main() {

	appConfig, err := cfg.LoadAppConfig("./teirxcfg.json")
	if err != nil {
		return
	}

	user := User{"", "", "", "", ""}

	fmt.Print("Enter First Name: ")
	fmt.Scanln(&user.FirstName)
	fmt.Print("Enter Last Name: ")
	fmt.Scanln(&user.LastName)
	fmt.Print("Enter Date Of Birth (yyyy-mm-dd): ")
	fmt.Scanln(&user.DateOfBirth)
	fmt.Print("Enter Zip: ")
	fmt.Scanln(&user.ZipCode)
	fmt.Print("Enter Email: ")
	fmt.Scanln(&user.Email)

	dbConn := NewDbConnection(&appConfig.Database)

    query := "INSERT INTO users (firstname, lastname, date_of_birth, zipcode, email) VALUES (?, ?, ?, ?, ?);" 
	_, err = dbConn.Db.Exec(query, user.FirstName, user.LastName, user.DateOfBirth, user.ZipCode, user.Email)
	if err != nil {
		panic(err.Error())
	}

	dbConn.CloseDbConnection()
}
