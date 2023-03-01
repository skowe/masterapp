package controler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/skowe/masterapp/libs/config"
)

const (
	SESSION_HOST_ADDRESS = "SESSION_HOST_ADDRESS"
)

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

func (l *LoginHandle) ServeHTTP(w http.ResponseWriter, r *http.Request){
	
	http.Post(l.SessionHost, "application/json", )
}

func New() *LoginHandle {
	env := config.Configure(append([]string{SESSION_HOST_ADDRESS}, config.GetDbEnvNames()...))
	return &LoginHandle{
		DbHost: env[config.DB_HOST],
		DbPort: env[config.DB_PORT],
		DbName: env[config.DB_NAME],
		DbUser: env[config.DB_USER],
		DbPass: env[config.DB_PASS],
		SessionHost: env[SESSION_HOST_ADDRESS],
	}
}