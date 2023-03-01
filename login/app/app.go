package app

import (
	"net/http"

	"github.com/skowe/masterapp/login/controler"
)

func StartApp() {
	http.Handle("/login", controler.New())
}