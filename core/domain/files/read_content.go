package files

import (
	"miservicegolang/core/pkg"
)

type ReadFileService interface {
	OpenFile(file []string) (string, error, pkg.Log)
}
