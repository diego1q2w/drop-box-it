package domain

import "os"

type Path string

func (p Path) ToString() string {
	return string(p)
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
