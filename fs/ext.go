package fs

import (
	"path/filepath"
	"strings"
)

func Ext(path string) string {
	return filepath.Ext(path)
}
func DotExt(ext string) string {
	if ext == "" {
		return ""
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return ext
}

func AddExt(path string, ext string) string {
	return path + DotExt(ext)
}

func RemoveExt(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

func ReplaceExt(path string, ext string) string {
	return AddExt(RemoveExt(path), ext)
}

func EqualExt(extA, extB string) bool {
	return DotExt(extA) == DotExt(extB)
}
