package gridfs

import (
	"encoding/hex"
	"hash"
	"io"
	"io/fs"
	"os"
	"sync/atomic"
	"time"
)

type copiedReader struct {
	cfs      *cacheFS
	source   File
	target   *os.File
	md5hash  hash.Hash
	tee      io.Reader
	tempPath string
	donePath string
	closed   atomic.Bool
}

func (cr *copiedReader) Stat() (fs.FileInfo, error) {
	return cr.source.Stat()
}

func (cr *copiedReader) Read(b []byte) (int, error) {
	return cr.tee.Read(b)
}

func (cr *copiedReader) Name() string {
	return cr.source.Name()
}

func (cr *copiedReader) Size() int64 {
	return cr.source.Size()
}

func (cr *copiedReader) Mode() fs.FileMode {
	return cr.source.Mode()
}

func (cr *copiedReader) ModTime() time.Time {
	return cr.source.ModTime()
}

func (cr *copiedReader) IsDir() bool {
	return cr.source.IsDir()
}

func (cr *copiedReader) Sys() any {
	return cr.source.Sys()
}

func (cr *copiedReader) ID() int64 {
	return cr.source.ID()
}

func (cr *copiedReader) MD5() string {
	return cr.source.MD5()
}

func (cr *copiedReader) ContentType() string {
	return cr.source.ContentType()
}

func (cr *copiedReader) ContentLength() string {
	return cr.source.ContentLength()
}

func (cr *copiedReader) Disposition() string {
	return cr.source.Disposition()
}

func (cr *copiedReader) Close() error {
	if !cr.closed.CompareAndSwap(false, true) {
		return nil
	}

	// 关闭目的文件和源文件
	_ = cr.source.Close()
	_ = cr.target.Close()

	// 如果 md5 一致代表文件缓存没有出错
	var err error
	matched := cr.matchedMD5()
	if matched {
		err = os.Rename(cr.tempPath, cr.donePath)
	}
	if !matched || err != nil {
		err = os.Remove(cr.tempPath)
	}
	cr.cfs.removeCopiedTask(cr.tempPath) // 删除缓存任务表。

	return err
}

func (cr *copiedReader) matchedMD5() bool {
	smd5 := cr.source.MD5() // 源文件 md5
	b := cr.md5hash.Sum(nil)
	tmd5 := hex.EncodeToString(b) // 目的文件的 md5

	return smd5 == tmd5
}
