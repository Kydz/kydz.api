package models

import (
	"crypto/sha256"
	"fmt"
	"github.com/Kydz/kydz.api/utils"
	"time"
)

type Admin struct {
	Email string `json:"email"`
	Salt string `json:"salt"`
	Password string `json:"password"`
	Name string `json:"name"`
	Token string `json:"token"`
	RememberToken string `json:"remember_token"`
}

func initAdmin(salt string, rp string) error {
	h := sha256.New()
	h.Write([]byte(salt + rp))
	p := fmt.Sprintf("%x",  h.Sum(nil))
	qs := `INSERT INTO admins (name, email, salt, password, token, remember_token, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	logQuery(qs + `psmask: %s, ps: %s`, salt, p)
	_, err := db.Exec(qs, "Kydz", "boatkyd@gmail.com", salt, p, utils.GenerateRandomString(32), utils.GenerateRandomString(64), time.Now(), time.Now())
	return err
}

func queryAdmin() (a Admin) {
	qs := `SELECT email, name, password, salt, token, remember_token from admins WHERE email = "boatkyd@gmail.com";`
	logQuery(qs)
	row := db.QueryRow(qs)
	err := row.Scan(&a.Email, &a.Name, &a.Password, &a.Salt, &a.Token, &a.RememberToken)
	if err != nil {
	}
	return a
}
