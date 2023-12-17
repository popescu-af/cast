package directory

import (
	"crypto/sha256"
	"fmt"
)

type File struct {
	id  string
	uri string
}

func new(uri string) (*File, error) {
	sum := sha256.Sum256([]byte(uri))
	return &File{
		id:  fmt.Sprintf("%x", sum[:16]),
		uri: uri,
	}, nil
}

func (f *File) ID() string {
	return f.id
}

func (f *File) URI() string {
	return f.uri
}
