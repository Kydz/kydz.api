package Kogger

import "fmt"

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
	k.DebugEnabled = true
}

func Debug(log string, v ...interface{}) {
	if k.DebugEnabled {
		fmt.Printf(log, v...)
	}
}

func Info() {

}

func Error() {

}
