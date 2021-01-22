// Code generated by go-bindata. DO NOT EDIT.
// sources:
// models/migrations/1489143364_initial_schema.down.sql (1.043kB)
// models/migrations/1489143364_initial_schema.up.sql (3.849kB)

package migrations

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __1489143364_initial_schemaDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x92\xcd\x6a\x85\x30\x10\x85\xf7\x79\x8a\xbc\x40\x9e\xc0\x55\x4b\x2d\x08\x85\x4a\x95\xfe\xac\x44\xd2\x51\x83\xe9\xa4\x98\x64\xe1\xdb\x97\x44\x7b\xef\x95\x9b\x04\x45\x97\x73\xe6\x9c\x8f\x9c\x91\x31\xfa\x5e\xe4\x1f\x15\x21\x8c\x51\x16\xfc\x9c\xc2\x87\x56\x4a\xc0\x1e\x74\x33\x08\x34\x3a\xb1\x4d\x9e\xde\x5e\x4b\x1f\x4a\x8b\x67\x9a\x7f\x16\x55\x5d\xdd\xf9\xb3\xfd\x3c\x29\x46\x38\xc3\xf3\xfe\x03\x3c\xad\xa4\x35\x42\xe1\x19\xe6\x25\x63\xe1\xd6\x0f\x8f\x2f\x79\xb0\x62\x37\xd2\x60\x8c\xc0\x3e\xc4\x5b\x41\xde\x7f\x43\xfa\x77\x64\xab\xfc\x55\x6e\x54\xae\x26\x81\x7d\x63\xe6\x5f\x08\xbe\xdb\x8d\x62\xad\x46\x88\xf1\x12\xdd\xa8\x45\x54\x16\x39\xfc\x40\xf8\xd7\x88\x84\x6e\x6c\xd1\xf0\xd4\x3d\x62\xfd\x6c\xeb\x4f\xde\x7b\x7f\xea\xd5\x13\xee\xfd\x5b\x74\x9d\xe0\x56\x9a\x39\x8a\x35\x6a\x84\x23\x0f\x59\xf6\xc3\x38\xaf\xa5\x8f\x6c\x35\x4c\x07\x68\x7e\x3d\x23\x7f\x01\x00\x00\xff\xff\x24\xb7\x34\x3a\x13\x04\x00\x00")

func _1489143364_initial_schemaDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__1489143364_initial_schemaDownSql,
		"1489143364_initial_schema.down.sql",
	)
}

func _1489143364_initial_schemaDownSql() (*asset, error) {
	bytes, err := _1489143364_initial_schemaDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1489143364_initial_schema.down.sql", size: 1043, mode: os.FileMode(0644), modTime: time.Unix(1594286535, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x8e, 0xc5, 0xc, 0x59, 0x81, 0x49, 0xd6, 0xdf, 0x50, 0x91, 0xa, 0x91, 0x96, 0x70, 0xcb, 0x90, 0xc1, 0x6f, 0x4e, 0x5, 0xd0, 0x91, 0x37, 0x32, 0x68, 0x9, 0x3, 0x30, 0x20, 0x18, 0x88, 0xae}}
	return a, nil
}

var __1489143364_initial_schemaUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x57\x51\x93\xb2\x36\x14\x7d\xe7\x57\xdc\x87\x9d\x51\x3b\xeb\x43\x9f\x77\xfa\x80\x7a\xd7\x8f\x56\xc1\x02\xf6\xfb\x9c\x4e\x87\x49\x21\xab\x99\x45\x60\x48\x6c\xeb\xbf\xef\x24\xe8\x92\x60\xa0\x75\xb7\xfb\x54\x9e\x22\x9c\xdc\xe4\xdc\x7b\x72\x6e\x9c\x4e\x21\x76\x67\x2b\x8c\x1c\x67\x3a\x85\xa9\xf9\xc8\x57\x27\x4e\x6b\x6e\xfb\xe6\xcc\x43\x74\x63\x6c\xa6\x83\xf7\x0c\x7e\x10\x03\x7e\xf3\xa2\x38\x72\xd4\x24\x18\x3b\x00\x2c\x03\xed\x89\x30\xf4\xdc\x95\x1c\x6d\x42\x6f\xed\x86\x3b\xf8\x09\x77\x8f\x0e\x40\x41\x8e\xb4\x85\xc5\xf8\x2d\x6e\x46\x32\xa6\xbf\x5d\xad\x60\xeb\x7b\x3f\x6f\x51\x42\xe9\x91\xb0\xfc\xdf\x41\x2b\xc2\xf9\x9f\x65\x9d\x25\x07\xc2\x0f\x16\xa8\xc4\x30\x9e\x90\x54\xb0\x3f\x88\xa0\x19\xc0\x2c\x08\x56\xe8\xfa\x1d\x0c\xfd\x4b\xd4\xe4\x6d\xc9\x1f\xa3\xc0\x9f\xdd\xc4\x49\x6b\x2a\x43\x24\x44\xa8\x6d\x79\x6b\x8c\x62\x77\xbd\x31\x30\xa7\x2a\x1b\xc4\x38\x93\xa7\xbe\x32\x88\xf2\x95\x16\xf6\x3a\x2c\x02\x78\x78\x70\x00\x66\xb8\xf4\x7c\x47\x46\x36\x8a\x01\xe3\x08\x57\x38\x8f\xe1\x7b\x78\x0e\x83\x35\x54\xfb\x44\x9c\x2b\x0a\x5f\xbf\x60\x88\x20\xce\x95\xca\xfd\x0f\x30\x52\x4b\xa8\x6f\xa3\x09\xc4\x5f\xb0\x89\x05\x70\xad\xf3\x6e\x83\xd0\x62\xc0\x8d\x00\xfd\xed\x7a\x3c\xba\xa6\x6f\xf4\x08\xa3\x9a\x72\x2a\x46\x93\x27\x35\x15\xfd\x05\x78\xcf\x72\x8c\xfe\xc2\x79\x78\x78\x1a\xd4\x4c\xc3\xb0\x2b\x9a\x56\x31\x5d\xc9\x48\x8d\x25\x0d\xd2\xf3\x63\x5c\x62\x68\xa8\x20\xc4\x67\x0c\xd1\x9f\x63\x04\x17\x35\xb2\x6c\x02\x81\x0f\x0b\x5c\x61\x8c\x30\x77\xa3\xb9\xbb\x50\x32\x51\x2b\x77\xe5\x64\x14\x4e\x11\x6e\x1e\x2d\x03\xa6\x44\x2a\x56\x53\x2e\x4b\xdb\xd6\xb5\x4f\x20\x16\xc4\x40\xe9\xd3\x03\xc9\x73\x5a\xec\xe9\xe7\x96\x3f\x63\x2f\x2f\x2c\x3d\xe5\xe2\xdc\x5f\xfe\x16\xd3\x96\x9f\x12\x7e\x96\xa5\x3f\xd2\x8c\x9d\x8e\x72\x74\x20\x75\xf6\x3e\x0d\xb4\x54\x6f\xcc\x43\x13\x82\x4d\x0d\x82\x89\x9c\xde\x78\xc2\x6d\x25\x53\x22\xe8\xbe\xac\x19\xe5\x0d\xee\xd7\xdf\xec\xb8\xaa\x64\x85\xe0\xcd\x7b\x4d\x5e\x37\xb8\x8c\xf2\xb4\x66\x95\x60\x65\x31\xb8\xae\x96\x39\x63\xdc\xc5\xbd\xe4\x64\xdf\xd8\xd5\x30\x0f\xc6\x93\xbc\x4c\x5f\xa5\x67\x69\xa6\x65\xdb\x5f\x45\x8b\x8c\x27\x65\x71\xcb\x43\x3b\x22\x7a\xe2\x59\x36\xe9\x1a\x9a\x26\xd8\x9b\x15\x74\x53\xeb\xc3\x0d\xa8\x9b\x97\xf9\x49\x66\xef\xee\x1e\xf3\x36\x51\x49\x45\x73\x03\x9d\xe8\x7b\xec\xe0\x2d\x17\x32\xde\x70\xa8\x4e\xda\x7a\xe2\xe9\x9d\xc1\xde\x18\x34\x31\xc3\xf8\x42\xe5\xd1\xd8\xc8\x64\x20\x85\xa4\x28\xca\x53\x91\xd2\x23\x2d\xc4\xdd\x69\x34\x26\xdf\xb6\xec\xde\x8e\xad\x1f\xb8\x9e\xd6\xfa\x7b\x99\x9d\x61\x18\xd2\x93\xeb\x4f\x48\xb1\xd1\x7b\xef\x6a\xbd\x39\x7b\xed\xb1\xde\x81\xb4\xaa\x49\xff\x7b\x65\x72\x2a\x04\x2b\xf6\x9f\xdb\xb8\x78\x5a\xd6\xac\xd8\xff\xc3\xcd\x45\x47\xb5\xcd\x2b\xcd\x09\xe7\x2c\x95\x5d\x2b\x3b\x17\xe4\xc8\xd2\xf7\x35\xae\x2b\x51\x55\x71\x00\x2e\x48\x2d\xf8\xe5\x9e\xd7\xb1\x45\x3d\xbb\x00\xca\x9c\xaf\xb8\x41\xa4\xde\xb9\xa0\xd3\xbc\x4c\xa4\x41\xb5\xf3\xd3\x8e\xac\x48\x4d\x8e\x5c\xbb\xd9\x5a\xbc\xfb\x17\x0f\xbf\x5a\xff\x25\x58\xaf\x2a\xc9\x90\xaf\x77\x0f\x50\x10\x42\x88\x9b\x95\x3b\x47\xb5\x8a\x35\x0c\xb8\x91\xd3\x08\xa2\x73\x10\x24\x93\x79\xb0\xf5\xe3\xf1\x77\x13\x59\xd8\xb4\x3c\x15\x42\xbe\x74\xc3\xd0\xdd\x25\xee\x72\x79\xd5\xad\xfa\xdc\xfc\x9f\x51\x8a\x6a\xb7\xb8\x0c\x83\xed\x06\x66\x3b\x23\xb0\x55\xd5\x76\xb2\x7d\x16\x71\x07\xd1\xc6\x30\xfe\x7b\x92\xcd\xd6\x3e\x4a\xf0\xc0\xec\xad\xe5\x0e\x82\x2a\xc4\x47\x08\x5e\xb8\x99\xdd\x4e\x71\x34\x5f\x35\x06\x61\x7a\x65\xd4\xca\xb9\x27\x15\x7f\x07\x00\x00\xff\xff\xb2\xb3\xb9\x00\x09\x0f\x00\x00")

func _1489143364_initial_schemaUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1489143364_initial_schemaUpSql,
		"1489143364_initial_schema.up.sql",
	)
}

func _1489143364_initial_schemaUpSql() (*asset, error) {
	bytes, err := _1489143364_initial_schemaUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1489143364_initial_schema.up.sql", size: 3849, mode: os.FileMode(0644), modTime: time.Unix(1611315902, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x1a, 0xa5, 0x48, 0xab, 0xd2, 0x53, 0x1, 0xbc, 0x7c, 0xa4, 0x7e, 0x16, 0x26, 0x7d, 0xd5, 0x8e, 0x50, 0x5, 0xfb, 0x19, 0x24, 0x8d, 0x57, 0x6e, 0x5a, 0x2e, 0x85, 0xa7, 0x24, 0x52, 0x4c, 0xa1}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
	"1489143364_initial_schema.down.sql": _1489143364_initial_schemaDownSql,
	"1489143364_initial_schema.up.sql":   _1489143364_initial_schemaUpSql,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	"1489143364_initial_schema.down.sql": {_1489143364_initial_schemaDownSql, map[string]*bintree{}},
	"1489143364_initial_schema.up.sql":   {_1489143364_initial_schemaUpSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
