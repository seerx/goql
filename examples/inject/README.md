# 自动注入及辅助能力
该例实现了参数自动注入功能，可以实现 Resolver 函数的参数或其承载结构中的成员变量自动注入。
只要合理使用注入功能，可以大量减少重复性开发工作，并降低错误出现率
<br>
所谓注入，还是需要自己实现注入函数，然后注册。注入的核心代码如下
<pre>
func init() {
    g := goql.Get()
    g.RegisterInject(InjectClass) // 注册注入函数
}

// 要注入的对象
type ClassInfo struct {
    Grade string `json:"grade"` // 年级
    Class string `json:"class"` // 班级
}

// 注入函数
func InjectClass(ctx context.Context, r *http.Request, gp *graphql.ResolveParams) *ClassInfo {
    return &ClassInfo{
        Grade: "一年级",
        Class: "1 班",
    }
}
</pre>
需要注意的是，注入函数必须是以下形式，返回的类型即为注入的类型，可以是指向结构的指针，也可以是一个 interface
<pre>
func (ctx context.Context, r *http.Request, gp *graphql.ResolveParams) *YourType {
</pre>
同一个类型只能注册一个注入函数。

# 注入到函数的参数中
<pre>
func Inject(class *ClassInfo) (*ClassInfo, error) {
    // 此时 class 的值就等于 InjectClass 返回的值
	return class, nil
}
... ...
goql.Get().RegisterQuery(Inject)
</pre>

# 验证
打开 GraphiQl 客户端，在查询编辑框中输入
<pre>
{
  InjectToLoader {
    grade
    class
  }
}
</pre>
点击执行，得到反馈
<pre>
{
  "data": {
    "InjectToLoader": {
      "class": "1 班",
      "grade": "一年级"
    }
  }
}
</pre>

同样，可以自行测试 Inject 方法

# 辅助功能
实现自动清理资源的功能，例如：自动关闭数据库连接<br>
参见 inject_closer.go 代码<br>
注入方式与普通注入一致，程序运行后，
打开 GraphiQl 客户端，在查询编辑框中输入
<pre>
{
  ReadFromDB
}
</pre>
可以再控制到看到输出
<pre>
准备工作 建立数据库连接
使用数据库连接
清理工作 关闭数据库连接
</pre>

当然，可做的事情，不仅仅是打开和关闭数据库连接，可以根据自己的场景去做一些实用的工作