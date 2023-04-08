package handler

import (
	core2 "github.com/dayeguilaiye/file-generator/core"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

const TypeCopy = "copy"

type CopyHandler struct {
	bufferSize int
}

func NewCopyHandler() *CopyHandler {
	return &CopyHandler{
		bufferSize: 32 * 1024,
	}
}

func (c *CopyHandler) SetBufferSize(size int) {
	c.bufferSize = size
}

type CopyParams struct {
	Name string
	Src  string
}

func (c *CopyHandler) GetHandleType() string {
	return TypeCopy
}

func (c *CopyHandler) GetHandlerFunc() core2.HandlerFunc {
	return func(generator *core2.Generator, targetDir string, data interface{}) error {
		params, ok := data.(CopyParams)
		if !ok {
			return core2.WrongDataTypeError
		}

		target := filepath.Join(targetDir, params.Name)
		if err := doCopy(params.Src, target, c.bufferSize); err != nil {
			return errors.WithMessagef(err, "failed to do copy from %s to %s", params.Src, target)
		}
		return nil
	}
}

func doCopy(src, dst string, bufferSize int) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if srcInfo.IsDir() {
		return copyDir(src, dst, bufferSize)
	}
	return copyFile(src, dst, bufferSize)
}

func copyDir(src, dst string, bufferSize int) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath, bufferSize); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath, bufferSize); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string, bufferSize int) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	buffer := make([]byte, bufferSize)
	if _, err = io.CopyBuffer(out, in, buffer); err != nil {
		return err
	}
	err = out.Sync()
	return err
}
