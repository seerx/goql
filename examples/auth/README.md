# 用户认证
由于涉及缓存等内容，在此不在写代码，只说明一下思路<br>
用户认证实就是用注入的功能实现的，类似下面的逻辑
<pre>

type AccountInfo struct {
    // 用户信息定义
}

func authentication(ctx context.Context, r *http.Request, gp *graphql.ResolveParams) *AccountInfo {
    // 1. 从 http.Request 获取请求带的 session
    // 2. 从缓存中找出 session 对应的 AccountInfo
    // 2.1 如果找不到，认证失败，直接 panic(errors.New("需要登录"))
    // 2.2 如果找到，则说明通过认证，直接 return AccountInfo 即可
	return nil
}
// 不要忘记注册 authentication 为注入函数
</pre>

如果某一个 Resolver 需要认证才能操作，可以直接加上 AccountInfo 参数，如下：
<pre>
func getInfo(account *AccountInfo)(*AccountInfo, error) {
    return account, nil
}
</pre>
如果一堆 Resolver 需要认证，那么就用承载结构，把 AccountInfo 定义为承载类的成员变量即可

