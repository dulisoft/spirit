package main

import (
	"fmt"
	"reflect"
)

// 定义一个结构体，并实现Creatable接口
type MyStruct struct {
	Name string
}

// 定义一个工厂函数映射
var factory = map[string]any{
	"MyStruct": &MyStruct{},
	"Name":     "duyq",
}

// 使用反射来根据类型名称创建对象
func createObject(typeName string) (interface{}, error) {
	if creator, ok := factory[typeName]; ok {
		objValue := reflect.ValueOf(creator)
		objType := reflect.TypeOf(creator)
		for objValue.Kind() == reflect.Pointer {
			objValue = objValue.Elem()
			objType = objType.Elem()
		}
		for i := 0; i < objValue.NumField(); i++ {
			fieldValue := objValue.Field(i)
			key := objType.Field(i).Name
			value, ok := factory[key]
			if !ok {
				return nil, fmt.Errorf("empty type value: %s.%s", typeName, key)
			}
			fieldValue.Set(reflect.ValueOf(value))
		}
		return creator, nil
	}
	return nil, fmt.Errorf("unknown type: %s", typeName)
}

func main() {
	// 根据类型名称创建对象
	obj, err := createObject("MyStruct")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 使用类型断言来检查对象的类型
	if myObj, ok := obj.(*MyStruct); ok {
		fmt.Println("Created object:", myObj)
	} else {
		fmt.Println("Unknown object type")
	}
}
