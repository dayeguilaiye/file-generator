package handler

import (
	core2 "github.com/dayeguilaiye/file-generator/core"
	"github.com/pkg/errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const TypeReplace = "replace"

type ReplaceHandler struct{}

func NewReplaceHandler() *ReplaceHandler {
	return &ReplaceHandler{}
}

type ReplaceParams struct {
	Name         string
	FileMode     fs.FileMode
	TemplatePath string
	Replaces     map[string]string
}

func (r *ReplaceHandler) GetHandleType() string {
	return TypeReplace
}

func (r *ReplaceHandler) GetHandlerFunc() core2.HandlerFunc {
	return func(generator *core2.Generator, targetDir string, data interface{}) error {
		params, ok := data.(ReplaceParams)
		if !ok {
			return core2.WrongDataTypeError
		}

		if params.FileMode == 0 {
			srcFile, err := os.Stat(params.TemplatePath)
			if err != nil {
				return errors.WithMessagef(err, "failed to check stat of file %s", params.TemplatePath)
			}
			params.FileMode = srcFile.Mode()
		}

		target := filepath.Join(targetDir, params.Name)
		srcContent, err := ioutil.ReadFile(params.TemplatePath)
		if err != nil {
			return errors.WithMessagef(err, "failed to read file %s", params.TemplatePath)
		}
		dstContent := string(srcContent)
		for k, v := range params.Replaces {
			dstContent = strings.ReplaceAll(dstContent, k, v)
		}
		err = ioutil.WriteFile(target, []byte(dstContent), params.FileMode)
		if err != nil {
			return errors.WithMessagef(err, "failed to write file %s", target)
		}
		return nil
	}
}
