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

func newJobCode(db *gorm.DB, opts ...gen.DOOption) jobCode {
	_jobCode := jobCode{}

	_jobCode.jobCodeDo.UseDB(db, opts...)
	_jobCode.jobCodeDo.UseModel(&model.JobCode{})

	tableName := _jobCode.jobCodeDo.TableName()
	_jobCode.ALL = field.NewAsterisk(tableName)
	_jobCode.ID = field.NewInt64(tableName, "id")
	_jobCode.Name = field.NewString(tableName, "name")
	_jobCode.Icon = field.NewBytes(tableName, "icon")
	_jobCode.Chunk = field.NewBytes(tableName, "chunk")
	_jobCode.Desc = field.NewString(tableName, "desc")
	_jobCode.Hash = field.NewString(tableName, "hash")
	_jobCode.CreatedAt = field.NewTime(tableName, "created_at")
	_jobCode.UpdatedAt = field.NewTime(tableName, "updated_at")

	_jobCode.fillFieldMap()

	return _jobCode
}

type jobCode struct {
	jobCodeDo jobCodeDo

	ALL       field.Asterisk
	ID        field.Int64
	Name      field.String
	Icon      field.Bytes
	Chunk     field.Bytes
	Desc      field.String
	Hash      field.String
	CreatedAt field.Time
	UpdatedAt field.Time

	fieldMap map[string]field.Expr
}

func (j jobCode) Table(newTableName string) *jobCode {
	j.jobCodeDo.UseTable(newTableName)
	return j.updateTableName(newTableName)
}

func (j jobCode) As(alias string) *jobCode {
	j.jobCodeDo.DO = *(j.jobCodeDo.As(alias).(*gen.DO))
	return j.updateTableName(alias)
}

func (j *jobCode) updateTableName(table string) *jobCode {
	j.ALL = field.NewAsterisk(table)
	j.ID = field.NewInt64(table, "id")
	j.Name = field.NewString(table, "name")
	j.Icon = field.NewBytes(table, "icon")
	j.Chunk = field.NewBytes(table, "chunk")
	j.Desc = field.NewString(table, "desc")
	j.Hash = field.NewString(table, "hash")
	j.CreatedAt = field.NewTime(table, "created_at")
	j.UpdatedAt = field.NewTime(table, "updated_at")

	j.fillFieldMap()

	return j
}

func (j *jobCode) WithContext(ctx context.Context) *jobCodeDo { return j.jobCodeDo.WithContext(ctx) }

func (j jobCode) TableName() string { return j.jobCodeDo.TableName() }

func (j jobCode) Alias() string { return j.jobCodeDo.Alias() }

func (j jobCode) Columns(cols ...field.Expr) gen.Columns { return j.jobCodeDo.Columns(cols...) }

func (j *jobCode) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := j.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (j *jobCode) fillFieldMap() {
	j.fieldMap = make(map[string]field.Expr, 8)
	j.fieldMap["id"] = j.ID
	j.fieldMap["name"] = j.Name
	j.fieldMap["icon"] = j.Icon
	j.fieldMap["chunk"] = j.Chunk
	j.fieldMap["desc"] = j.Desc
	j.fieldMap["hash"] = j.Hash
	j.fieldMap["created_at"] = j.CreatedAt
	j.fieldMap["updated_at"] = j.UpdatedAt
}

func (j jobCode) clone(db *gorm.DB) jobCode {
	j.jobCodeDo.ReplaceConnPool(db.Statement.ConnPool)
	return j
}

func (j jobCode) replaceDB(db *gorm.DB) jobCode {
	j.jobCodeDo.ReplaceDB(db)
	return j
}

type jobCodeDo struct{ gen.DO }

func (j jobCodeDo) Debug() *jobCodeDo {
	return j.withDO(j.DO.Debug())
}

func (j jobCodeDo) WithContext(ctx context.Context) *jobCodeDo {
	return j.withDO(j.DO.WithContext(ctx))
}

func (j jobCodeDo) ReadDB() *jobCodeDo {
	return j.Clauses(dbresolver.Read)
}

func (j jobCodeDo) WriteDB() *jobCodeDo {
	return j.Clauses(dbresolver.Write)
}

func (j jobCodeDo) Session(config *gorm.Session) *jobCodeDo {
	return j.withDO(j.DO.Session(config))
}

func (j jobCodeDo) Clauses(conds ...clause.Expression) *jobCodeDo {
	return j.withDO(j.DO.Clauses(conds...))
}

func (j jobCodeDo) Returning(value interface{}, columns ...string) *jobCodeDo {
	return j.withDO(j.DO.Returning(value, columns...))
}

func (j jobCodeDo) Not(conds ...gen.Condition) *jobCodeDo {
	return j.withDO(j.DO.Not(conds...))
}

func (j jobCodeDo) Or(conds ...gen.Condition) *jobCodeDo {
	return j.withDO(j.DO.Or(conds...))
}

func (j jobCodeDo) Select(conds ...field.Expr) *jobCodeDo {
	return j.withDO(j.DO.Select(conds...))
}

func (j jobCodeDo) Where(conds ...gen.Condition) *jobCodeDo {
	return j.withDO(j.DO.Where(conds...))
}

func (j jobCodeDo) Order(conds ...field.Expr) *jobCodeDo {
	return j.withDO(j.DO.Order(conds...))
}

func (j jobCodeDo) Distinct(cols ...field.Expr) *jobCodeDo {
	return j.withDO(j.DO.Distinct(cols...))
}

func (j jobCodeDo) Omit(cols ...field.Expr) *jobCodeDo {
	return j.withDO(j.DO.Omit(cols...))
}

func (j jobCodeDo) Join(table schema.Tabler, on ...field.Expr) *jobCodeDo {
	return j.withDO(j.DO.Join(table, on...))
}

func (j jobCodeDo) LeftJoin(table schema.Tabler, on ...field.Expr) *jobCodeDo {
	return j.withDO(j.DO.LeftJoin(table, on...))
}

func (j jobCodeDo) RightJoin(table schema.Tabler, on ...field.Expr) *jobCodeDo {
	return j.withDO(j.DO.RightJoin(table, on...))
}

func (j jobCodeDo) Group(cols ...field.Expr) *jobCodeDo {
	return j.withDO(j.DO.Group(cols...))
}

func (j jobCodeDo) Having(conds ...gen.Condition) *jobCodeDo {
	return j.withDO(j.DO.Having(conds...))
}

func (j jobCodeDo) Limit(limit int) *jobCodeDo {
	return j.withDO(j.DO.Limit(limit))
}

func (j jobCodeDo) Offset(offset int) *jobCodeDo {
	return j.withDO(j.DO.Offset(offset))
}

func (j jobCodeDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *jobCodeDo {
	return j.withDO(j.DO.Scopes(funcs...))
}

func (j jobCodeDo) Unscoped() *jobCodeDo {
	return j.withDO(j.DO.Unscoped())
}

func (j jobCodeDo) Create(values ...*model.JobCode) error {
	if len(values) == 0 {
		return nil
	}
	return j.DO.Create(values)
}

func (j jobCodeDo) CreateInBatches(values []*model.JobCode, batchSize int) error {
	return j.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (j jobCodeDo) Save(values ...*model.JobCode) error {
	if len(values) == 0 {
		return nil
	}
	return j.DO.Save(values)
}

func (j jobCodeDo) First() (*model.JobCode, error) {
	if result, err := j.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.JobCode), nil
	}
}

func (j jobCodeDo) Take() (*model.JobCode, error) {
	if result, err := j.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.JobCode), nil
	}
}

func (j jobCodeDo) Last() (*model.JobCode, error) {
	if result, err := j.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.JobCode), nil
	}
}

func (j jobCodeDo) Find() ([]*model.JobCode, error) {
	result, err := j.DO.Find()
	return result.([]*model.JobCode), err
}

func (j jobCodeDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.JobCode, err error) {
	buf := make([]*model.JobCode, 0, batchSize)
	err = j.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (j jobCodeDo) FindInBatches(result *[]*model.JobCode, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return j.DO.FindInBatches(result, batchSize, fc)
}

func (j jobCodeDo) Attrs(attrs ...field.AssignExpr) *jobCodeDo {
	return j.withDO(j.DO.Attrs(attrs...))
}

func (j jobCodeDo) Assign(attrs ...field.AssignExpr) *jobCodeDo {
	return j.withDO(j.DO.Assign(attrs...))
}

func (j jobCodeDo) Joins(fields ...field.RelationField) *jobCodeDo {
	for _, _f := range fields {
		j = *j.withDO(j.DO.Joins(_f))
	}
	return &j
}

func (j jobCodeDo) Preload(fields ...field.RelationField) *jobCodeDo {
	for _, _f := range fields {
		j = *j.withDO(j.DO.Preload(_f))
	}
	return &j
}

func (j jobCodeDo) FirstOrInit() (*model.JobCode, error) {
	if result, err := j.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.JobCode), nil
	}
}

func (j jobCodeDo) FirstOrCreate() (*model.JobCode, error) {
	if result, err := j.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.JobCode), nil
	}
}

func (j jobCodeDo) FindByPage(offset int, limit int) (result []*model.JobCode, count int64, err error) {
	result, err = j.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = j.Offset(-1).Limit(-1).Count()
	return
}

func (j jobCodeDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = j.Count()
	if err != nil {
		return
	}

	err = j.Offset(offset).Limit(limit).Scan(result)
	return
}

func (j jobCodeDo) Scan(result interface{}) (err error) {
	return j.DO.Scan(result)
}

func (j jobCodeDo) Delete(models ...*model.JobCode) (result gen.ResultInfo, err error) {
	return j.DO.Delete(models)
}

func (j *jobCodeDo) withDO(do gen.Dao) *jobCodeDo {
	j.DO = *do.(*gen.DO)
	return j
}
