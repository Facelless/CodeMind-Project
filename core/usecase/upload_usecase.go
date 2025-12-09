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
	ReadService files.ReadFileService
}

func NewUploadUsecase(f files.FileService, z files.ZipService, r files.ReadFileService) *UploadUsecase {
	return &UploadUsecase{
		FileService: f,
		ZipService:  z,
		ReadService: r,
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

	contentFile, err, _ := u.ReadService.OpenFile(files)
	if err != nil {
		return nil, pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "Error open file",
				"err":     err,
			},
		}
	}

	return nil, pkg.Log{
		Error: false,
		Body: map[string]any{
			"message": "File extracted successfully.",
			"files":   files,
			"content": contentFile,
		},
	}
}
