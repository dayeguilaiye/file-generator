English | [简体中文](README.zh-CN.md)

# File Generator

This library is used to generate files conveniently.

If you want to generate a file, and the file contains multiple directories or even multiple layers of compressed files, such as `.tar.gz` and `.zip`, you will want a tool like this.

# use example
You need to generate a complex zip file:

The outermost layer is a compressed file named `a.tar.gz`, which contains a directory named `bDir`, which contains a file named `c.sh`, `[NEED_REPLACE]` in this file needs to be in different In addition, there is a `classroom.txt` file in `bDir`, and the names and ages of the whole class are to be printed out.

If you manually create this file every time, you need to decompress the file, modify the file, and then pack and compress it. In some cases, you may encounter more troublesome situations than this, such as nesting of multi-layer compressed files.

So we need a tool to help us generate this file.

Taking the above requirement as an example, we can use this library to generate the file.

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

There are already detailed comments in the code, so I won't repeat them here. Here's what each type means.
#Node
In this repository, `Node` is an abstract concept that represents a file or directory. There are two properties inside this node, namely `Type` and `Data`.

## Node.Type
Node.Type is used to identify the type of the node, which is not limited to files and folders. For example, the types supported by default in this warehouse include `copy`, `dir`, `gotemplate`, `replace`, `tgz`.
- Nodes of type `copy` will copy the file or folder specified in `Data` directly to the target path;
- A node of type `dir` creates a directory;
- The `gotemplate` type node will generate a file according to the `go template` template file and data specified in `Data`;
- `replace` type of node will generate a file according to the template file and data specified in `Data`, but the content in the file will be replaced with the content specified in `Data`;
- Nodes of type `tgz` will pack the specified subfiles or directories in `Data` into a `tar.gz` file.

## Node.Data
`Node.Data` is an `interface{}` type, its content is directly related to `Node.Type`, and different `Node.Type` will have different `Data` types.
If `Node.Type == dir`, then the corresponding `Node.Data` should store the following structure:

```go
type DirParams struct {
    Name     string
    FileMode fs.FileMode
    Children []core2.Node
}
```
Among them, `Name` represents the name of the directory, `FileMode` represents the permissions of the directory, and `Children` represents the child nodes under the directory.

# Handler
Handler is a function used to handle different types of Node, and its definition is as follows:
```go
type HandlerFunc func(generator *Generator, targetDir string, data interface{}) error

type Handler interface {
    GetHandleType() string
    GetHandlerFunc() HandlerFunc
}
```
This repository provides some default `Handler`, as mentioned above, supports `copy`, `dir`, `gotemplate`, `replace`, `tgz` types of `Handler`.
Users can also refer to the provided `Handler` to customize `Handler`.

# Generator
`Generator` is the executor of the library, and instances of this class are used to generate files. `DefaultGenerator` has some built-in handlers to handle different file types.
Users can also add custom handlers through `generator.Handle(handler)`.

## more examples

See the [example](example) directory for more usage information.