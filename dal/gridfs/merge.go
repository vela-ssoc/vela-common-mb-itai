package gridfs

import (
	"bytes"
	"io"
	"io/fs"
	"strconv"
	"time"
)

func Merge(f File, d []byte) File {
	if len(d) == 0 {
		return f
	}

	br := bytes.NewReader(d)
	multi := io.MultiReader(f, br)

	return &mergeFile{
		file:  f,
		read:  multi,
		size:  f.Size() + br.Size(),
		mtime: time.Now(),
	}
}

type mergeFile struct {
	file  File
	read  io.Reader
	size  int64
	mtime time.Time
}

func (m *mergeFile) Stat() (fs.FileInfo, error) {
	return m, nil
}

func (m *mergeFile) Read(p []byte) (int, error) {
	return m.read.Read(p)
}

func (m *mergeFile) Close() error {
	return m.file.Close()
}

func (m *mergeFile) Name() string {
	return m.file.Name()
}

func (m *mergeFile) Size() int64 {
	return m.size
}

func (m *mergeFile) Mode() fs.FileMode {
	return m.file.Mode()
}

func (m *mergeFile) ModTime() time.Time {
	return m.mtime
}

func (m *mergeFile) IsDir() bool {
	return false
}

func (m *mergeFile) Sys() any {
	return nil
}

func (m *mergeFile) ID() int64 {
	return m.file.ID()
}

func (m *mergeFile) MD5() string {
	return ""
}

func (m *mergeFile) ContentType() string {
	return m.file.ContentType()
}

func (m *mergeFile) ContentLength() string {
	return strconv.FormatInt(m.size, 10)
}

func (m *mergeFile) Disposition() string {
	return m.file.Disposition()
}
