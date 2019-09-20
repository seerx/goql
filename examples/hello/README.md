# Hello world 示例
<pre>
    g := goql.Get()
	g.RegisterQuery(func() (string, error) {
		return "Hello goql!", nil
	})
	util.StartService(8080)
</pre>
