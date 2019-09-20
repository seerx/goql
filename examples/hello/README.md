# Hello world 示例
<pre>
func main() {
    g := goql.Get()
    g.RegisterQuery(func() (string, error) {
	return "Hello goql!", nil
    })
    util.StartService(8080)
}
</pre>

非常简单，这已经完成了一个 Hello world! 示例，运行程序，打开浏览器，在地址栏输入：
<pre>
http://localhost:8080/
</pre>
即可以打开客户端，然后就可以开始开心的测试了

虽然很简单，但是已经可以开始了。
