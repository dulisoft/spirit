package main

import (
	"bear/controller/example"
	"bear/controller/user"

	"github.com/dulisoft/spirit/annotation"
)

func InitObjectPool() {
	annotation.Put(new(example.GreeterService))
	annotation.Put(new(user.UserService))
}
