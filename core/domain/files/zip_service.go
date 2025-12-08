package files

import "miservicegolang/core/pkg"

type ZipService interface {
	Unzip(src string, dst string) ([]string, error, pkg.Log)
}
