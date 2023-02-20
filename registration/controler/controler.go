package controler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/skowe/masterapp/libs/config"
	"github.com/skowe/masterapp/models"
)

type RegistrationHandle struct {
	DbHost string
	DbPort string
	DbName string
	DbUser string
	DbPass string
}

// Vraća isecak u redosledu user, password, host, port, db
func (reg *RegistrationHandle) Slice() []string {
	return []string{reg.DbUser, reg.DbPass, reg.DbHost, reg.DbPort, reg.DbName}
}

func (reg *RegistrationHandle) Open() (*sql.DB, error) {
	dbConnFmt := "%s:%s@tcp(%s:%s)/%s"
	dbConnString := fmt.Sprintf(dbConnFmt, reg.DbUser, reg.DbPass, reg.DbHost, reg.DbPort, reg.DbName)
	return sql.Open("mysql", dbConnString)
}

func (reg *RegistrationHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	database, err := reg.Open()
	encoder := json.NewEncoder(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encodeData := &models.ApiError{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Error 500: Internal server error",
			ErrorType:    "internal server error",
		}

		encoder.Encode(encodeData)
		return
	}
	defer database.Close()

	switch r.Method {
	case http.MethodPost:
		x := json.NewDecoder(r.Body)
		defer r.Body.Close()

		data := &models.User{}
		err = x.Decode(data)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			encodeData := &models.ApiError{
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "Error 400: Bad request",
				ErrorType:    "bad request",
			}
			encoder.Encode(encodeData)
			return
		}
		selectQuery := "SELECT username, email FROM users WHERE username = ? OR email = ?"

		res, err := database.Query(selectQuery, data.Username, data.Email)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			encodeData := &models.ApiError{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "Error 500: Internal server error",
				ErrorType:    "internal server error",
			}

			encoder.Encode(encodeData)
			return
		}

		if res.Next() {
			apiError := &models.ApiError{
				ErrorCode:    http.StatusSeeOther,
				ErrorMessage: "Error 303: See other",
				ErrorType:    "see other",
			}
			chk := models.User{}
			w.WriteHeader(http.StatusSeeOther) // Prepravi ApiError Model, treba da sadrzi polje extra info
			res.Scan(&(chk.Username), &(chk.Email))
			if chk.Email == data.Email {
				apiError.Info = "The e-mail entered is already in use"
			} else {
				apiError.Info = "The username entered is already in use"
			}

			encoder.Encode(apiError)
			return
		}

		query := "INSERT INTO users(username, email, password) VALUES(?, ?, ?)"

		_, err = database.Exec(query, data.Username, data.Email, data.Password)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encodeData := &models.ApiError{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "Error 500: Internal server error",
				ErrorType:    "internal server error",
			}

			encoder.Encode(encodeData)
			return
		}

		w.WriteHeader(http.StatusOK)

		encoder.Encode(struct {
			Message string `json:"message"`
		}{
			Message: "User created successfully.",
		})
	default:
		apiError := &models.ApiError{
			ErrorCode:    http.StatusMethodNotAllowed,
			ErrorMessage: "Error 405: Method not allowed",
			ErrorType:    "method not allowed",
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		encoder.Encode(apiError)
	}
}

func New() *RegistrationHandle {
	env := config.Configure(config.GetDbEnvNames())
	return &RegistrationHandle{
		DbHost: env[config.DB_HOST],
		DbPort: env[config.DB_PORT],
		DbName: env[config.DB_NAME],
		DbUser: env[config.DB_USER],
		DbPass: env[config.DB_PASS],
	}
}
