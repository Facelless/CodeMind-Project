package file

import (
	"archive/zip"
	"io"
	"miservicegolang/core/domain/files"
	"miservicegolang/core/pkg"
	"os"
	"path/filepath"
)

type UnzipService struct{}

func NewLocalUnzipService() files.ZipService {
	return &UnzipService{}
}

func (z *UnzipService) Unzip(src string, dst string) ([]string, error, pkg.Log) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return nil, nil, pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "Error reading zip file.",
			},
		}
	}
	defer r.Close()
	os.MkdirAll(dst, 0755)
	var files []string

	for _, f := range r.File {
		fpath := filepath.Join(dst, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
			continue
		}

		os.MkdirAll(filepath.Dir(fpath), 0755)
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, nil, pkg.Log{
				Error: true,
				Body: map[string]any{
					"message": "Error creating file",
					"err":     err,
				},
			}
		}
		rc, err := f.Open()
		if err != nil {
			return nil, nil, pkg.Log{
				Error: true,
				Body: map[string]any{
					"message": "Error open file.",
					"err":     err,
				},
			}
		}
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return nil, nil, pkg.Log{
				Error: true,
				Body: map[string]any{
					"message": "Erro copy file.",
					"err":     err,
				},
			}
		}
		files = append(files, fpath)
	}
	return files, nil, pkg.Log{}
}
