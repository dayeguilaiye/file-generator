package handler

import (
	"bytes"
	"github.com/dayeguilaiye/file-generator/core"
	"github.com/pkg/errors"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"text/template"
)

const TypeGoTemplate = "goTemplate"

type GoTemplateHandler struct {
	defaultFileMode fs.FileMode
}

func NewGoTemplateHandler() *GoTemplateHandler {
	return &GoTemplateHandler{}
}

func (g *GoTemplateHandler) SetDefaultFileMode(mode fs.FileMode) {
	g.defaultFileMode = mode
}

type GoTemplateParam struct {
	Name         string
	FileMode     fs.FileMode
	Interface    interface{}
	TemplatePath string
}

func (g *GoTemplateHandler) GetHandleType() string {
	return TypeGoTemplate
}

func (g *GoTemplateHandler) GetHandlerFunc() core.HandlerFunc {
	return func(_ *core.Generator, dir string, data interface{}) error {
		params, ok := data.(GoTemplateParam)
		if !ok {
			return core.WrongDataTypeError
		}

		if params.FileMode == 0 {
			params.FileMode = g.defaultFileMode
		}

		target := filepath.Join(dir, params.Name)
		tmpl, err := template.ParseFiles(params.TemplatePath)
		if err != nil {
			return errors.WithMessagef(err, "failed to parse file %s into gotemplate", params.TemplatePath)
		}
		content := &bytes.Buffer{}
		err = tmpl.Execute(content, params.Interface)
		if err != nil {
			return errors.WithMessagef(err, "failed to execute gotemplate %s", params.TemplatePath)
		}
		err = ioutil.WriteFile(target, content.Bytes(), params.FileMode)
		if err != nil {
			return errors.WithMessagef(err, "failed to write file %s", target)
		}
		return nil
	}
}
