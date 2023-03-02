package app

import (
	"log"
	"net/http"

	"github.com/skowe/masterapp/registration/controler"
)

func StartApp() {
	handler := controler.New()
	http.Handle("/register", handler)

	log.Println(http.ListenAndServe(":8080", nil))
}
