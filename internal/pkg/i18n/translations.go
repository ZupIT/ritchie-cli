// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// resources/i18n/en.toml
// resources/i18n/pt_BR.toml
package i18n

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _resourcesI18nEnToml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8a\x2e\x4f\xcd\x49\xce\xcf\x4d\xd5\x2b\xca\x2c\x49\xce\xc8\x4c\x8d\xe5\xca\x2f\xc9\x48\x2d\x52\xb0\x55\x50\x0a\x87\xc8\x28\x94\xe4\x2b\x04\x41\x24\x15\x95\x00\x01\x00\x00\xff\xff\x44\x64\x44\xad\x2f\x00\x00\x00")

func resourcesI18nEnTomlBytes() ([]byte, error) {
	return bindataRead(
		_resourcesI18nEnToml,
		"resources/i18n/en.toml",
	)
}

func resourcesI18nEnToml() (*asset, error) {
	bytes, err := resourcesI18nEnTomlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/i18n/en.toml", size: 47, mode: os.FileMode(420), modTime: time.Unix(1610136932, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesI18nPt_brToml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8a\x2e\x4f\xcd\x49\xce\xcf\x4d\xd5\x2b\xca\x2c\x49\xce\xc8\x4c\x8d\xe5\xca\x2f\xc9\x48\x2d\x52\xb0\x55\x50\x72\x4a\xcd\xd5\x2d\xcb\xcc\x4b\xc9\x57\x48\xcc\x57\x08\x82\x48\x2b\x2a\x01\x02\x00\x00\xff\xff\xf6\x61\xbf\x0c\x31\x00\x00\x00")

func resourcesI18nPt_brTomlBytes() ([]byte, error) {
	return bindataRead(
		_resourcesI18nPt_brToml,
		"resources/i18n/pt_BR.toml",
	)
}

func resourcesI18nPt_brToml() (*asset, error) {
	bytes, err := resourcesI18nPt_brTomlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/i18n/pt_BR.toml", size: 49, mode: os.FileMode(420), modTime: time.Unix(1610136932, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"resources/i18n/en.toml":    resourcesI18nEnToml,
	"resources/i18n/pt_BR.toml": resourcesI18nPt_brToml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"resources": &bintree{nil, map[string]*bintree{
		"i18n": &bintree{nil, map[string]*bintree{
			"en.toml":    &bintree{resourcesI18nEnToml, map[string]*bintree{}},
			"pt_BR.toml": &bintree{resourcesI18nPt_brToml, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
