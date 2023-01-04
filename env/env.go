package env

import (
	"strings"
)

const (
	Local = "local"
	Dev = "dev"
	Test = "test"
	Prod = "prod"
)

var env = Test

func Init(e string){
	if strings.TrimSpace(e) == "" {
		env = Test
	}
	env = e
}

func IsLocal() bool { return env == Local }
func IsTest() bool  { return env == Test }
func IsDev() bool   { return env == Dev }
func IsProd() bool  { return env == Prod }