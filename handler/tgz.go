package handler

import (
	"archive/tar"
	"compress/gzip"
	core2 "github.com/dayeguilaiye/file-generator/core"
	"github.com/pkg/errors"

	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const TypeTgz = "tgz"

type TgzHandler struct {
	defaultFileMode fs.FileMode
	TmpDir          string
}

func NewTgzHandler() *TgzHandler {
	return &TgzHandler{
		defaultFileMode: 0666,
		TmpDir:          "./",
	}
}

type TgzParams struct {
	Name     string
	FileMode fs.FileMode
	Children []core2.Node
}

func (t *TgzHandler) GetHandleType() string {
	return TypeTgz
}

func (t *TgzHandler) GetHandlerFunc() core2.HandlerFunc {
	return func(generator *core2.Generator, targetDir string, data interface{}) error {
		params, ok := data.(TgzParams)
		if !ok {
			return core2.WrongDataTypeError
		}
		if params.FileMode == 0 {
			params.FileMode = t.defaultFileMode
		}
		if err := recreateDir(t.TmpDir); err != nil {
			return errors.WithMessagef(err, "failed to recreate dir %s", t.TmpDir)
		}
		for i, child := range params.Children {
			err := generator.Generate(t.TmpDir, child)
			if err != nil {
				return errors.WithMessagef(err, "failed to generate child of %s, index: %d", params.Name, i)
			}
		}
		if err := compressFilesInDirToTgz(t.TmpDir, filepath.Join(targetDir, params.Name), params.FileMode); err != nil {
			return errors.WithMessagef(err, "failed to compress files in dir %s to tgz file", t.TmpDir)
		}

		return nil
	}
}

// compressFilesInDirToTgz compress all file in srcDir(except the srcDir itself) to target tgz file
func compressFilesInDirToTgz(srcDir, target string, fileMode os.FileMode) error {
	// check if srcDir is a directory
	if fi, err := os.Stat(srcDir); err != nil {
		return err
	} else if !fi.IsDir() {
		return errors.New("srcDir is not a directory")
	}
	t, _ := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
	defer t.Close()
	gw := gzip.NewWriter(t)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == srcDir {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = path[len(srcDir)+1:]
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
		return nil
	})
}

// recreateDir create dir if not exist, delete and create if exist
func recreateDir(dir string) error {
	if _, err := os.Stat(dir); err == nil {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return os.MkdirAll(dir, 0777)
}

func (t *TgzHandler) SetDefaultFileMode(mode fs.FileMode) {
	t.defaultFileMode = mode
}
