package directory

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"

	"model/library"
	"utils/fileserver"
)

var (
	ErrorNotDirectory = errors.New("not a directory")
)

type Library struct {
	id                  string
	path                string
	supportedExtensions map[string]struct{}
	fileserver          *fileserver.Instance
}

func New(path string, supportedExtensions []string) (*Library, error) {
	s, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !s.IsDir() {
		return nil, ErrorNotDirectory
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	exts := make(map[string]struct{})
	for _, ext := range supportedExtensions {
		exts[ext] = struct{}{}
	}
	lib := &Library{
		id:                  id.String(),
		path:                path,
		supportedExtensions: exts,
	}
	if lib.fileserver, err = fileserver.New(path); err != nil {
		return nil, err
	}
	return lib, nil
}

func (l *Library) ID() string {
	return l.id
}

func (l *Library) List() ([]library.Media, error) {
	return l.list(l.path, 0, 4)
}

func (l *Library) list(dir string, depth int, maxDepth int) ([]library.Media, error) {
	if depth >= maxDepth {
		return nil, nil
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var media []library.Media
	for _, file := range files {
		if file.IsDir() {
			recMedia, err := l.list(dir+"/"+file.Name(), depth+1, maxDepth)
			if err != nil {
				return nil, err
			}
			media = append(media, recMedia...)
			continue
		}
		ext := path.Ext(file.Name())
		if _, ok := l.supportedExtensions[ext]; !ok {
			continue
		}
		fileName := dir + "/" + file.Name()
		uri := strings.Replace(fileName, l.path, l.fileserver.ServingURL, 1)
		f, err := new(uri)
		if err != nil {
			return nil, err
		}
		media = append(media, f)
	}
	return media, nil
}
