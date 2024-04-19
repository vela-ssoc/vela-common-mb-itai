package gridfs

import (
	"crypto/md5"
	"database/sql"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func NewCache(db *sql.DB, path string) FS {
	cfs := NewFS(db)
	intPID := os.Getpid()
	pid := strconv.FormatInt(int64(intPID), 10)
	if path == "" {
		path = os.TempDir()
	}
	_ = os.MkdirAll(path, os.ModePerm)

	return &cacheFS{
		dbfs:   cfs,
		path:   path,
		pid:    pid,
		caches: make(map[string]*copiedReader, 32),
	}
}

type cacheFS struct {
	dbfs   FS
	path   string
	pid    string
	mutex  sync.Mutex
	caches map[string]*copiedReader
}

func (cf *cacheFS) Open(name string) (fs.File, error) {
	fl, err := cf.dbfs.Open(name)
	if err != nil {
		return nil, err
	}

	f, ok := fl.(File)
	if !ok {
		_ = fl.Close()
		return nil, os.ErrNotExist
	}

	of, err := cf.openFile(f)
	if err != nil {
		return nil, err
	}

	return of, nil
}

func (cf *cacheFS) OpenID(id int64) (File, error) {
	if f, err := cf.dbfs.OpenID(id); err != nil {
		return nil, err
	} else {
		return cf.openFile(f)
	}
}

func (cf *cacheFS) Remove(id int64) error {
	return cf.dbfs.Remove(id)
}

func (cf *cacheFS) Write(r io.Reader, name string) (File, error) {
	return cf.dbfs.Write(r, name)
}

func (cf *cacheFS) openFile(dbFile File) (File, error) {
	localPath := cf.localDiskPath(dbFile)
	if disk, err := os.Open(localPath); err == nil { // 打开成功说明文件已经缓存完毕了。
		f := &diskFile{disk: disk, db: dbFile}
		return f, nil
	}

	tempPath := localPath + ".caching" // 缓存文件临时名字
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	teeFile := cf.caches[tempPath]
	if teeFile != nil { // 已经创建了缓存任务，并且还在运行中，当前这个下载任务就通过数据库下载。
		return dbFile, nil
	}
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return nil, err
	}

	md5hash := md5.New()
	writer := io.MultiWriter(md5hash, tempFile)
	tee := io.TeeReader(dbFile, writer)
	copiedTask := &copiedReader{
		cfs:      cf,
		source:   dbFile,
		target:   tempFile,
		md5hash:  md5hash,
		tee:      tee,
		tempPath: tempPath,
		donePath: localPath,
	}
	cf.caches[tempPath] = copiedTask

	return copiedTask, nil
}

// cachedName 缓存到本地的文件名，特定的规则防止冲突。
// 防止多个 manager/broker 程序部署在同一台机器上，用 pid 作为区分。
func (cf *cacheFS) localDiskPath(f File) string {
	// {pid}-{id}-{hash}-name
	id, hash, name := f.ID(), f.MD5(), f.Name()
	fid := strconv.FormatInt(id, 10)
	str := cf.pid + "-" + fid + "-" + hash + "-" + name

	return filepath.Join(cf.path, str)
}

func (cf *cacheFS) removeCopiedTask(tempPath string) {
	cf.mutex.Lock()
	delete(cf.caches, tempPath)
	cf.mutex.Unlock()
}
