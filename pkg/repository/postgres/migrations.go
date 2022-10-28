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

	info := bindataFileInfo{name: "1_initialize_schema.down.sql", size: 18, mode: os.FileMode(436), modTime: time.Unix(1666964845, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1_initialize_schemaUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8d\x31\x4b\xc7\x30\x10\x47\xf7\x7c\x8a\x1b\xff\x05\x27\xa1\x93\x53\xd4\x08\x62\xac\x52\x52\xb0\x53\x39\x92\x80\xc1\x26\x0d\xc9\xc5\xc1\x4f\x2f\x6d\xa1\x36\x20\x62\xb6\x5f\x78\xef\xde\x5d\x2f\xb8\x12\xa0\xf8\xad\x14\x40\x98\x3f\x32\xbb\x30\x00\x00\x67\xa0\x7a\xa5\xd4\x3f\xaf\xfd\xe3\x33\xef\x47\x78\x12\xe3\xd5\x26\x04\xf4\xf6\x0c\x7c\x62\xd2\xef\x98\x2e\xd7\x6d\xdb\xac\xbb\x7b\x51\xd0\x0d\x52\xee\xb4\xb1\x59\x27\x17\xc9\x2d\x61\xa3\x95\x78\x53\x67\xfb\x5e\x3c\xf0\x41\xd6\x06\x9a\xd9\x85\xa3\x61\x90\xaa\x5e\x7d\x5f\x2f\x3e\xce\x96\xac\x99\x90\xd6\x4d\xce\xdb\x4c\xe8\x23\x7d\xfd\x12\x80\xdd\x49\x16\x7f\x8c\x3f\x1c\x5d\x52\xb2\x81\xa6\x03\xd8\x9b\x25\x9a\xff\xf9\x6b\x93\x35\x37\xec\x3b\x00\x00\xff\xff\xc9\xdc\x56\xfb\x7d\x01\x00\x00")

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

	info := bindataFileInfo{name: "1_initialize_schema.up.sql", size: 381, mode: os.FileMode(436), modTime: time.Unix(1666969780, 0)}
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
