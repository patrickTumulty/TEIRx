package core

import (
	"database/sql"
	"net/http"
	"strings"
	"teirxserver/src/dbapi"
	"teirxserver/src/omdb"
	"teirxserver/src/security"
	"teirxserver/src/txlog"

	"github.com/gin-gonic/gin"
)

type LogoutRequest struct {
	AuthToken string `json:"token"`
}

func IsError(err error) bool {
	return err != nil && err != sql.ErrNoRows
}

func handleLogout(c *gin.Context) {
	var requestBody LogoutRequest

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		txlog.TxLogError("Unable to parse JSON: %s", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	userId, err := dbapi.GetDBConnection().GetUserIDFromAuthToken(requestBody.AuthToken)
	if err != nil {
		txlog.TxLogError("Unable to retreive user ID with auth token: %s", err.Error())
		c.Status(http.StatusUnauthorized)
		return
	}

	err = dbapi.GetDBConnection().RemoveAuthToken(userId)
	if err != nil {
		txlog.TxLogError("Unable remove auth token: %s", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

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

	dbHash, userId, err := dbapi.GetDBConnection().RetrievePasswordHashAndID(requestBody.Username)
	if err != nil {
		txlog.TxLogError("Unable to retrieve user info: %s", err.Error())
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNoContent, gin.H{"error": "no user found"})
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
		return
	}

	token, err := dbapi.GetDBConnection().GetAuthTokenFromUserID(userId)
	if err != nil {
		token, err = security.GenerateAuthToken()
		if err != nil {
			txlog.TxLogError("Unable to generate auth token: %s", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		err = dbapi.GetDBConnection().StoreAuthToken(userId, token)
		if err != nil {
			txlog.TxLogError("Unable to store auth token: %s", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
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

func handleSearch(c *gin.Context) {
    query := strings.TrimSpace(c.DefaultQuery("query", ""))
    if query == "" {
        c.Status(http.StatusBadRequest)
        return 
    }

	// if strings.HasPrefix(query, "tt") { // Maybe IMDb ID
	// 	movie, err := dbapi.GetDBConnection().GetMovieFromID(query)
	// 	if err != nil {
	// 		if err != sql.ErrNoRows {
	// 			txlog.TxLogError("Error finding movie from ID %s", query)
	// 		}
	// 	} else {
	// 		c.JSON(http.StatusOK, gin.H{"results": []gin.H{
	// 			movie.ToJSON(),
	// 		}})
	// 		return
	// 	}
	// }
	//
	// movies, err := dbapi.GetDBConnection().GetMoviesFromTitleOrID(query, 15)
	// if err != nil {
	// 	if err != sql.ErrNoRows {
	// 		txlog.TxLogError("Error finding movies from ID or Title '%s'", query)
	// 	}
	// } else {
 //        json_movies := []gin.H{}
 //        for _, m := range movies {
 //            json_movies = append(json_movies, m.ToJSON())
 //        }
	// 	c.JSON(http.StatusOK, gin.H{"results": json_movies})
	// 	return
	// }

    // No we need to query OMDb to fill in gaps in our database

    items, err := omdb.OmdbSearch(query)
    if err != nil {
        txlog.TxLogError("Error searching OMDb: %s", err.Error())
        c.Status(http.StatusInternalServerError)
        return
    }

    jsonItems := []gin.H{}
    for _, item := range items {
        jsonItems = append(jsonItems, item.ToJson())
    } 

	c.JSON(http.StatusOK, jsonItems)
}

func RegisterRoutes(router *gin.Engine) {
	router.POST("/login", handleLogin)
	router.POST("/logout", handleLogout)
	router.POST("/register-user", handleRegisterUser)
	router.GET("/search", handleSearch)
}
