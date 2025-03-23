package main

import (
	"log"
	"net/http"
	"teirxserver/src/cfg"
	"teirxserver/src/txlog"

	"github.com/gin-gonic/gin"
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

type GinLogForwarder struct{}

func (g GinLogForwarder) Write(p []byte) (n int, err error) {
	txlog.TxLogInfo(string(p))
	return n, nil
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

	gin.DisableConsoleColor()
	gin.DefaultWriter = GinLogForwarder{}
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello")
	})

	router.Run("localhost:8080")

	txlog.TxLogInfo("**** Teirx Server Close")
}
