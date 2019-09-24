# 提交参数和返回数据
核心代码,打开 submit.go 查看详细代码
<pre>
// StudentRequest 查询参数
type StudentRequest struct {
    ID int
}

// Student 返回的数据结构
type Student struct {
    ID    int
    Name  string
    Class string
}

// StudentByID 功能函数
func StudentByID(req *StudentRequest) (*Student, error) {
	for _, s := range students {
		if s.ID == req.ID {
			return s, nil
		}
	}
	return nil, errors.New("没有找到")
}
</pre>
上面的代码实现了根据 ID 查询学生的功能，不用那些繁琐的配置信息，只需要关注具体业务即可，怎么样？轻松吧?

## 注意，接收提交的参数必须使用 struct，且必须定义为指针类型，这是硬性规定
## 鉴于 graphql 中不允许出现相同的类型名称，所以参数和返回的数据类型分别不能重复，当有不同包内有相同的结构名称同时作为返回类型时程序将会 panic ，可以通过错误信息找到原因

现在测试一下：
打开 GraphiQl 客户端，在查询编辑框中输入，注意其参数是由 in:{} 包裹的，在第 3 节中将会提供剥去 in 包裹的方法，去与不去就看自己的喜好了 
<pre>
{ 
  StudentByID(in:{ID:2}) {
    ID
    Name
    Class
  }
}
</pre>
点击执行，得到反馈
<pre>
{
  "data": {
    "StudentByID": {
      "Class": "1(1)",
      "ID": 2,
      "Name": "小红"
    }
  }
}
</pre>

除了查询功能，代码中还提供了列表功能，很简单，可自行查看。