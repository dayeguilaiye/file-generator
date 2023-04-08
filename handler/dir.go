package handler

import (
	core2 "github.com/dayeguilaiye/file-generator/core"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"path/filepath"
)

const TypeDir = "dir"

type DirHandler struct {
	defaultFileMode fs.FileMode
}

func NewDirHandler() *DirHandler {
	return &DirHandler{}
}

func (d *DirHandler) SetDefaultFileMode(mode fs.FileMode) {
	d.defaultFileMode = mode
}

type DirParams struct {
	Name     string
	FileMode fs.FileMode
	Children []core2.Node
}

func (d *DirHandler) GetHandleType() string {
	return TypeDir
}

func (d *DirHandler) GetHandlerFunc() core2.HandlerFunc {
	return func(generator *core2.Generator, targetDir string, data interface{}) error {
		params, ok := data.(DirParams)
		if !ok {
			return core2.WrongDataTypeError
		}

		if params.FileMode == 0 {
			params.FileMode = d.defaultFileMode
		}

		target := filepath.Join(targetDir, params.Name)
		if err := os.MkdirAll(target, params.FileMode); err != nil {
			return errors.WithMessagef(err, "failed to make dir %s", target)
		}
		for i, child := range params.Children {
			err := generator.Generate(target, child)
			if err != nil {
				return errors.WithMessagef(err, "failed to generate child of %s, index %d", params.Name, i)
			}
		}
		return nil
	}
}
