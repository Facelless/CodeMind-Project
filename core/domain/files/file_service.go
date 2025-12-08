package files

import (
	"mime/multipart"
	"miservicegolang/core/pkg"
)

type FileService interface {
	SaveTemp(file multipart.File) (string, error, pkg.Log)
}
