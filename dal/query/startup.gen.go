// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
)

func newStartup(db *gorm.DB, opts ...gen.DOOption) startup {
	_startup := startup{}

	_startup.startupDo.UseDB(db, opts...)
	_startup.startupDo.UseModel(&model.Startup{})

	tableName := _startup.startupDo.TableName()
	_startup.ALL = field.NewAsterisk(tableName)
	_startup.ID = field.NewInt64(tableName, "id")
	_startup.Node = field.NewField(tableName, "node")
	_startup.Logger = field.NewField(tableName, "logger")
	_startup.Console = field.NewField(tableName, "console")
	_startup.Extends = field.NewField(tableName, "extends")
	_startup.Failed = field.NewBool(tableName, "failed")
	_startup.Reason = field.NewString(tableName, "reason")
	_startup.CreatedAt = field.NewTime(tableName, "created_at")
	_startup.UpdatedAt = field.NewTime(tableName, "updated_at")

	_startup.fillFieldMap()

	return _startup
}

type startup struct {
	startupDo startupDo

	ALL       field.Asterisk
	ID        field.Int64
	Node      field.Field
	Logger    field.Field
	Console   field.Field
	Extends   field.Field
	Failed    field.Bool
	Reason    field.String
	CreatedAt field.Time
	UpdatedAt field.Time

	fieldMap map[string]field.Expr
}

func (s startup) Table(newTableName string) *startup {
	s.startupDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s startup) As(alias string) *startup {
	s.startupDo.DO = *(s.startupDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *startup) updateTableName(table string) *startup {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewInt64(table, "id")
	s.Node = field.NewField(table, "node")
	s.Logger = field.NewField(table, "logger")
	s.Console = field.NewField(table, "console")
	s.Extends = field.NewField(table, "extends")
	s.Failed = field.NewBool(table, "failed")
	s.Reason = field.NewString(table, "reason")
	s.CreatedAt = field.NewTime(table, "created_at")
	s.UpdatedAt = field.NewTime(table, "updated_at")

	s.fillFieldMap()

	return s
}

func (s *startup) WithContext(ctx context.Context) *startupDo { return s.startupDo.WithContext(ctx) }

func (s startup) TableName() string { return s.startupDo.TableName() }

func (s startup) Alias() string { return s.startupDo.Alias() }

func (s startup) Columns(cols ...field.Expr) gen.Columns { return s.startupDo.Columns(cols...) }

func (s *startup) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *startup) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 9)
	s.fieldMap["id"] = s.ID
	s.fieldMap["node"] = s.Node
	s.fieldMap["logger"] = s.Logger
	s.fieldMap["console"] = s.Console
	s.fieldMap["extends"] = s.Extends
	s.fieldMap["failed"] = s.Failed
	s.fieldMap["reason"] = s.Reason
	s.fieldMap["created_at"] = s.CreatedAt
	s.fieldMap["updated_at"] = s.UpdatedAt
}

func (s startup) clone(db *gorm.DB) startup {
	s.startupDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s startup) replaceDB(db *gorm.DB) startup {
	s.startupDo.ReplaceDB(db)
	return s
}

type startupDo struct{ gen.DO }

func (s startupDo) Debug() *startupDo {
	return s.withDO(s.DO.Debug())
}

func (s startupDo) WithContext(ctx context.Context) *startupDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s startupDo) ReadDB() *startupDo {
	return s.Clauses(dbresolver.Read)
}

func (s startupDo) WriteDB() *startupDo {
	return s.Clauses(dbresolver.Write)
}

func (s startupDo) Session(config *gorm.Session) *startupDo {
	return s.withDO(s.DO.Session(config))
}

func (s startupDo) Clauses(conds ...clause.Expression) *startupDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s startupDo) Returning(value interface{}, columns ...string) *startupDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s startupDo) Not(conds ...gen.Condition) *startupDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s startupDo) Or(conds ...gen.Condition) *startupDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s startupDo) Select(conds ...field.Expr) *startupDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s startupDo) Where(conds ...gen.Condition) *startupDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s startupDo) Order(conds ...field.Expr) *startupDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s startupDo) Distinct(cols ...field.Expr) *startupDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s startupDo) Omit(cols ...field.Expr) *startupDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s startupDo) Join(table schema.Tabler, on ...field.Expr) *startupDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s startupDo) LeftJoin(table schema.Tabler, on ...field.Expr) *startupDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s startupDo) RightJoin(table schema.Tabler, on ...field.Expr) *startupDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s startupDo) Group(cols ...field.Expr) *startupDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s startupDo) Having(conds ...gen.Condition) *startupDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s startupDo) Limit(limit int) *startupDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s startupDo) Offset(offset int) *startupDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s startupDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *startupDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s startupDo) Unscoped() *startupDo {
	return s.withDO(s.DO.Unscoped())
}

func (s startupDo) Create(values ...*model.Startup) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s startupDo) CreateInBatches(values []*model.Startup, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s startupDo) Save(values ...*model.Startup) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s startupDo) First() (*model.Startup, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Startup), nil
	}
}

func (s startupDo) Take() (*model.Startup, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Startup), nil
	}
}

func (s startupDo) Last() (*model.Startup, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Startup), nil
	}
}

func (s startupDo) Find() ([]*model.Startup, error) {
	result, err := s.DO.Find()
	return result.([]*model.Startup), err
}

func (s startupDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Startup, err error) {
	buf := make([]*model.Startup, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s startupDo) FindInBatches(result *[]*model.Startup, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s startupDo) Attrs(attrs ...field.AssignExpr) *startupDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s startupDo) Assign(attrs ...field.AssignExpr) *startupDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s startupDo) Joins(fields ...field.RelationField) *startupDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s startupDo) Preload(fields ...field.RelationField) *startupDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s startupDo) FirstOrInit() (*model.Startup, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Startup), nil
	}
}

func (s startupDo) FirstOrCreate() (*model.Startup, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Startup), nil
	}
}

func (s startupDo) FindByPage(offset int, limit int) (result []*model.Startup, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s startupDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s startupDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s startupDo) Delete(models ...*model.Startup) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *startupDo) withDO(do gen.Dao) *startupDo {
	s.DO = *do.(*gen.DO)
	return s
}
