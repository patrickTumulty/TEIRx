package main

import (
	"fmt"
	"log"
	"net/http"
	"teirxserver/src/cfg"
	"teirxserver/src/txlog"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"slices"
)

type Person struct {
	Name string
	Age  int
}

type User struct {
	FirstName    string
	LastName     string
	DateOfBirth  string // yyyy-MM-dd
	ZipCode      string
	Email        string
	PasswordHash string
	Reputation   int
}

type GinLogForwarder struct{}

func (g GinLogForwarder) Write(p []byte) (n int, err error) {
	txlog.TxLogInfo(string(p))
	return n, nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	gin.DefaultWriter = GinLogForwarder{}
	router := gin.Default()
	router.Use(corsMiddleware())

	router.POST("/login", func(c *gin.Context) {
		var requestBody LoginRequest

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		txlog.TxLogInfo("Received: name=%s pass=%s", requestBody.Username, requestBody.Password)
		fmt.Printf("Received: name=%s pass=%s\n", requestBody.Username, requestBody.Password)

		c.Status(http.StatusOK)
	})

	router.Run("localhost:8080")

	txlog.TxLogInfo("**** Teirx Server Close")
}
