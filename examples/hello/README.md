# Hello world 示例
<pre>
    g := goql.Get()
    g.RegisterQuery(func() (string, error) {
	return "Hello goql!", nil
    })
    util.StartService(8080)
</pre>

非常简单，这已经完成了一个 Hello world! 示例，运行程序，打开浏览器，在地址栏输入：
<pre>
http://localhost:8080/
</pre>
即可以打开客户端，然后就可以开始开心的测试了
