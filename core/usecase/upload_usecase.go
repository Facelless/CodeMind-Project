package usecase

import (
	"mime/multipart"
	"miservicegolang/core/domain/files"
	"miservicegolang/core/pkg"
	"os"
	"path/filepath"
)

type UploadUsecase struct {
	FileService files.FileService
	ZipService  files.ZipService
}

func NewUploadUsecase(f files.FileService, z files.ZipService) *UploadUsecase {
	return &UploadUsecase{
		FileService: f,
		ZipService:  z,
	}
}

func (u *UploadUsecase) Execute(file multipart.File) (error, pkg.Log) {
	tempFile, err, _ := u.FileService.SaveTemp(file)
	if err != nil {
		return nil, pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "Error saving file.",
				"err":     err,
			},
		}
	}
	outputDir := filepath.Join(os.TempDir(), "project_unzipped")
	files, err, _ := u.ZipService.Unzip(tempFile, outputDir)
	if err != nil {
		return nil, pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "Error extracting folder.",
				"err":     err,
			},
		}
	}

	return nil, pkg.Log{
		Error: false,
		Body: map[string]any{
			"message": "File extracted successfully.",
			"files":   files,
		},
	}
}
