package controler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/skowe/masterapp/libs/config"
	"github.com/skowe/masterapp/models"
)

const (
	SESSION_HOST_ADDRESS = "SESSION_HOST_ADDRESS"
)

var ErrorConnection = errors.New("failed to connect to the database")

type LoginHandle struct {
	DbHost      string
	DbPort      string
	DbName      string
	DbUser      string
	DbPass      string
	SessionHost string
}

func extractUserFromRequest(r *http.Request) (*models.User, error) {

	decoder := json.NewDecoder(r.Body)

	user := &models.User{}
	err := decoder.Decode(user)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer r.Body.Close()

	return user, nil
}

func (reg *LoginHandle) Open() (*sql.DB, error) {
	dbConnFmt := "%s:%s@tcp(%s:%s)/%s"
	dbConnString := fmt.Sprintf(dbConnFmt, reg.DbUser, reg.DbPass, reg.DbHost, reg.DbPort, reg.DbName)
	return sql.Open("mysql", dbConnString)
}

func (l *LoginHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	responseEncoder := json.NewEncoder(w)

	user, err := extractUserFromRequest(r)

	if err != nil {
		responseEncoder.Encode(&models.ApiError{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Error 400: bad request",
			ErrorType:    "bad request",
			Info:         "Malformed request body, can't extract user data",
		})
		return
	}
	user, err = Login(user, l)
	log.Println(err)
	if err != nil {
		if err != ErrorConnection {
			responseEncoder.Encode(&models.ApiError{
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: "Error 401: unauthorized",
				ErrorType:    "unauthorized",
				Info:         "Provided credidentials are invalid",
			})
		} else {
			responseEncoder.Encode(&models.ApiError{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "Error 500: internal server error",
				ErrorType:    "internal server error",
				Info:         "Failed to connect to the users database",
			})
		}
		return
	}
	responseEncoder.Encode(user)

}

func Login(u *models.User, l *LoginHandle) (*models.User, error) {

	selectQuery := "SELECT username, email, password FROM users WHERE username = ? OR email = ?"
	resU := &models.User{}
	db, err := l.Open()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(selectQuery, u.Username, u.Email)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var username, email, pass string
		rows.Scan(&username, &email, &pass)
		if pass != u.Password {
			return nil, errors.New("bad password")
		}

		resU.Username = username
		resU.Email = email
	} else {
		return nil, errors.New("no user")
	}
	return resU, nil
}

func New() *LoginHandle {
	env := config.Configure(append([]string{SESSION_HOST_ADDRESS}, config.GetDbEnvNames()...))
	return &LoginHandle{
		DbHost:      env[config.DB_HOST],
		DbPort:      env[config.DB_PORT],
		DbName:      env[config.DB_NAME],
		DbUser:      env[config.DB_USER],
		DbPass:      env[config.DB_PASS],
		SessionHost: env[SESSION_HOST_ADDRESS],
	}
}
