// Code generated for package postgres by go-bindata DO NOT EDIT. (@generated)
// sources:
// pkg/repository/postgres/migrations/1_initialize_schema.down.sql
// pkg/repository/postgres/migrations/1_initialize_schema.up.sql
package postgres

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

var __1_initialize_schemaDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\x28\x49\x2c\xce\x2e\xb6\xe6\x02\x04\x00\x00\xff\xff\xdc\xde\x52\xad\x12\x00\x00\x00")

func _1_initialize_schemaDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__1_initialize_schemaDownSql,
		"1_initialize_schema.down.sql",
	)
}

func _1_initialize_schemaDownSql() (*asset, error) {
	bytes, err := _1_initialize_schemaDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1_initialize_schema.down.sql", size: 18, mode: os.FileMode(436), modTime: time.Unix(1668418308, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1_initialize_schemaUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8d\xb1\x4a\xc6\x30\x14\x85\xf7\x3c\xc5\x1d\xff\x82\x93\xd0\xc9\x29\x6a\x04\x31\x56\x29\x29\xd8\xa9\x5c\x92\x80\xc1\x26\x0d\xc9\x8d\x83\x4f\x2f\x6d\xa1\x36\x88\x60\xb6\x93\x7b\xbe\xf3\xdd\xf5\x82\x2b\x01\x8a\xdf\x4a\x01\x84\xf9\x23\xb3\x0b\x03\x00\x70\x06\xaa\x57\x4a\xfd\xf3\xda\x3f\x3e\xf3\x7e\x84\x27\x31\x5e\x6d\x40\x40\x6f\xcf\x85\x4f\x4c\xfa\x1d\xd3\xe5\xba\x6d\x9b\x35\x77\x2f\x0a\xba\x41\xca\xbd\x6d\x6c\xd6\xc9\x45\x72\x4b\xd8\xda\x4a\xbc\xa9\x33\x7d\x2f\x1e\xf8\x20\x6b\x02\xcd\xec\xc2\xe1\x30\x48\x95\xaf\xde\xd7\x8b\x8f\xb3\x25\x6b\x26\xa4\x35\x93\xf3\x36\x13\xfa\x48\x5f\xbf\xf7\x61\x47\x92\xc5\x1f\xe0\x6f\x44\x97\x94\x6c\xa0\xe9\xb8\xef\xc6\x12\xcd\xbf\xf0\xd5\xc8\x9a\x1b\xf6\x1d\x00\x00\xff\xff\x3f\xb0\x90\x6b\x7a\x01\x00\x00")

func _1_initialize_schemaUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1_initialize_schemaUpSql,
		"1_initialize_schema.up.sql",
	)
}

func _1_initialize_schemaUpSql() (*asset, error) {
	bytes, err := _1_initialize_schemaUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1_initialize_schema.up.sql", size: 378, mode: os.FileMode(436), modTime: time.Unix(1668593892, 0)}
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
	"1_initialize_schema.down.sql": _1_initialize_schemaDownSql,
	"1_initialize_schema.up.sql":   _1_initialize_schemaUpSql,
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
	"1_initialize_schema.down.sql": &bintree{_1_initialize_schemaDownSql, map[string]*bintree{}},
	"1_initialize_schema.up.sql":   &bintree{_1_initialize_schemaUpSql, map[string]*bintree{}},
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
