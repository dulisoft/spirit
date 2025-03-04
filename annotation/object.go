package annotation

// Object 被注解对象
type Func struct {
	Name     string    //名称
	pkg      *Pkg      //包信息
	FuncSign *FuncSign //方法签名
}

type Pkg struct {
	Name string //包名称, 不包含名称
	Path string //包路径, 全路径
	File string //文件路径,相对于项目的相对路径
}

type FuncSign struct {
	Reveiver IdentType    //方法的接受者
	Args     []*IdentType //参数类型数组
	Returns  []*IdentType //返回值
}
