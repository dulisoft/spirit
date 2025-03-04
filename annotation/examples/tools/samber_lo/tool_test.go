package samber_lo

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"testing"
)

type User struct {
	Name string
	Age  int
}

var users = []*User{
	{
		Name: "Tom",
		Age:  12,
	},
	{
		Name: "Jerry",
		Age:  23,
	},
	{
		Name: "Alice",
		Age:  16,
	},
}

func TestT2(t *testing.T) {
	bytes, err := json.Marshal(users)
	if err != nil {
		return
	}
	fmt.Printf(string(bytes))
}

func TestSliceToMap(t *testing.T) {
	//map[string]*User
	userDict := lo.SliceToMap(users, func(u *User) (string, *User) {
		return u.Name, u
	})
	fmt.Printf("%s", string(lo.T2(json.Marshal(userDict)).A))
}

func TestTimes(t *testing.T) {
	//[]string
	names := lo.Times(len(users), func(index int) string {
		return users[index].Name
	})
	fmt.Printf("%v", names)
}

func TestFilter(t *testing.T) {
	adults := lo.Filter(users, func(item *User, index int) bool {
		return item.Age >= 18
	})
	fmt.Printf("%v", adults)
}
