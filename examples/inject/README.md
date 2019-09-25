# 自动注入及辅助能力
该例实现了参数自动注入功能，可以实现 Resolver 函数的参数或其承载结构中的参数自动注入。
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
需要注意的是，注入函数必须是一下形式，返回的类型即为注入的类型，可以是指向结构的指针，也可以是一个 interface
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

# 注入到结构中的字段中
<pre>
type Loader struct {
	Class *ClassInfo
}

func (l Loader) InjectToLoader() (*ClassInfo, error) {
    // 执行到此处时 l.Class 的值由 InjectClass 函数提供
	return l.Class, nil
}

... ...
goql.Get().RegisterQuery(Loader{})
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