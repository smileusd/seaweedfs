package filer2

import (
	"errors"
	"os"
	"time"
	"path/filepath"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"strings"
)

type FullPath string

func NewFullPath(dir, name string) FullPath {
	if strings.HasSuffix(dir, "/") {
		return FullPath(dir + name)
	}
	return FullPath(dir + "/" + name)
}

func (fp FullPath) DirAndName() (string, string) {
	dir, name := filepath.Split(string(fp))
	if dir == "/" {
		return dir, name
	}
	if len(dir) < 1 {
		return "/", ""
	}
	return dir[:len(dir)-1], name
}

func (fp FullPath) Name() (string) {
	_, name := filepath.Split(string(fp))
	return name
}

type Attr struct {
	Mtime  time.Time   // time of last modification
	Crtime time.Time   // time of creation (OS X only)
	Mode   os.FileMode // file mode
	Uid    uint32      // owner uid
	Gid    uint32      // group gid
}

func (attr Attr) IsDirectory() (bool) {
	return attr.Mode&os.ModeDir > 0
}

type Entry struct {
	FullPath

	Attr

	// the following is for files
	Chunks []*filer_pb.FileChunk `json:"chunks,omitempty"`
}

func (entry Entry) Size() uint64 {
	return TotalSize(entry.Chunks)
}

func (entry Entry) Timestamp() time.Time {
	if entry.IsDirectory() {
		return entry.Crtime
	} else {
		return entry.Mtime
	}
}

type AbstractFiler interface {
	CreateEntry(*Entry) (error)
	AppendFileChunk(FullPath, []*filer_pb.FileChunk) (err error)
	FindEntry(FullPath) (found bool, fileEntry *Entry, err error)
	DeleteEntry(FullPath) (fileEntry *Entry, err error)

	ListDirectoryEntries(dirPath FullPath) ([]*Entry, error)
	UpdateEntry(*Entry) (error)
}

var ErrNotFound = errors.New("filer: no entry is found in filer store")

type FilerStore interface {
	InsertEntry(*Entry) (error)
	UpdateEntry(*Entry) (err error)
	FindEntry(FullPath) (found bool, entry *Entry, err error)
	DeleteEntry(FullPath) (fileEntry *Entry, err error)

	ListDirectoryEntries(dirPath FullPath, startFileName string, inclusive bool, limit int) ([]*Entry, error)
}