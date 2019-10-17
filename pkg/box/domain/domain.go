package domain

import "os"

type Path string

func (p Path) ToString() string {
	return string(p)
}

type File struct {
	Path    Path
	Mode    os.FileMode
	Content []byte
}

type Action int

const (
	Write Action = iota
	Delete
)
