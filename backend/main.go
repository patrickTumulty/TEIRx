package main

import (
	"log"
	"teirxserver/src/cfg"
	"teirxserver/src/core"
	"teirxserver/src/dbapi"
	"teirxserver/src/txlog"

	"slices"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type GinLogForwarder struct{}

func (g GinLogForwarder) Write(p []byte) (n int, err error) {
	txlog.TxLogInfo(string(p))
	return n, nil
}

// CORS middleware function definition
func corsMiddleware() gin.HandlerFunc {

	var allowedOrigins = []string{
		"http://localhost:3000",
	}

	// Return the actual middleware handler function
	return func(c *gin.Context) {

		origin := c.Request.Header.Get("Origin")

		if slices.Contains(allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
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
	// gin.DefaultWriter = GinLogForwarder{}
	router := gin.Default()
	router.Use(corsMiddleware())

	core.RegisterRoutes(router)

	router.Run("localhost:8080")

	txlog.TxLogInfo("**** Teirx Server Closing")

	dbapi.GetDBConnection().CloseDbConnection()
}
