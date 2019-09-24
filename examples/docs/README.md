# API 及参数说明
这个例子是 submit.go 的复刻版，功能是一样的。

# 不同之处
<ol>
<li>API 和参数都增加了说明信息，这才更graphql</li>
<li>Resolver 函数不再是独立的 func，而是属于一个结构</li>
<li>Resolver 函数参数可以直接使用 id 了，不需要写成 in:{id: 1}</li>
</ol>

打开 GraphiQl 客户端，在查询编辑框中输入
<pre>
{ 
  Query(id:1) {
    id
    name
    class
  }
}
</pre>
点击执行，得到反馈
<pre>
{
  "data": {
    "Query": {
      "class": "1(1)",
      "id": 1,
      "name": "小明"
    }
  }
}
</pre>

除了查询功能，代码中还提供了列表功能，很简单，可自行查看。