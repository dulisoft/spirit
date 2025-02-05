package annotation

import (
	"reflect"
	"testing"
)

type Human struct {
	Class string
	Name  string
}

func TestNew(t *testing.T) {
	p := new(Human)
	t.Log(p == nil)
	vp := reflect.ValueOf(p)
	vp = vp.Elem()
	for i := 0; i < vp.NumField(); i++ {
		fieldValue := vp.Field(i)
		fieldValue.Set(reflect.ValueOf("duyq"))
	}
	t.Logf(p.Class)
	t.Logf(p.Name)
}
