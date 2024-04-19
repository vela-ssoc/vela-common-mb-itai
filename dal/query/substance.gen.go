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

func newSubstance(db *gorm.DB, opts ...gen.DOOption) substance {
	_substance := substance{}

	_substance.substanceDo.UseDB(db, opts...)
	_substance.substanceDo.UseModel(&model.Substance{})

	tableName := _substance.substanceDo.TableName()
	_substance.ALL = field.NewAsterisk(tableName)
	_substance.ID = field.NewInt64(tableName, "id")
	_substance.Name = field.NewString(tableName, "name")
	_substance.Icon = field.NewBytes(tableName, "icon")
	_substance.Hash = field.NewString(tableName, "hash")
	_substance.Desc = field.NewString(tableName, "desc")
	_substance.Chunk = field.NewBytes(tableName, "chunk")
	_substance.Links = field.NewField(tableName, "links")
	_substance.MinionID = field.NewInt64(tableName, "minion_id")
	_substance.Version = field.NewInt64(tableName, "version")
	_substance.CreatedID = field.NewInt64(tableName, "created_id")
	_substance.UpdatedID = field.NewInt64(tableName, "updated_id")
	_substance.CreatedAt = field.NewTime(tableName, "created_at")
	_substance.UpdatedAt = field.NewTime(tableName, "updated_at")

	_substance.fillFieldMap()

	return _substance
}

type substance struct {
	substanceDo substanceDo

	ALL       field.Asterisk
	ID        field.Int64
	Name      field.String
	Icon      field.Bytes
	Hash      field.String
	Desc      field.String
	Chunk     field.Bytes
	Links     field.Field
	MinionID  field.Int64
	Version   field.Int64
	CreatedID field.Int64
	UpdatedID field.Int64
	CreatedAt field.Time
	UpdatedAt field.Time

	fieldMap map[string]field.Expr
}

func (s substance) Table(newTableName string) *substance {
	s.substanceDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s substance) As(alias string) *substance {
	s.substanceDo.DO = *(s.substanceDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *substance) updateTableName(table string) *substance {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewInt64(table, "id")
	s.Name = field.NewString(table, "name")
	s.Icon = field.NewBytes(table, "icon")
	s.Hash = field.NewString(table, "hash")
	s.Desc = field.NewString(table, "desc")
	s.Chunk = field.NewBytes(table, "chunk")
	s.Links = field.NewField(table, "links")
	s.MinionID = field.NewInt64(table, "minion_id")
	s.Version = field.NewInt64(table, "version")
	s.CreatedID = field.NewInt64(table, "created_id")
	s.UpdatedID = field.NewInt64(table, "updated_id")
	s.CreatedAt = field.NewTime(table, "created_at")
	s.UpdatedAt = field.NewTime(table, "updated_at")

	s.fillFieldMap()

	return s
}

func (s *substance) WithContext(ctx context.Context) *substanceDo {
	return s.substanceDo.WithContext(ctx)
}

func (s substance) TableName() string { return s.substanceDo.TableName() }

func (s substance) Alias() string { return s.substanceDo.Alias() }

func (s substance) Columns(cols ...field.Expr) gen.Columns { return s.substanceDo.Columns(cols...) }

func (s *substance) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *substance) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 13)
	s.fieldMap["id"] = s.ID
	s.fieldMap["name"] = s.Name
	s.fieldMap["icon"] = s.Icon
	s.fieldMap["hash"] = s.Hash
	s.fieldMap["desc"] = s.Desc
	s.fieldMap["chunk"] = s.Chunk
	s.fieldMap["links"] = s.Links
	s.fieldMap["minion_id"] = s.MinionID
	s.fieldMap["version"] = s.Version
	s.fieldMap["created_id"] = s.CreatedID
	s.fieldMap["updated_id"] = s.UpdatedID
	s.fieldMap["created_at"] = s.CreatedAt
	s.fieldMap["updated_at"] = s.UpdatedAt
}

func (s substance) clone(db *gorm.DB) substance {
	s.substanceDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s substance) replaceDB(db *gorm.DB) substance {
	s.substanceDo.ReplaceDB(db)
	return s
}

type substanceDo struct{ gen.DO }

func (s substanceDo) Debug() *substanceDo {
	return s.withDO(s.DO.Debug())
}

func (s substanceDo) WithContext(ctx context.Context) *substanceDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s substanceDo) ReadDB() *substanceDo {
	return s.Clauses(dbresolver.Read)
}

func (s substanceDo) WriteDB() *substanceDo {
	return s.Clauses(dbresolver.Write)
}

func (s substanceDo) Session(config *gorm.Session) *substanceDo {
	return s.withDO(s.DO.Session(config))
}

func (s substanceDo) Clauses(conds ...clause.Expression) *substanceDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s substanceDo) Returning(value interface{}, columns ...string) *substanceDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s substanceDo) Not(conds ...gen.Condition) *substanceDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s substanceDo) Or(conds ...gen.Condition) *substanceDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s substanceDo) Select(conds ...field.Expr) *substanceDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s substanceDo) Where(conds ...gen.Condition) *substanceDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s substanceDo) Order(conds ...field.Expr) *substanceDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s substanceDo) Distinct(cols ...field.Expr) *substanceDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s substanceDo) Omit(cols ...field.Expr) *substanceDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s substanceDo) Join(table schema.Tabler, on ...field.Expr) *substanceDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s substanceDo) LeftJoin(table schema.Tabler, on ...field.Expr) *substanceDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s substanceDo) RightJoin(table schema.Tabler, on ...field.Expr) *substanceDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s substanceDo) Group(cols ...field.Expr) *substanceDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s substanceDo) Having(conds ...gen.Condition) *substanceDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s substanceDo) Limit(limit int) *substanceDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s substanceDo) Offset(offset int) *substanceDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s substanceDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *substanceDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s substanceDo) Unscoped() *substanceDo {
	return s.withDO(s.DO.Unscoped())
}

func (s substanceDo) Create(values ...*model.Substance) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s substanceDo) CreateInBatches(values []*model.Substance, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s substanceDo) Save(values ...*model.Substance) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s substanceDo) First() (*model.Substance, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Substance), nil
	}
}

func (s substanceDo) Take() (*model.Substance, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Substance), nil
	}
}

func (s substanceDo) Last() (*model.Substance, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Substance), nil
	}
}

func (s substanceDo) Find() ([]*model.Substance, error) {
	result, err := s.DO.Find()
	return result.([]*model.Substance), err
}

func (s substanceDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Substance, err error) {
	buf := make([]*model.Substance, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s substanceDo) FindInBatches(result *[]*model.Substance, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s substanceDo) Attrs(attrs ...field.AssignExpr) *substanceDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s substanceDo) Assign(attrs ...field.AssignExpr) *substanceDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s substanceDo) Joins(fields ...field.RelationField) *substanceDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s substanceDo) Preload(fields ...field.RelationField) *substanceDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s substanceDo) FirstOrInit() (*model.Substance, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Substance), nil
	}
}

func (s substanceDo) FirstOrCreate() (*model.Substance, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Substance), nil
	}
}

func (s substanceDo) FindByPage(offset int, limit int) (result []*model.Substance, count int64, err error) {
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

func (s substanceDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s substanceDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s substanceDo) Delete(models ...*model.Substance) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *substanceDo) withDO(do gen.Dao) *substanceDo {
	s.DO = *do.(*gen.DO)
	return s
}
