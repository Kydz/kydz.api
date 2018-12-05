package handlers

import (
	"encoding/json"
	"github.com/Kydz/kydz.api/models"
	"io/ioutil"
	"log"
	"net/http"
)

func GetInit(w http.ResponseWriter, r *http.Request) {
	err := models.InitAdmin()
	if err != nil {
		log.Println("[Error] Init admin failed: " + err.Error())
	} else {
		normalResponse(`{"success": true}`, w)
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
	}
	var a = new(models.Admin)
	err = json.Unmarshal(d, &a)
	if err != nil {
	}
	token, err := models.Login(a.Password)
	if err != nil {
		log.Println("[Error] Login failed: " + err.Error())
		normalResponse(`{"success": false}`, w)
	} else {
		normalResponse(`{"token": "` + token + `"}`, w)
	}
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := r.Header.Get("x-auth-token")
		if models.CheckAdminToken(t) {
			next(w, r)
		} else {
			http.Error(w, "access deny", http.StatusUnauthorized)
		}
	}
}
