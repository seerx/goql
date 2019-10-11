# API 分组
当接口比较多时，我们希望对接口进行分类，这里提供了 api 名称加前缀的方法来实现类似能力。
<br>
只要在 Resolver 函数的承载结构中定义 Prefix 成员变量（类型不限），然后定义响应的 tag 。
那么该结构所有被解析为 Resolver 的函数对应的 api 名称都会自动加上前缀。如下：
<pre>
// StudentLoader 学生信息操作承载类
type StudentLoader struct {
	// 定义 api 前缀为 student
	Prefix string `gql:"prefix=student"`
}

// Hello hello
func (StudentLoader) Hello() (string, error) {
	return "Hello Student", nil
}

...
goql.Get().RegisterQuery(StudentLoader{})
...

</pre>
Hello 函数对应的 api，将会更名为 studentHello，而不是 Hello。
