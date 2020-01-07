# ZeroDB代码规范建议稿

## 项目的基本说明

- 本项目统一使用官方包管理工具`dep`作为第三方库管理工具。

- 各子项目的编译输出统一为`.out`后缀名，方便做`git`排除。

## 代码源文件命名

字母全部小写，如果是多个单词，用_隔开，比如「msg_receiver.go」

## 代码格式

`Go`本身有官方代码格式，在提交代码前务必使用`gofmt`工具格式代码。

## 包名命名

在`Go`中包是以目录形式组织，在包名的命名上尽量是单个单词，且全部字母小写

## 局部变量命名

1. 变量命名要纯英文，不可拼音，要词达意

2. 结构体中的成员变量

    - 在`Go`中的规范是首字母大写则为`public`成员变量，首字母小写则为`private`成员变量。

3. 变量命名的原则是驼峰命名，比如`var logLevel int`

## 函数命名

1. 函数命名不允许使用close、append、new、copy、delete、len、cap、make、complex、real、imag、panic、print、recover、println …。以免跟系统函数冲突,产生潜在的错误。

## `struct`命名

1. 首字母大写则为可导出，首字母小写则为不可导出。

2. 命名的原则是驼峰命名。

3. `Go`中有一个语法糖，一个`struct`中包含另外一个`struct`则默认导入被包含的`struct`的全部成员变量和成员函数。

```go
    type person struct {
        name string
    }

    func (p *person) SayHi() {
        fmt.Println("hello")
    }

    type student struct {
        person
        id int
    }
```

4. 在`Go`中故意模糊了指针跟结构体值的区别。在`struct`的成员函数的时候尽量用指针。

5. 成员函数接收器不要使用诸如`this`,`self`等面向对象的参数。

```go
type My struct {}

func (m *My) Say() {}
```

## `goroutine`的使用

1.`channel`原则上要成对使用（要有写有读），否则有可能造成死锁。使用完后要注意`close`。

2.需要等待`goroutine`退出的，则尽量用`sync.WaitGroup`来阻塞主`goroutine`。

## 单元测试

`Go`的单元测试原则上是以文件为单元的。命名上以文件名加`_test`来命名。比如：`session.go`，`session_test.go`。

单元测试模板：

```go
package session

import "testing"

func TestFunctionName(t *testing.T) {
}
```

## 返回值问题

* 在`Go`中可以在一个函数中返回一个局部变量地址，但在实际编程中尽量不要这样。

```go
func ReturnString() *string {
    str := "hello"
    return &str;
}
```

虽然这种在`Go`是正确的，但我们不要这么返回。

* 定义函数或者方法返回值时尽量在函数体中创建返回值,而不是在定义的时候创建.

  - Example A:

    ```go
    func（a *A）exampleA（name string)(*task, Error){
        ...
        t := &task{...}
        e := ...
        return t, e
    }
    ```

  - Example B:

    ```go
    func（a *B）exampleB（name string)(t *task, e Error){
        ...

        return
    }
    ```

推荐使用**Example A**的方式,尽量不要使用`Example B`的方式,使用B的方式如果忘记给`t`、 `e`赋值可能会导致返回`nil`引起运行时的空指针`panic`。
