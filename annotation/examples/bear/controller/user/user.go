package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	Name string
}

type userExample struct {
	Name string `json:"greeter_name"`
	ID   string `json:"greeter_id"`
}

var examples []*userExample

// GetExample godoc
// @Summary      PostExample Summary
// @Description  PostExample Description
// @Accept       json
// @Produce      json
// @Tags         Greeter1
// @Param        message  body      userExample  true  "Greeter Info"
// @Success      200      {string}  string         "success"
// @Failure      500      {string}  string         "fail"
// @Router       /greeter/user [post]
func (s *UserService) GetExample(c *gin.Context) {
	var e = &userExample{}

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
