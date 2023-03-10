package app

import (
	"log"
	"net/http"

	"github.com/skowe/masterapp/login/controler"
)

func StartApp() {
	http.Handle("/login", controler.New())

	log.Println(http.ListenAndServe(":8081", nil))
}
