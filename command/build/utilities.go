package build

import "path/filepath"

func normalizePath(f, rootDir string) (string, error) {
	d, err := filepath.Abs(rootDir)

	if err != nil {
		return "", err
	}

	f, err = filepath.Abs(f)

	if err != nil {
		return "", err
	}

	return filepath.Rel(d, f)
}
