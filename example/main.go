package main

import (
	"fmt"
	"github.com/dayeguilaiye/file-generator/core"
	"github.com/dayeguilaiye/file-generator/generators"
	"github.com/dayeguilaiye/file-generator/handler"
)

func main() {
	g := generators.DefaultGenerator()
	err := g.Generate("example/gitIgnore_result", testNode)
	if err != nil {
		fmt.Printf("failed to generate file, err: %v", err)
	} else {
		fmt.Println("success")
	}
}

// result: a.tar.gz/bDir/c.sh
// c.sh is copied from example/template.sh, and replaced "[NEED_REPLACE]" into "c"
var testNode = core.Node{
	Type: "tgz",
	Data: handler.TgzParams{
		Name: "a.tar.gz",
		Children: []core.Node{
			{
				Type: "dir",
				Data: handler.DirParams{
					Name: "bDir",
					Children: []core.Node{
						{
							Type: "replace",
							Data: handler.ReplaceParams{
								Name:         "c.sh",
								FileMode:     0777,
								TemplatePath: "example/template.sh",
								Replaces:     map[string]string{"[NEED_REPLACE]": "c"},
							},
						},
						{
							Type: "goTemplate",
							Data: handler.GoTemplateParam{
								Name:         "classroom.txt",
								FileMode:     0777,
								TemplatePath: "example/classroom.gotemplate",
								Interface:    classRoom,
							},
						},
					},
				},
			},
		},
	},
}

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

type ClassRoom struct {
	Students []Student
}

type Student struct {
	Name string
	Age  string
}
