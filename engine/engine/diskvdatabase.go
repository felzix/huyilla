package engine

import (
	"github.com/gogo/protobuf/proto"
	"github.com/peterbourgon/diskv/v3"
	"regexp"
	"strings"
)

type DiskVDatabase struct {
	diskV *diskv.Diskv
}

func NewDisKVDatabase(saveDir, tempDir string, cacheSize uint64) *DiskVDatabase {
	return &DiskVDatabase{
		diskV: diskv.New(diskv.Options{
			BasePath:          saveDir,
			TempDir:           tempDir,
			AdvancedTransform: filesystemTransform,
			InverseTransform:  filesystemInverseTransform,
			CacheSizeMax:      cacheSize,
		})}
}

func (db *DiskVDatabase) Get(key string, thing proto.Unmarshaler) error {
	if blob, err := db.diskV.Read(key); err == nil {
		if err := thing.Unmarshal(blob); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (db *DiskVDatabase) GetByPrefix(prefix string) <-chan string {
	return db.diskV.KeysPrefix(prefix, nil)
}

func (db *DiskVDatabase) Set(key string, thing proto.Marshaler) error {
	if blob, err := thing.Marshal(); err == nil {
		if err := db.diskV.Write(key, blob); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (db *DiskVDatabase) Has(key string) bool {
	return db.diskV.Has(key)
}

func (db *DiskVDatabase) End(key string) error {
	return db.diskV.Erase(key)
}

func (db *DiskVDatabase) EndAll() error {
	return db.diskV.EraseAll()
}

func filesystemTransform(key string) *diskv.PathKey {
	path := strings.Split(key, ".")
	last := len(path) - 1
	return &diskv.PathKey{
		Path:     path[:last],
		FileName: path[last],
	}
}

func filesystemInverseTransform(pathKey *diskv.PathKey) (key string) {
	return strings.Join(pathKey.Path, "/") + pathKey.FileName
}

var regexFileNotFound = regexp.MustCompile("no such file or directory")

func fileIsNotFound(err error) bool {
	return regexFileNotFound.MatchString(err.Error())
}
