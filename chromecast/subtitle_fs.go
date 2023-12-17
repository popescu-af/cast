package chromecast

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"model/library"
	"utils/fileserver"

	"github.com/asticode/go-astisub"
)

type subtitleFS struct {
	fileserver *fileserver.Instance
	outputDir  string
}

func newSubtitleFS() (*subtitleFS, error) {
	tempDir, err := os.MkdirTemp("", "cast-*")
	if err != nil {
		return nil, err
	}
	fileserver, err := fileserver.New(tempDir)
	if err != nil {
		return nil, err
	}
	return &subtitleFS{
		fileserver: fileserver,
		outputDir:  tempDir,
	}, nil
}

func (s *subtitleFS) convertToVTT(subtitle library.Media) (library.Media, error) {
	file, err := os.CreateTemp(s.outputDir, "temp-subtitle-*.vtt")
	if err != nil {
		return nil, err
	}
	outputPath := file.Name()
	file.Close()
	subtitlePath, err := s.downloadSubtitle(subtitle.URI())
	if err != nil {
		return nil, err
	}
	sub, err := astisub.Open(astisub.Options{Filename: subtitlePath})
	if err != nil {
		return nil, err
	}
	err = sub.Write(outputPath)
	if err != nil {
		return nil, err
	}
	return &vttSubtitle{
		id:  subtitle.ID(),
		uri: strings.Replace(outputPath, s.outputDir, s.fileserver.ServingURL, 1),
	}, nil
}

func (s *subtitleFS) downloadSubtitle(uri string) (string, error) {
	file, err := os.CreateTemp(s.outputDir, "downloaded-subtitle-*"+filepath.Ext(uri))
	if err != nil {
		return "", err
	}
	defer file.Close()

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	_, err = file.Write(b)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

type vttSubtitle struct {
	id  string
	uri string
}

func (v *vttSubtitle) ID() string {
	return v.id
}

func (v *vttSubtitle) URI() string {
	return v.uri
}
