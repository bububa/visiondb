package utils

import (
	"os/user"
	"path/filepath"
	"strings"
)

// PathJoin returns joined paths
func PathJoin(dirs ...string) string {
	l := make([]string, 0, len(dirs))
	for _, path := range dirs {
		l = append(l, normalizeUserDir(path))
	}
	return filepath.Join(l...)
}

func normalizeUserDir(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if path == "~" {
		return dir
	} else if strings.HasPrefix(path, "~/") {
		return filepath.Join(dir, path[2:])
	}
	return path
}
