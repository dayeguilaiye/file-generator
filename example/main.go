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
