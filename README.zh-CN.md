[English](README.md) | 简体中文

# File Generator

该库用于方便地生成文件。

如果你想要生成一个文件，而且该文件包含了多层目录甚至多层的压缩文件，比如 `.tar.gz` 和 `.zip`，你就会想要一个这样的工具。

# 使用示例
你需要生成一个复杂的压缩文件：

最外层是个名为 `a.tar.gz` 的压缩文件，该文件包含一个名为 bDir 的目录，该目录下包含一个名为 c.sh 的文件，该文件中的 `[NEED_REPLACE]` 需要在不同情况下更换值; 此外，`bDir` 内还有一个 `classroom.txt` 文件，要打印出全班的名称和年龄。

如果每次都手动制作这个文件，需要解压该文件，修改文件，然后再打包压缩，某些情况下也许会遇到比这还要麻烦的情况，比如多层压缩文件嵌套。

所以我们需要一个工具来帮我们生成这个文件。

以上述需求为例，我们可以使用该库来生成该文件。

```go
package main

import (
	"fmt"
	"github.com/dayeguilaiye/file-generator/core"
	"github.com/dayeguilaiye/file-generator/generators"
	"github.com/dayeguilaiye/file-generator/handler"
)

func main() {
	// set some data
	var classRoom = ClassRoom{
		Students: []Student{
			{
				Name: "Sam",
				Age:  "10",
			},
			{
				Name: "Amy",
				Age:  "12",
			},
		},
	}

	err := GenerateMyStrangeFile(classRoom, "echo hello world")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success")
	}
}

// some data structure

type Student struct {
	Name string
	Age  string
}

type ClassRoom struct {
	Students []Student
}

// GenerateMyStrangeFile is an example to generate a file with extremely deep layers,
// and the data in file have to change with the input data (classroom).
func GenerateMyStrangeFile(room ClassRoom, roomName string) error {
	// set node structure
	node := core.Node{
		// a tgz type file with name "a.tar.gz"
		Type: "tgz",
		Data: handler.TgzParams{
			Name: "a.tar.gz",
			// the a.tar.gz contains a directory with name "bDir"
			Children: []core.Node{
				{
					Type: "dir",
					Data: handler.DirParams{
						Name: "bDir",
						// in bDir, there are two files, one is c.sh, the other is classroom.txt
						Children: []core.Node{
							{
								// c.sh is a file with content "echo [NEED_REPLACE]",
								// and the "[NEED_REPLACE]" will be replaced with the input data (roomName)
								Type: "replace",
								Data: handler.ReplaceParams{
									Name:         "c.sh",
									FileMode:     0777,
									TemplatePath: "example/template.sh",
									Replaces:     map[string]string{"[NEED_REPLACE]": roomName},
								},
							},
							{
								// classroom.txt is a file with content "
								// {{ range .Students }}
								//    My name is {{ .Name }}, my age is {{ .Age }}.
								// {{ end }}
								//", it is wrote by go template, and the data in file will be replaced with the input data (room)
								Type: "goTemplate",
								Data: handler.GoTemplateParam{
									Name:         "classroom.txt",
									FileMode:     0777,
									TemplatePath: "example/classroom.gotemplate",
									Interface:    room,
								},
							},
						},
					},
				},
			},
		},
	}
	// get the default generator
	g := generators.DefaultGenerator()
	// use the generator generate the node
	return g.Generate("example/gitIgnore_result", node)
}

```

代码中已经有了详细的注释，这里就不再赘述。以下是各个类型的含义。
# Node
在本仓库中，`Node` 是一个抽象的概念，它代表了一个文件或者目录。该节点内部有两个属性，分别是 `Type` 和 `Data`。

## Node.Type
Node.Type 用于标识该节点的类型，该类型并不局限于文件和文件夹，比如本仓库默认支持的类型包括 `copy`, `dir`, `gotemplate`, `replace`, `tgz`。
- `copy` 类型的节点会将 `Data` 中指定的文件或文件夹直接复制到目标路径中；
- `dir` 类型的节点会创建一个目录；
- `gotemplate` 类型的节点会根据 `Data` 中指定的 `go template` 模板文件和数据生成一个文件；
- `replace` 类型的节点会根据 `Data` 中指定的模板文件和数据生成一个文件，但是该文件中的内容会被替换为 `Data` 中指定的内容；
- `tgz` 类型的节点会将 `Data` 中指定子文件或目录打包成一个 `tar.gz` 文件。

## Node.Data
`Node.Data` 是个 `interface{}` 类型，它的内容与 `Node.Type` 直接相关，不同的 `Node.Type` 会有不同的 `Data` 类型。
如 `Node.Type == dir` ，则对应的 Node.Data 中存放的就应该是如下的结构体：

```go
type DirParams struct {
	Name     string
	FileMode fs.FileMode
	Children []core2.Node
}
```
其中，`Name代表了该目录的名称`，`FileMode` 代表了该目录的权限，`Children` 代表了该目录下的子节点。

# Handler
Handler为用于处理不同类型的 Node 的函数，它的定义如下：
```go
type HandlerFunc func(generator *Generator, targetDir string, data interface{}) error

type Handler interface {
	GetHandleType() string
	GetHandlerFunc() HandlerFunc
}
```
本仓库提供了一些默认的 `Handler`，如上文所说，支持 `copy`, `dir`, `gotemplate`, `replace`, `tgz` 类型的 `Handler`。
用户也可参考已提供的 `Handler` 来自定义 `Handler`。

# Generator
`Generator` 是该库的执行器，由该类的实例来执行生成文件的操作。`DefaultGenerator` 内置了一些 handler，用来处理不同的文件类型。
用户也可以通过 `generator.Handle(handler)` 来添加自定义的 handler。

## 更多示例

查看 [example](example) 目录以获取更多用法信息。