package controler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

func (reg *LoginHandle) Slice() []string {
	return []string{reg.DbUser, reg.DbPass, reg.DbHost, reg.DbPort, reg.DbName}
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
	err = Login(user, l)

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
	user.Password = ""
	responseEncoder.Encode(user)

}

func Login(u *models.User, l *LoginHandle) error {

	selectQuery := "SELECT password FROM users WHERE username = ? OR email = ?"

	db, err := l.Open()
	if err != nil {
		return err
	}

	rows, err := db.Query(selectQuery, u.Username, u.Password)
	if err != nil {
		return err
	}
	if rows.Next() {
		var pass string
		rows.Scan(&pass)
		if pass != u.Password {
			return errors.New("credidential missmatch")
		}
	} else {
		return errors.New("credidential missmatch")
	}
	return nil
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

func extractUserFromRequest(r *http.Request) (*models.User, error) {

	decoder := json.NewDecoder(r.Body)

	user := &models.User{}
	err := decoder.Decode(user)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return user, nil
}

func GetSessionParams(user *models.User, l *LoginHandle) (*http.Response, error) {
	userEncoded, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		return nil, err
	}
	// Ovde je neophodno napraviti pozeljni reader objekat za post metod ka session menadzeru
	RequestBody := bytes.NewReader(userEncoded)
	return http.Post(l.SessionHost+"/start", "application/json", RequestBody)
}
