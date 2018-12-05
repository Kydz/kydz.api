package Kogger

import (
	"github.com/Kydz/kydz.api/Konfigurator"
	"log"
)

var k *kogger

type kogger struct {
	DebugEnabled  bool
	InfoFilePath  string
	ErrorFilePath string
	DBFilePath    string
	DebugFilePath string
}

func init() {
	k = new(kogger)
	k.DebugEnabled = !Konfigurator.GetKon().IsProd()
}

func Debug(s string, v ...interface{}) {
	if k.DebugEnabled {
		if len(v) > 0 {
			log.Printf(s, v...)
		} else {
			log.Println(s)
		}
	}
}

func Info() {

}

func Error() {

}
