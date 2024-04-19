package gridfs

import (
	"database/sql"
	"errors"
	"io"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type File interface {
	fs.File
	fs.FileInfo
	ID() int64
	MD5() string
	ContentType() string
	ContentLength() string
	Disposition() string
}

// file 文件表
type file struct {
	id        int64
	name      string
	size      int64
	checksum  string
	createdAt time.Time

	db     *sql.DB // 数据库连接
	serial int     // 即将读取的文件分片序号
	buff   []byte  // 文件分片缓存
	eof    bool    // 是否读完了
	count  int
}

func (f *file) Stat() (fs.FileInfo, error) { return f, nil }
func (f *file) Close() error               { return nil }
func (f *file) Name() string               { return f.name }
func (f *file) Size() int64                { return f.size }
func (f *file) Mode() fs.FileMode          { return os.ModePerm }
func (f *file) ModTime() time.Time         { return f.createdAt }
func (f *file) IsDir() bool                { return false }
func (f *file) Sys() any                   { return nil }
func (f *file) ID() int64                  { return f.id }
func (f *file) MD5() string                { return f.checksum }

func (f *file) ContentType() string {
	ct := mime.TypeByExtension(filepath.Ext(f.name))
	// 当 Content-Type 是 text/html 时，就算设置了 Content-Length 长度，
	// 浏览器下载的时候也不会正常显示进度条，改一下 Content-Type 即可。
	if ct == "" || strings.HasPrefix(ct, "text/html") {
		ct = "application/octet-stream"
	}

	return ct
}

func (f *file) ContentLength() string {
	return strconv.FormatInt(f.size, 10)
}

func (f *file) Disposition() string {
	return mime.FormatMediaType("attachment", map[string]string{"filename": f.name})
}

func (f *file) Read(b []byte) (int, error) {
	if f.db == nil {
		return 0, io.ErrUnexpectedEOF
	}
	if f.eof {
		return 0, io.EOF
	}

	var err error
	var n int
	bsz := len(b)
	for !f.eof && bsz > n {
		if len(f.buff) == 0 {
			if err = f.nextSerial(); err != nil {
				break
			}
		}

		i := copy(b[n:], f.buff)
		n += i
		f.buff = f.buff[i:]
	}
	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

func (f *file) nextSerial() error {
	rawSQL := "SELECT data FROM gridfs_chunk WHERE file_id = $1 AND serial = $2"

	var chk chunk
	if err := f.db.QueryRow(rawSQL, f.id, f.serial).Scan(&chk.data); err != nil {
		f.eof = true
		if errors.Is(err, sql.ErrNoRows) {
			return io.EOF
		}
		return err
	}

	f.count += len(chk.data)

	f.serial++ // 序号增加
	f.buff = chk.data

	return nil
}

type chunk struct {
	id     int64  // 表 ID
	fileID int64  // 所属文件 ID
	data   []byte // 文件分片
	serial int    // 文件分片序号 0-n
}
