package file

import (
	"io"
	"mime/multipart"
	"miservicegolang/core/domain/files"
	"miservicegolang/core/pkg"
	"os"
	"path/filepath"
)

type LocalFileService struct{}

func NewLocalFileService() files.FileService {
	return &LocalFileService{}
}

func (l *LocalFileService) SaveTemp(file multipart.File) (string, error, pkg.Log) {
	tempPath := filepath.Join(os.TempDir(), "upload.zip")
	out, err := os.Create(tempPath)

	if err != nil {
		return "", nil, pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "Error creating folder.",
			},
		}
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return "", nil, pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "Error copy folder.",
			},
		}
	}

	return tempPath, nil, pkg.Log{}
}
