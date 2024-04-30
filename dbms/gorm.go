package dbms

import (
	"database/sql"
	"net/url"
	"time"

	_ "gitee.com/opengauss/openGauss-connector-go-pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Config struct {
	MaxOpenConn int               `json:"max_open_conn" yaml:"max_open_conn"`                          // 最大连接数
	MaxIdleConn int               `json:"max_idle_conn" yaml:"max_idle_conn"`                          // 最大空闲连接数
	MaxLifeTime time.Duration     `json:"max_life_time" yaml:"max_life_time"`                          // 连接最大存活时长
	MaxIdleTime time.Duration     `json:"max_idle_time" yaml:"max_idle_time"`                          // 空闲连接最大时长
	DSN         string            `json:"dsn"           yaml:"dsn"`                                    // 数据源
	User        string            `json:"user"          yaml:"user"   validate:"required_without=DSN"` // 数据库用户名
	Passwd      string            `json:"passwd"        yaml:"passwd" validate:"required_without=DSN"` // 密码
	Net         string            `json:"net"           yaml:"net"`                                    // 连接协议
	Addr        string            `json:"addr"          yaml:"addr"   validate:"required_without=DSN"` // 连接地址
	DBName      string            `json:"dbname"        yaml:"dbname" validate:"required_without=DSN"` // 库名
	Params      map[string]string `json:"params"        yaml:"params"`                                 // 参数
}

// FormatDSN 生成数据库连接
func (db Config) FormatDSN() string {
	if dsn := db.DSN; dsn != "" {
		return dsn
	}

	params := make(url.Values, 8)
	for k, v := range db.Params {
		params.Set(k, v)
	}
	u := &url.URL{
		Scheme:   "opengauss",
		User:     url.UserPassword(db.User, db.Passwd),
		Host:     db.Addr,
		Path:     db.DBName,
		RawQuery: params.Encode(),
	}

	return u.String()
}

// Open 连接数据库
func Open(cfg Config, lgi logger.Interface) (*gorm.DB, *sql.DB, error) {
	dsn := cfg.FormatDSN()

	gauss := postgres.Config{DriverName: "opengauss", DSN: dsn}
	db, err := gorm.Open(postgres.New(gauss), &gorm.Config{Logger: lgi})
	if err != nil {
		return nil, nil, err
	}
	sdb, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	// ----------[ 设置连接参数 ]----------
	sdb.SetMaxIdleConns(cfg.MaxIdleConn)
	sdb.SetMaxOpenConns(cfg.MaxOpenConn)
	sdb.SetConnMaxLifetime(cfg.MaxLifeTime)
	sdb.SetConnMaxIdleTime(cfg.MaxIdleTime)

	callbackConfig := &callbacks.Config{
		CreateClauses: []string{"INSERT", "VALUES", "ON CONFLICT"},
		UpdateClauses: []string{"UPDATE", "SET", "FROM", "WHERE"},
		DeleteClauses: []string{"DELETE", "FROM", "WHERE"},
	}
	if !gauss.WithoutReturning {
		callbackConfig.CreateClauses = append(callbackConfig.CreateClauses, "RETURNING")
		callbackConfig.UpdateClauses = append(callbackConfig.UpdateClauses, "RETURNING")
		callbackConfig.DeleteClauses = append(callbackConfig.DeleteClauses, "RETURNING")
	}
	_ = db.Callback().
		Create().
		Replace("gorm:create", Create(callbackConfig))
	// 雪花 ID 生成器
	sn := newSnow()
	_ = db.Callback().
		Create().
		Before("gorm:create").
		Register("generate_id", sn.autoID)
	db.ClauseBuilders["ON CONFLICT"] = onConflictFunc

	return db, sdb, nil
}

type rewriteCreate struct {
	crt func(*gorm.DB)
}

func (rc *rewriteCreate) rewrite(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	if db.Statement.Schema != nil {
		if len(db.Statement.Schema.FieldsWithDefaultDBValue) > 0 {
			if _, ok := db.Statement.Clauses["ON CONFLICT"]; ok {
				db.Statement.AddClause(clause.Returning{})
			}
		}
	}

	rc.crt(db)
}
