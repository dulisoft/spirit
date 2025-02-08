package example

import (
	"fmt"

	ginx "github.com/gin-gonic/gin"
)

type GreeterService struct {
	Name string
}

type greeterExample struct {
	Name string `json:"greeter_name"`
	ID   string `json:"greeter_id"`
}

var examples []*greeterExample

// PostExample godoc
// @Summary      PostExample Summary
// @Description  PostExample Description
// @Accept       json
// @Produce      json
// @Tags         Greeter1
// @Param        message  body      greeterExample  true  "Greeter Info"
// @Success      200      {string}  string         "success"
// @Failure      500      {string}  string         "fail"
// @Router       /greeter/post [post]
func (s *GreeterService) PostExample(c *ginx.Context, ass *greeterExample) {
	var e = &greeterExample{}

	if err := c.ShouldBind(e); err != nil {
		fmt.Println(e)
		return
	}
	examples = append(examples, e)
	c.JSON(200, "success")
}

//this is c comment

func xx() string {
	return ""
}
