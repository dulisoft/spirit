package main

import (
	"bear/controller/example"
	"bear/controller/user"

	"github.com/dulisoft/spirit/annotation"
	"github.com/gin-gonic/gin"
)

func main() {
	InitObjectPool()
	r := gin.Default()
	r.Handle("POST", "/greeter/post", annotation.GetInstance[example.GreeterService]().PostExample)
	r.Handle("POST", "/greeter/user", annotation.GetInstance[user.UserService]().GetExample)
	if err := r.Run("0.0.0.0:8888"); err != nil {
		panic(err)
	}
}
