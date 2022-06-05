package replacefile

import (
	"path/filepath"
)

func RealPath(path string) (string,error) {
	return filepath.EvalSymlinks(path)
}
