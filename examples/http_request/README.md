# 得到 http.Request
没错，就是通过注入功能得到当前请求的 http.Request 参数<br>
注入函数不能更简单，如下

<pre>
func injectRequest(ctx context.Context, r *http.Request, gp *graphql.ResolveParams) *http.Request {
	return r
}
// 注册
... ...
goql.Get().RegisterInject(injectRequest)
</pre>

不管在哪里使用，只要在 Resolver 函数中声明 *http.Request 参数，即可以拿到当前 http 请求的信息
<br>
可以自行验证执行效果。