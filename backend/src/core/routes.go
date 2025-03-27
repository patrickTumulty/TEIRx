package core

import (
	"database/sql"
	"net/http"
	"teirxserver/src/dbapi"
	"teirxserver/src/security"
	"teirxserver/src/txlog"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handleLogin(c *gin.Context) {
	var requestBody LoginRequest

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		txlog.TxLogError("Unable to parse JSON: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbHash, err := dbapi.GetDBConnection().RetrievePasswordHash(requestBody.Username)
	if err != nil {
		txlog.TxLogError("Unable to retrieve user info: %s", err.Error())
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

    err = security.AuthenticatePassword(requestBody.Password, dbHash)
    if err != nil {
        txlog.TxLogError("Failed to authenticate password: %s", err.Error())
        if err == security.ErrPasswordsDoNotMatch {
            c.Status(http.StatusUnauthorized)
        } else {
            c.Status(http.StatusInternalServerError)
        }
    }

	c.JSON(http.StatusOK, gin.H{"token": "0xCAFE"})
}

type RegisterUserRequest struct {
	Username     string `json:"username"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}

func handleRegisterUser(c *gin.Context) {
	var requestBody RegisterUserRequest

	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		txlog.TxLogError("Register User Error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = dbapi.GetDBConnection().RegisterUser(
		requestBody.Username,
		requestBody.Firstname,
		requestBody.Lastname,
		requestBody.Email,
		requestBody.PasswordHash,
	)

	if err != nil {
		txlog.TxLogError("Unable to add user to database: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txlog.TxLogInfo("Status OK")

	c.Status(http.StatusOK)
}

func RegisterRoutes(router *gin.Engine) {
	router.POST("/login", handleLogin)
	router.POST("/register-user", handleRegisterUser)
}
