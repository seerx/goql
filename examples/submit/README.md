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

现在测试一下：
打开 GraphiQl 客户端，在查询编辑框中输入
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