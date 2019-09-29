# 参数验证
该包实现了对提交参数的简单验证功能，包括整数和浮点型的范围、字符串的长度、字符串正则表达式验证。<br>
该包不提供 graphql 的!表达的 require 验证。但是提供了 require.Requirement 参数来提供类似的功能。需要的时候，只要在Resolver函数中声明参数或承载结构中
定义成员变量即可。
# 验证表达式
<ol>
    <li>整数
        <ol>
            <li>大于 n  => limit=n<$v </li>
            <li>大于等于 n  => limit=n<=$v </li>
            <li>小于 n  => limit=$v<n </li>
            <li>小于等于 n  => limit=$v<=n </li>
            <li>小于m大于n  => limit=n<$v<m </li>
            <li>小于等于m大于等于n  => limit=n<=$v<=m </li>
        </ol>
    </li>
    <li>浮点型
        <ol>
            <li>大于 n  => limit=n<$v </li>
            <li>大于等于 n  => limit=n<=$v </li>
            <li>小于 n  => limit=$v<n </li>
            <li>小于等于 n  => limit=$v<=n </li>
            <li>小于m大于n  => limit=n<$v<m </li>
            <li>小于等于m大于等于n  => limit=n<=$v<=m </li>
        </ol>
    </li>
    <li>字符串长度
        <ol>
            <li>大于 n  => limit=n<$v </li>
            <li>大于等于 n  => limit=n<=$v </li>
            <li>小于 n  => limit=$v<n </li>
            <li>小于等于 n  => limit=$v<=n </li>
            <li>小于m大于n  => limit=n<$v<m </li>
            <li>小于等于m大于等于n  => limit=n<=$v<=m </li>
        </ol>
    </li>
    <li>字符串正则表达式
        <ol>
            <li>regexp=正则表达式<br>
            切记：正则表达式中的转义符 '\' 要写成 '\\'</li>
        </ol>
    </li>
    <li> 错误提示
        <ol>
            <li>error=说明文本<br>
            设置后，如果验证失败将使用 error 的内容作为提示内容，否则自动组织，参见代码</li>
        </ol>
    </li>
</ol>

$v 代表了数值或字符串长度，所有验证大小的都使用 < 或者 <= 设定。

# require.Requirement
必填参数由 require.Requirement 来决定。只要在 Resolver 函数参数中添加 *require.Requirement 即可以
使用。需要注意的是 Requirement.Requires 函数接收的名称是 json Tag，如果没有定义 json Tag 则与成员变量名称一致；
另外，Requirement.Requires 是自动解除最外层参数 in 的。<br>
具体的，还要自己去实验一下。