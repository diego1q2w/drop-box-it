package domain

import "os"

type Path string

func (p Path) ToString() string {
	return string(p)
}
func (p Path) ToBytes() []byte {
	return []byte(p)
}

type File struct {
	Path Path
	Mode os.FileMode
}

type FileStatus int

const (
	Created FileStatus = iota
	Updated
	Deleted
)
