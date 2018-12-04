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
}

var kon *Kon

func init() {
	configJson, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("Failed to load Kouter config: %+v", err)
	}
	c := new(Kon)
	err = json.Unmarshal(configJson, c)
	if err != nil {
		log.Fatalf("Fialed to parse Kouter config: %+v", err)
	}
	kon = c
}

func GetKon() *Kon {
	return kon
}
