package domain

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Path string

func (p Path) ToString() string {
	return string(p)
}

func (p Path) RemoveBasePath(base Path) Path {
	removed := strings.ReplaceAll(p.ToString(), fmt.Sprintf("%s/", base), "")
	return Path(removed)
}

func (p Path) WithRoot(root Path) Path {
	return Path(filepath.Join(root.ToString(), p.ToString()))
}

func (p Path) ToBytes() []byte {
	return []byte(p)
}

type Content []byte

type HashContent [32]byte

type File struct {
	Path    Path
	Mode    os.FileMode
	Content HashContent
}

type FileStatus int

const (
	Created FileStatus = iota
	Updated
	Deleted
	Synced
)

func (c Content) Hash() HashContent {
	return sha256.Sum256(c)
}

func (f File) Equal(file File) bool {
	return f.Content == file.Content
}
