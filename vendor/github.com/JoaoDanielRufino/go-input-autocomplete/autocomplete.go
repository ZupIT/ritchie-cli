package input_autocomplete

import (
	"runtime"
	"strings"
)

type autocomplete struct {
	cmd DirListChecker
}

func Autocomplete(path string) string {
	os := runtime.GOOS
	a := autocomplete{
		cmd: Cmd{},
	}

	switch os {
	case "linux", "darwin":
		return a.unixAutocomplete(path)
	case "windows":
		return a.windowsAutocomplete(path)
	default:
		return path
	}
}

// Return if the string starts with prefix, case insensitive
func hasInsensitivePrefix(s string, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[0:len(prefix)], prefix)
}

func (a autocomplete) unixAutocomplete(path string) string {
	if path == "" || path[len(path)-1] == ' ' {
		return path
	}

	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 || (path[0] != '/' && path[:2] != "./") {
		path = "./" + path
		if !strings.Contains(path, "/") {
			lastSlash = 1
		} else {
			lastSlash = strings.LastIndex(path, "/")
		}
	}

	path = a.findFromPrefix(path, lastSlash)
	ok, err := a.cmd.IsDir(path)
	if ok && err == nil {
		path = path + "/"
	}

	return path
}

func (a autocomplete) windowsAutocomplete(path string) string {
	if path == "" || path[len(path)-1] == ' ' {
		return path
	}

	lastSlash := strings.LastIndex(path, "\\")
	if !strings.Contains(path, ":") && !strings.Contains(path, ".\\") {
		path = ".\\" + path
		if !strings.Contains(path, "\\") {
			lastSlash = 1
		} else {
			lastSlash = strings.LastIndex(path, "\\")
		}
	}

	path = a.findFromPrefix(path, lastSlash)
	ok, err := a.cmd.IsDir(path)
	if ok && err == nil {
		path = path + "\\"
	}

	return path
}

func (a autocomplete) findFromPrefix(prefix string, lastSlash int) string {
	contents, err := a.cmd.ListContent(prefix[:lastSlash+1])
	if err != nil {
		return prefix
	}

	for _, content := range contents {
		if hasInsensitivePrefix(content, prefix[lastSlash+1:]) {
			return prefix[:lastSlash+1] + content
		}
	}

	return prefix
}
