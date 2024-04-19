package gridfs

import (
	"database/sql"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// NewCDN 带本地缓存的。
//
// Deprecated: 该缓存在线上运行时会出现文件无法下载的 BUG，
// 已废弃请勿使用，新的接口为 [NewCache]。
func NewCDN(db *sql.DB, dir string, min int64) FS {
	gfs := NewFS(db)
	if dir == "" {
		dir = os.TempDir()
	}
	files := make(map[int64]*cdnFile, 64)

	return &cdnFS{gfs: gfs, dir: dir, min: min, files: files}
}

// cdnFS 带文件缓存的 FS 管理器
type cdnFS struct {
	gfs   FS                 // 底层文件管理器
	dir   string             // CDN 文件缓存的目录
	min   int64              // filesize 小于或等于 min 时不会经过 CDN 缓存
	mutex sync.Mutex         // files lock
	files map[int64]*cdnFile // cdnFS 缓存文件映射
}

// Open implement fs.FS
func (cn *cdnFS) Open(name string) (fs.File, error) {
	fl, err := cn.gfs.Open(name)
	if err != nil {
		return nil, err
	}
	f, ok := fl.(File)
	if !ok {
		_ = fl.Close()
		return nil, os.ErrPermission
	}

	return cn.openFile(f)
}

// OpenID 通过 ID 打开文件
func (cn *cdnFS) OpenID(id int64) (File, error) {
	// 先查询数据库
	fl, err := cn.gfs.OpenID(id)
	if err != nil {
		return nil, err
	}

	return cn.openFile(fl)
}

// Remove 删除文件
func (cn *cdnFS) Remove(id int64) error {
	cn.removeID(id)
	return cn.gfs.Remove(id)
}

// Write 保存文件
func (cn *cdnFS) Write(r io.Reader, name string) (File, error) {
	return cn.gfs.Write(r, name)
}

func (cn *cdnFS) fromCDN(mfl File) (File, error) {
	cn.mutex.Lock()
	defer cn.mutex.Unlock()

	fileID := mfl.ID()
	name := mfl.Name()
	cf, ok := cn.files[fileID]
	if ok {
		// 如果 CDN 本地文件已经缓存完毕，就直接读取本地磁盘的文件，
		// 不再从数据库下载。
		if cf.done.Load() {
			return cf.Open()
		}

		// 如果文件已经加入了 CDN 缓存，但是还未缓存完毕，那么此次
		// 继续走数据库下载。
		return mfl, nil
	}

	// 如果 CDN 还未创建，就创建 CDN 缓存任务
	filename := strconv.FormatInt(fileID, 10) + "-" + name
	disk := filepath.Join(cn.dir, filename)
	dfl, err := os.Create(disk)
	if err != nil {
		return mfl, nil
	}

	// 客户端下载的同时进行磁盘缓存
	cfo := &cdnFile{file: mfl, disk: disk, tmp: dfl}
	cn.files[fileID] = cfo

	ch := cn.newCaching(mfl, cfo, dfl)

	return ch, nil
}

func (cn *cdnFS) openFile(fl File) (File, error) {
	size := fl.Size()
	if size <= cn.min {
		return fl, nil
	}
	return cn.fromCDN(fl)
}

func (cn *cdnFS) newCaching(raw File, cf *cdnFile, disk *os.File) *cdnCaching {
	teeRead := io.TeeReader(raw, disk)

	return &cdnCaching{
		cdn:     cn,
		cdnFile: cf,
		rawFile: raw,
		teeRead: teeRead,
	}
}

func (cn *cdnFS) removeID(id int64) {
	cn.mutex.Lock()
	defer cn.mutex.Unlock()

	cf, ok := cn.files[id]
	if ok {
		delete(cn.files, id)
		_ = os.Remove(cf.disk)
	}
}

type cdnFile struct {
	file File
	disk string
	done atomic.Bool
	tmp  *os.File
}

func (cf *cdnFile) Open() (File, error) {
	open, err := os.Open(cf.disk)
	if err != nil {
		return nil, err
	}

	wf := &warpedFile{diskFile: open, rawFile: cf.file}

	return wf, nil
}

type warpedFile struct {
	diskFile *os.File
	rawFile  File
}

func (wf *warpedFile) Stat() (fs.FileInfo, error) { return wf.diskFile.Stat() }
func (wf *warpedFile) Read(p []byte) (int, error) { return wf.diskFile.Read(p) }
func (wf *warpedFile) Close() error               { return wf.diskFile.Close() }
func (wf *warpedFile) Name() string               { return wf.rawFile.Name() }
func (wf *warpedFile) Size() int64                { return wf.rawFile.Size() }
func (wf *warpedFile) Mode() fs.FileMode          { return wf.rawFile.Mode() }
func (wf *warpedFile) ModTime() time.Time         { return wf.rawFile.ModTime() }
func (wf *warpedFile) IsDir() bool                { return wf.rawFile.IsDir() }
func (wf *warpedFile) Sys() any                   { return wf.rawFile.Sys() }
func (wf *warpedFile) MD5() string                { return wf.rawFile.MD5() }
func (wf *warpedFile) ContentType() string        { return wf.rawFile.ContentType() }
func (wf *warpedFile) ContentLength() string      { return wf.rawFile.ContentLength() }
func (wf *warpedFile) Disposition() string        { return wf.rawFile.Disposition() }
func (wf *warpedFile) ID() int64                  { return wf.rawFile.ID() }

// Seek 实现 io.Seeker 用于支持断点续传
func (wf *warpedFile) Seek(offset int64, whence int) (int64, error) {
	return wf.diskFile.Seek(offset, whence)
}

type cdnCaching struct {
	cdn     *cdnFS
	cdnFile *cdnFile
	rawFile File
	teeRead io.Reader
}

func (cc *cdnCaching) Stat() (fs.FileInfo, error) { return cc.rawFile.Stat() }
func (cc *cdnCaching) Name() string               { return cc.rawFile.Name() }
func (cc *cdnCaching) Size() int64                { return cc.rawFile.Size() }
func (cc *cdnCaching) Mode() fs.FileMode          { return cc.rawFile.Mode() }
func (cc *cdnCaching) ModTime() time.Time         { return cc.rawFile.ModTime() }
func (cc *cdnCaching) IsDir() bool                { return cc.rawFile.IsDir() }
func (cc *cdnCaching) Sys() any                   { return cc.rawFile.Sys() }
func (cc *cdnCaching) MD5() string                { return cc.rawFile.MD5() }
func (cc *cdnCaching) ContentType() string        { return cc.rawFile.ContentType() }
func (cc *cdnCaching) ContentLength() string      { return cc.rawFile.ContentLength() }
func (cc *cdnCaching) Disposition() string        { return cc.rawFile.Disposition() }
func (cc *cdnCaching) ID() int64                  { return cc.rawFile.ID() }

func (cc *cdnCaching) Read(p []byte) (n int, err error) {
	if n, err = cc.teeRead.Read(p); err == io.EOF {
		// 如果 io.EOF 代表文件已经读取完毕，将 CDN 缓存任务状态设置为 done
		cc.cdnFile.done.Store(true)
		_ = cc.cdnFile.tmp.Close()
	}

	return
}

func (cc *cdnCaching) Close() error {
	if !cc.cdnFile.done.Load() {
		_ = cc.cdnFile.tmp.Close()
		cc.cdn.removeID(cc.rawFile.ID())
	}

	return cc.rawFile.Close()
}
