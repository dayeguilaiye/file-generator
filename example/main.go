package main

import (
	"fmt"
	"github.com/dayeguilaiye/file-generator/core"
	"github.com/dayeguilaiye/file-generator/generators"
	"github.com/dayeguilaiye/file-generator/handler"
)

func main() {
	g := generators.DefaultGenerator()
	node := core.Node{
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
									TemplatePath: "./template.sh",
									Replaces:     map[string]string{"[NEED_REPLACE]": "c"},
								},
							},
						},
					},
				},
			},
		},
	}
	err := g.Generate("./", node)
	if err != nil {
		fmt.Printf("failed to generate file, err: %v", err)
	} else {
		fmt.Println("success")
	}
}
