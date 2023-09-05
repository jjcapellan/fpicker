package fpicker

import (
	"path/filepath"
	"strings"
)

func getRootName(path string) string {
	if path == home {
		return "home"
	}
	for _, drive := range drives {
		if drive.Path == path {
			return drive.Name
		}
	}

	return ""
}

func makeBreadcrumb(root string, path string) []string {
	bc := []string{}

	// Checks to avoid infinite loops
	if !strings.Contains(path, root) {
		return bc
	}
	if path == root {
		bc = append(bc, getRootName(root))
		return bc
	}
	bc = append(bc, filepath.Base(path))

	parent := path
	for {

		parent = filepath.Dir(parent)
		parent = filepath.ToSlash(parent)
		if parent == root {
			bc = append(bc, getRootName(root))
			break
		}
		bc = append(bc, filepath.Base(parent))

	}
	return reverseArray(bc)
}

func pathToAbs(path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(currentDir, path)
	}
	return path
}

func reverseArray(arr []string) []string {
	length := len(arr)
	reversed := make([]string, length)

	for i := 0; i < length; i++ {
		reversed[i] = arr[length-i-1]
	}

	return reversed
}
