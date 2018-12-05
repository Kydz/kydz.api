package Konfigurator

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Kon struct {
	Site string `json:"site"`
	Db   struct {
		Account string `json:"account"`
		Pass    string `json:"pass"`
		Schema string `json:"schema"`
	} `json:"db"`
	AdminPass string `json:"admin_pass"`
	Env string `json:"env"`
}

var kon *Kon

func init() {
	configJson, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalln("Failed to load Kouter config:" + err.Error())
	}
	c := new(Kon)
	err = json.Unmarshal(configJson, c)
	if err != nil {
		log.Fatalln("Fialed to parse Kouter config:" + err.Error())
	}
	kon = c
}

func GetKon() *Kon {
	return kon
}

func (k *Kon) IsProd() bool {
	return k.Env == "production"
}
