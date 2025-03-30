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

func (db *DbConnection) RetrievePasswordHashAndID(usernameOrEmail string) (string, int, error) {
	query := "SELECT user_id, password_hash FROM users WHERE email = ? OR username = ? LIMIT 1"
	stmt, err := db.Db.Prepare(query)
	if err != nil {
		return "", -1, err
	}
	var passwordHash string
	var userID int
	err = stmt.QueryRow(usernameOrEmail, usernameOrEmail).Scan(&userID, &passwordHash)
	if err != nil {
		return "", -1, err
	}
	return passwordHash, userID, nil
}

func (db *DbConnection) StoreAuthToken(userId int, authToken string) error {
	query := "INSERT INTO session_tokens (user_id, session_token) VALUES (?, ?);"
	_, err := db.Db.Exec(query, userId, authToken)
	if err != nil {
		return err
	}
	return nil
}

func (db *DbConnection) RemoveAuthToken(userId int) error {
	query := "DELETE FROM session_tokens WHERE user_id = ?"
	_, err := db.Db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (db *DbConnection) GetAuthTokenFromUserID(userId int) (string, error) {
	query := "SELECT session_token FROM session_tokens WHERE user_id = ? LIMIT 1"
	stmt, err := db.Db.Prepare(query)
	if err != nil {
		return "", err
	}
	var authToken string
	err = stmt.QueryRow(userId).Scan(&authToken)
	if err != nil {
		return "", err
	}
	return authToken, nil
}

func (db *DbConnection) GetUserIDFromAuthToken(authToken string) (int, error) {
	query := "SELECT user_id FROM session_tokens WHERE session_token = ? LIMIT 1"
	stmt, err := db.Db.Prepare(query)
	if err != nil {
		return -1, err
	}
	var userID int
	err = stmt.QueryRow(authToken).Scan(&userID)
	if err != nil {
		return -1, err
	}
	return userID, nil
}

type DbMovie struct {
	ImdbID        string
	Title         string
	ReleaseDate   string // YYYY-MM-DD
	Genre         string
	Actors        string
	Director      string
	Writer        string
	Plot          string
	PosterImgPath string
	MediaType     string // "movie", "series", "episode"
}

const STATIC_PATH = "/images/movies/"

func (movie *DbMovie) ToJson() map[string]any {
	m := make(map[string]any)
	m["imdb_id"] = movie.ImdbID
	m["title"] = movie.Title
	m["release_date"] = movie.ReleaseDate
	m["genre"] = movie.Genre
	m["actors"] = movie.Actors
	m["director"] = movie.Director
	m["writer"] = movie.Writer
	m["plot"] = movie.Plot
	m["poster_img_path"] = movie.PosterImgPath
	m["media_type"] = movie.MediaType
	return m
}

func (db *DbConnection) GetMovieFromID(id string) (DbMovie, error) {
	movie := new(DbMovie)
	query := "SELECT imdb_id, title, release_date, genre, actors, director, writer, plot, poster_image_path, media_type FROM movies WHERE imdb_id = ? LIMIT 1"
	stmt, err := db.Db.Prepare(query)
	if err != nil {
		return *movie, err
	}
	err = stmt.QueryRow(id).Scan(
		&movie.ImdbID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.Genre,
		&movie.Actors,
		&movie.Director,
		&movie.Writer,
		&movie.Plot,
		&movie.PosterImgPath,
		&movie.MediaType,
	)
	if err != nil {
		return *movie, err
	}
	return *movie, nil
}

func (db *DbConnection) GetMoviesFromTitleOrID(titleOrID string, query_limit uint8) ([]DbMovie, error) {
	var movies []DbMovie
	rows, err := db.Db.Query(
		"SELECT imdb_id, title, release_date, genre, actors, director, writer, plot, poster_image_path, media_type FROM movies WHERE title LIKE ? OR imdb_id LIKE ? LIMIT ?",
		"%"+titleOrID+"%",
		"%"+titleOrID+"%",
		string(query_limit),
	)
	if err != nil {
		return movies, err
	}
	for rows.Next() {
		var movie DbMovie
		err = rows.Scan(
			&movie.ImdbID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.Genre,
			&movie.Actors,
			&movie.Director,
			&movie.Writer,
			&movie.Plot,
			&movie.PosterImgPath,
			&movie.MediaType,
		)
		if err != nil {
			txlog.TxLogError("Unable to parse movie DB row to object: %s", err.Error())
			continue
		}
		movies = append(movies, movie)
	}
	return movies, nil
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

func (db *DbConnection) AddMovie(movie DbMovie) error {
	query := `INSERT INTO movies (imdb_id, title, release_date, genre, actors, director, writer, plot, poster_image_path, media_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Db.Exec(
		query,
		&movie.ImdbID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.Genre,
		&movie.Actors,
		&movie.Director,
		&movie.Writer,
		&movie.Plot,
		&movie.PosterImgPath,
		&movie.MediaType,
	)
	if err != nil {
		return err
	}
	return nil
}
