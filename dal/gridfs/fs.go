package gridfs

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"io"
	"io/fs"
	"time"
)

type FS interface {
	fs.FS
	OpenID(id int64) (File, error)
	Remove(id int64) error
	Write(io.Reader, string) (File, error)
}

func NewFS(db *sql.DB) FS {
	return &gridFS{db: db, burst: 60 * 1024}
}

type gridFS struct {
	db    *sql.DB // 数据库连接
	burst int     // 60K
}

// Open 通过文件名打开文件
func (f *gridFS) Open(name string) (fs.File, error) {
	rawSQL := "SELECT id, size, name, checksum, created_at FROM gridfs_file WHERE name = $1 ORDER BY id DESC LIMIT 1"
	fl := &file{db: f.db}
	if err := f.db.QueryRow(rawSQL, name).
		Scan(&fl.id, &fl.size, &fl.name, &fl.checksum, &fl.createdAt); err != nil {
		return nil, err
	}

	return fl, nil
}

func (f *gridFS) OpenID(id int64) (File, error) {
	rawSQL := "SELECT id, size, name, checksum, created_at FROM gridfs_file WHERE id = $1"
	fl := &file{db: f.db}
	if err := f.db.QueryRow(rawSQL, id).
		Scan(&fl.id, &fl.size, &fl.name, &fl.checksum, &fl.createdAt); err != nil {
		return nil, err
	}

	return fl, nil
}

func (f *gridFS) Remove(id int64) error {
	fileSQL := "DELETE FROM gridfs_file WHERE id = $1"
	chkSQL := "DELETE FROM gridfs_chunk WHERE file_id = $1"

	tx, err := f.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer tx.Rollback()

	var over bool
	defer func() {
		if err != nil || !over {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.Exec(fileSQL, id); err == nil {
		if _, err = tx.Exec(chkSQL, id); err == nil {
			err = tx.Commit()
		}
	}
	over = true

	return err
}

func (f *gridFS) Write(r io.Reader, name string) (File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	// 开启事务
	opt := &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	tx, err := f.db.BeginTx(ctx, opt)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer tx.Rollback()

	burst := f.burst
	createdAt := time.Now()
	insertFile := "INSERT INTO gridfs_file(name, checksum, created_at) VALUES ($1, $2, $3) RETURNING id"
	var fileID int64
	// see: https://github.com/lib/pq/issues/24
	if err = tx.QueryRowContext(ctx, insertFile, name, "", createdAt).Scan(&fileID); err != nil {
		return nil, err
	}

	insertPart := "INSERT INTO gridfs_chunk (file_id, serial, data) VALUES ($1, $2, $3)"
	buf := make([]byte, burst)

	digest := md5.New()
	tr := io.TeeReader(r, digest)

	var n, serial int
	var filesize int64
	for {
		if n, err = tr.Read(buf); err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		if _, err = tx.ExecContext(ctx, insertPart, fileID, serial, buf[:n]); err != nil {
			break
		}
		serial++
		filesize += int64(n)
	}

	if err == nil {
		sum := hex.EncodeToString(digest.Sum(nil))
		updateFile := "UPDATE gridfs_file SET size = $1, checksum = $2 WHERE id = $3"
		if _, err = tx.ExecContext(ctx, updateFile, filesize, sum, fileID); err == nil {
			if err = tx.Commit(); err == nil {
				fl := &file{
					id:        fileID,
					name:      name,
					size:      filesize,
					checksum:  sum,
					createdAt: createdAt,
					db:        f.db,
				}

				return fl, nil
			}
		}
	}

	return nil, err
}
