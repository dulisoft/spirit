package main

import (
	"bear/controller/example"
	"bear/controller/user"

	"github.com/dulisoft/spirit/annotation"
)

func InitObjectPool() {

	annotation.Put[example.GreeterService](PostExample)
	annotation.Put[user.UserService](GetExample)
}
