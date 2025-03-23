package main

import (
	"log"
	"teirxserver/src/cfg"
	"teirxserver/src/txlog"

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

func main() {

    err := cfg.LoadAppConfig("teirxcfg.json")
    if err != nil {
        log.Panic(err)
        return
    }

    appConfig := cfg.GetAppConfig()

    logLevel := txlog.Str2LogLevel(appConfig.Logging.Level)
    err = txlog.InitLogging(logLevel, appConfig.Logging.Filepath)
    if err != nil {
        log.Panic(err)
        return
    }

    txlog.TxLogInfo("**** Teirx Server Starting")

    txlog.TxLogInfo("Hello, World %d!!!", 5);
    txlog.TxLogWarn("Hello, World %d!!!", 5);
    txlog.TxLogError("Hello, World %d!!!", 5);
    txlog.TxLogDebug("Hello, World %d!!!", 5);

	// appConfig, err := cfg.LoadAppConfig("./teirxcfg.json")
	// if err != nil {
	// 	return
	// }
	//
	// user := User{"", "", "", "", ""}
	//
	// fmt.Print("Enter First Name: ")
	// fmt.Scanln(&user.FirstName)
	// fmt.Print("Enter Last Name: ")
	// fmt.Scanln(&user.LastName)
	// fmt.Print("Enter Date Of Birth (yyyy-mm-dd): ")
	// fmt.Scanln(&user.DateOfBirth)
	// fmt.Print("Enter Zip: ")
	// fmt.Scanln(&user.ZipCode)
	// fmt.Print("Enter Email: ")
	// fmt.Scanln(&user.Email)
	//
	// dbConn := NewDbConnection(&appConfig.Database)
	//
	// query := "INSERT INTO users (firstname, lastname, date_of_birth, zipcode, email) VALUES (?, ?, ?, ?, ?);"
	// _, err = dbConn.Db.Exec(query, user.FirstName, user.LastName, user.DateOfBirth, user.ZipCode, user.Email)
	// if err != nil {
	// 	panic(err.Error())
	// }
	//
	// dbConn.CloseDbConnection()

    txlog.TxLogInfo("**** Teirx Server Close")
}
