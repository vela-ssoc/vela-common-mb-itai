package gridfs

import (
	"io/fs"
	"os"
	"time"
)

type diskFile struct {
	disk *os.File
	db   File
}

func (df *diskFile) Stat() (fs.FileInfo, error) {
	return df.db.Stat()
}

func (df *diskFile) Read(b []byte) (int, error) {
	return df.disk.Read(b)
}

func (df *diskFile) Close() error {
	_ = df.db.Close()
	return df.disk.Close()
}

func (df *diskFile) Name() string {
	return df.db.Name()
}

func (df *diskFile) Size() int64 {
	return df.db.Size()
}

func (df *diskFile) Mode() fs.FileMode {
	return df.db.Mode()
}

func (df *diskFile) ModTime() time.Time {
	return df.db.ModTime()
}

func (df *diskFile) IsDir() bool {
	return df.db.IsDir()
}

func (df *diskFile) Sys() any {
	return df.db.Sys()
}

func (df *diskFile) ID() int64 {
	return df.db.ID()
}

func (df *diskFile) MD5() string {
	return df.db.MD5()
}

func (df *diskFile) ContentType() string {
	return df.db.ContentType()
}

func (df *diskFile) ContentLength() string {
	return df.db.ContentLength()
}

func (df *diskFile) Disposition() string {
	return df.db.Disposition()
}
