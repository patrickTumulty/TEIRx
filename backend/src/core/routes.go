package core

import (
	"database/sql"
	"math"
	"math/rand/v2"
	"net/http"
	"pen_daemon/src/dbapi"
	"pen_daemon/src/omdb"
	"pen_daemon/src/security"
	"pen_daemon/src/txlog"
	"strings"

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

func handleGetFilm(c *gin.Context) {

	txlog.TxLogInfo("Getting film")

	id := strings.TrimSpace(c.DefaultQuery("id", ""))
	if id == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	film, err := omdb.OmdbGetByID(id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	ranks, err := dbapi.GetDBConnection().GetFilmRanks(id)
	stats := make(map[rune]int)

	for _, item := range ranks {
		tier := item.GetTierAsRune()
		value, ok := stats[tier]
		if !ok {
			stats[tier] = 1
		} else {
			stats[tier] = value + 1
		}
	}

	avg := 0
	totalRanks := 0
	for key, value := range stats {
		totalRanks += value
		mult := TierStr2Int(key)
		avg += value * mult
	}

	tier := int(math.Ceil(float64(avg) / float64(totalRanks)))

	filmStats := make(map[string]any)
	if totalRanks == 0 {
		filmStats["rank"] = "Not Yet Ranked"
	} else {
		filmStats["rank"] = string(TierInt2Str(tier))
	}

	filmStats["total_count"] = totalRanks
	filmStats["s_count"] = stats['S']
	filmStats["a_count"] = stats['A']
	filmStats["b_count"] = stats['B']
	filmStats["c_count"] = stats['C']
	filmStats["d_count"] = stats['D']
	filmStats["f_count"] = stats['F']

	filmData := make(map[string]any)
	filmData["title"] = film.Title
	filmData["plot"] = film.Plot
	filmData["rated"] = film.Rated
	filmData["year"] = film.Year
	filmData["poster"] = film.Poster
	filmData["stats"] = filmStats

	c.JSON(http.StatusOK, filmData)
}

func handleSearch(c *gin.Context) {

	txlog.TxLogInfo("Searching...")

	query := strings.TrimSpace(c.DefaultQuery("query", ""))
	if query == "" {
		c.Status(http.StatusBadRequest)
		return
	}

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

func handleGetFeatured(c *gin.Context) {

	featured := []string{
		"Mickey 17",
		"Blade Runner",
		"Interstellar",
		"Alien",
		"Dune",
		"Casino Royale",
	}

	randomInt := rand.IntN(len(featured))

	film := featured[randomInt]

	items, err := omdb.OmdbSearch(film)
	if err != nil {
		txlog.TxLogError("Error searching OMDb for featured: %s", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(items) == 0 {
		txlog.TxLogError("No films found for query: '%s'", film)
		c.Status(http.StatusInternalServerError)
		return
	}

	item := items[0]

	filmData := make(map[string]any)
	filmData["title"] = item.Title
	filmData["year"] = item.Year
	filmData["poster"] = item.Poster

	c.JSON(http.StatusOK, filmData)
}

func RegisterRoutes(router *gin.Engine) {
	router.POST("/login", handleLogin)
	router.POST("/logout", handleLogout)
	router.POST("/register-user", handleRegisterUser)
	router.GET("/search", handleSearch)
	router.GET("/get-film", handleGetFilm)
	router.GET("/featured", handleGetFeatured)
}
