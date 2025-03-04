# 基础功能
读取注解功能
   1：得到注解的名称
   2：得到注解的方法信息
   3：得到注解后面的参数

对象池：
   1：使用注解
   2：使用反射构造对象
   3：初始化到对象池

自动代码
   1：构造方法
   2：模版方法
   3：自定义 

项目文档：
   1：swag文档
   2：swag文档注解自定义
   

然后就是处理注解，对注解的处理可以绑定很多方法，这个是一连串的方法，挨个执行即可



// Object 被注解对象
type Object struct {
	Name     string    //名称
	PkgPath  string    //包路径
	PkgName  string    //包名称
	FuncSign *FuncSign //方法签名
}

type FuncSign struct {
	Reveiver IdentType    //方法的接受者
	Args     []*IdentType //参数类型数组
	Returns  []*IdentType //返回值
}

type StructInfo struct {
	Fields []*FieldInfo  
}

type FieldInfo struct {
	Name string
	Type string
}

type IdentType struct {
	Path string
	Name string
}

