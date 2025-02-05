package annotation

import (
	"fmt"
	"reflect"
)

// 对象池,key是对象的路径，package/名称
// 此处的初始化，应该给出默认的new(T)对象
var pool = make(map[string]InitObject)

func Init[T any]() {
	uniqueKey := TypePack[T]()
	pool[uniqueKey] = InitObject{
		Has: false,
		Obj: new(T),
	}
}

func Put[T any](t *T) {
	uniqueKey := TypePack[T]()
	pool[uniqueKey] = InitObject{
		Has: true,
		Obj: t,
	}
}

type InitObject struct {
	Obj any
	Has bool
}

// 使用反射来根据类型名称创建对象
func createObject[T any]() (*T, error) {
	uniqueKey := TypePack[T]()
	if _, ok := pool[uniqueKey]; !ok {
		return nil, fmt.Errorf("unknown type: %s", uniqueKey)
	}
	obj := new(T)
	objValue := reflect.ValueOf(obj)
	objType := reflect.TypeOf(obj)
	for objValue.Kind() == reflect.Pointer {
		objValue = objValue.Elem()
		objType = objType.Elem()
	}
	for i := 0; i < objValue.NumField(); i++ {
		fieldValue := objValue.Field(i)
		key := objType.Field(i).Name
		value := getInstanceVlaueOf(key)

		fieldValue.Set(reflect.ValueOf(value))
	}
	return obj, nil
}

func getInstanceVlaueOf(uniqueKey string) reflect.Value {
	t := getInstanceByUniqueKey(uniqueKey)
	objValue := reflect.ValueOf(t)
	for objValue.Kind() == reflect.Pointer {
		objValue = objValue.Elem()
	}
	return objValue
}

// getInstanceByUniqueKey    获取某个struct的实例
func getInstanceByUniqueKey(uniqueKey string) any {
	initObj, has := pool[uniqueKey]
	if !has || !initObj.Has {
		panic(fmt.Sprintf("invalid type %v", uniqueKey))
	}
	return initObj.Obj
}

// GetInstance 获取某个struct的实例
func GetInstance[T any]() *T {
	uniqueKey := TypePack[T]()
	initObj, has := pool[uniqueKey]
	if !has {
		panic(fmt.Sprintf("invalid type %v", uniqueKey))
	}
	if initObj.Has {
		t, ok := initObj.Obj.(*T)
		if !ok {
			panic(fmt.Sprintf("invalid type %v", uniqueKey))
		}
		if t == new(T) {
			panic(fmt.Sprintf("empty type %v", uniqueKey))
		}
		return t
	}
	//自动生成逻辑
	obj, err := createObject[T]()
	if err != nil {
		panic(fmt.Sprintf("create type %v,error %v", uniqueKey, err))
	}
	return obj
}
