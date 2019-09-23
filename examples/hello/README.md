# Hello world 示例
<pre>
func hello() (string, error) {
	return "Hello graphql!", nil
}
func main() {
    g := goql.Get()
    g.RegisterQuery(hello)
    util.StartService(8080)
}
</pre>

非常简单，这已经完成了一个 Hello world! 示例，运行程序，打开浏览器，在地址栏输入：
<pre>
http://localhost:8080/
</pre>
即可以打开 GraphiQl 客户端，在查询编辑框中输入
<pre>
{ hello }
</pre>
点击执行，可以看到反馈结果了
<pre>
{
  "data": {
    "hello": "Hello graphql!"
  }
}
</pre>

这是一个没有参数的例子，虽然很简单，但是 graphql 接口已经可以运行了。
