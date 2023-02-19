package app

import (
	"log"
	"net/http"

	"github.com/skowe/masterapp/registration/controler"
)

func StartApp() {
	http.Handle("/register", controler.New())

	log.Println(http.ListenAndServe(":8080", nil))
}
