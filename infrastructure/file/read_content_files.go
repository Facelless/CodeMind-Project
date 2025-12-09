package file

import (
	"miservicegolang/core/domain/files"
	"miservicegolang/core/pkg"
	"os"
)

type LocalReadFileService struct{}

func NewLocalReadFileService() files.ReadFileService {
	return &LocalReadFileService{}
}

func (l *LocalReadFileService) OpenFile(files []string) (string, error, pkg.Log) {
	var contentFile []string
	for _, f := range files {
		data, err := os.ReadFile(f)
		contentFile = append(contentFile, string(data))
		if err != nil {
			return "", nil, pkg.Log{
				Error: true,
				Body: map[string]any{
					"message": "Error read file: ",
					"err":     err,
				},
			}
		}
	}
	return contentFile[0], nil, pkg.Log{}
}
