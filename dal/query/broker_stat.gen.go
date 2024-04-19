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

func newBrokerStat(db *gorm.DB, opts ...gen.DOOption) brokerStat {
	_brokerStat := brokerStat{}

	_brokerStat.brokerStatDo.UseDB(db, opts...)
	_brokerStat.brokerStatDo.UseModel(&model.BrokerStat{})

	tableName := _brokerStat.brokerStatDo.TableName()
	_brokerStat.ALL = field.NewAsterisk(tableName)
	_brokerStat.ID = field.NewInt64(tableName, "id")
	_brokerStat.Name = field.NewString(tableName, "name")
	_brokerStat.MemUsed = field.NewUint64(tableName, "mem_used")
	_brokerStat.MemTotal = field.NewUint64(tableName, "mem_total")
	_brokerStat.CPUPercent = field.NewFloat64(tableName, "cpu_percent")
	_brokerStat.CreatedAt = field.NewTime(tableName, "created_at")
	_brokerStat.UpdatedAt = field.NewTime(tableName, "updated_at")

	_brokerStat.fillFieldMap()

	return _brokerStat
}

type brokerStat struct {
	brokerStatDo brokerStatDo

	ALL        field.Asterisk
	ID         field.Int64
	Name       field.String
	MemUsed    field.Uint64
	MemTotal   field.Uint64
	CPUPercent field.Float64
	CreatedAt  field.Time
	UpdatedAt  field.Time

	fieldMap map[string]field.Expr
}

func (b brokerStat) Table(newTableName string) *brokerStat {
	b.brokerStatDo.UseTable(newTableName)
	return b.updateTableName(newTableName)
}

func (b brokerStat) As(alias string) *brokerStat {
	b.brokerStatDo.DO = *(b.brokerStatDo.As(alias).(*gen.DO))
	return b.updateTableName(alias)
}

func (b *brokerStat) updateTableName(table string) *brokerStat {
	b.ALL = field.NewAsterisk(table)
	b.ID = field.NewInt64(table, "id")
	b.Name = field.NewString(table, "name")
	b.MemUsed = field.NewUint64(table, "mem_used")
	b.MemTotal = field.NewUint64(table, "mem_total")
	b.CPUPercent = field.NewFloat64(table, "cpu_percent")
	b.CreatedAt = field.NewTime(table, "created_at")
	b.UpdatedAt = field.NewTime(table, "updated_at")

	b.fillFieldMap()

	return b
}

func (b *brokerStat) WithContext(ctx context.Context) *brokerStatDo {
	return b.brokerStatDo.WithContext(ctx)
}

func (b brokerStat) TableName() string { return b.brokerStatDo.TableName() }

func (b brokerStat) Alias() string { return b.brokerStatDo.Alias() }

func (b brokerStat) Columns(cols ...field.Expr) gen.Columns { return b.brokerStatDo.Columns(cols...) }

func (b *brokerStat) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := b.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (b *brokerStat) fillFieldMap() {
	b.fieldMap = make(map[string]field.Expr, 7)
	b.fieldMap["id"] = b.ID
	b.fieldMap["name"] = b.Name
	b.fieldMap["mem_used"] = b.MemUsed
	b.fieldMap["mem_total"] = b.MemTotal
	b.fieldMap["cpu_percent"] = b.CPUPercent
	b.fieldMap["created_at"] = b.CreatedAt
	b.fieldMap["updated_at"] = b.UpdatedAt
}

func (b brokerStat) clone(db *gorm.DB) brokerStat {
	b.brokerStatDo.ReplaceConnPool(db.Statement.ConnPool)
	return b
}

func (b brokerStat) replaceDB(db *gorm.DB) brokerStat {
	b.brokerStatDo.ReplaceDB(db)
	return b
}

type brokerStatDo struct{ gen.DO }

func (b brokerStatDo) Debug() *brokerStatDo {
	return b.withDO(b.DO.Debug())
}

func (b brokerStatDo) WithContext(ctx context.Context) *brokerStatDo {
	return b.withDO(b.DO.WithContext(ctx))
}

func (b brokerStatDo) ReadDB() *brokerStatDo {
	return b.Clauses(dbresolver.Read)
}

func (b brokerStatDo) WriteDB() *brokerStatDo {
	return b.Clauses(dbresolver.Write)
}

func (b brokerStatDo) Session(config *gorm.Session) *brokerStatDo {
	return b.withDO(b.DO.Session(config))
}

func (b brokerStatDo) Clauses(conds ...clause.Expression) *brokerStatDo {
	return b.withDO(b.DO.Clauses(conds...))
}

func (b brokerStatDo) Returning(value interface{}, columns ...string) *brokerStatDo {
	return b.withDO(b.DO.Returning(value, columns...))
}

func (b brokerStatDo) Not(conds ...gen.Condition) *brokerStatDo {
	return b.withDO(b.DO.Not(conds...))
}

func (b brokerStatDo) Or(conds ...gen.Condition) *brokerStatDo {
	return b.withDO(b.DO.Or(conds...))
}

func (b brokerStatDo) Select(conds ...field.Expr) *brokerStatDo {
	return b.withDO(b.DO.Select(conds...))
}

func (b brokerStatDo) Where(conds ...gen.Condition) *brokerStatDo {
	return b.withDO(b.DO.Where(conds...))
}

func (b brokerStatDo) Order(conds ...field.Expr) *brokerStatDo {
	return b.withDO(b.DO.Order(conds...))
}

func (b brokerStatDo) Distinct(cols ...field.Expr) *brokerStatDo {
	return b.withDO(b.DO.Distinct(cols...))
}

func (b brokerStatDo) Omit(cols ...field.Expr) *brokerStatDo {
	return b.withDO(b.DO.Omit(cols...))
}

func (b brokerStatDo) Join(table schema.Tabler, on ...field.Expr) *brokerStatDo {
	return b.withDO(b.DO.Join(table, on...))
}

func (b brokerStatDo) LeftJoin(table schema.Tabler, on ...field.Expr) *brokerStatDo {
	return b.withDO(b.DO.LeftJoin(table, on...))
}

func (b brokerStatDo) RightJoin(table schema.Tabler, on ...field.Expr) *brokerStatDo {
	return b.withDO(b.DO.RightJoin(table, on...))
}

func (b brokerStatDo) Group(cols ...field.Expr) *brokerStatDo {
	return b.withDO(b.DO.Group(cols...))
}

func (b brokerStatDo) Having(conds ...gen.Condition) *brokerStatDo {
	return b.withDO(b.DO.Having(conds...))
}

func (b brokerStatDo) Limit(limit int) *brokerStatDo {
	return b.withDO(b.DO.Limit(limit))
}

func (b brokerStatDo) Offset(offset int) *brokerStatDo {
	return b.withDO(b.DO.Offset(offset))
}

func (b brokerStatDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *brokerStatDo {
	return b.withDO(b.DO.Scopes(funcs...))
}

func (b brokerStatDo) Unscoped() *brokerStatDo {
	return b.withDO(b.DO.Unscoped())
}

func (b brokerStatDo) Create(values ...*model.BrokerStat) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Create(values)
}

func (b brokerStatDo) CreateInBatches(values []*model.BrokerStat, batchSize int) error {
	return b.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (b brokerStatDo) Save(values ...*model.BrokerStat) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Save(values)
}

func (b brokerStatDo) First() (*model.BrokerStat, error) {
	if result, err := b.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.BrokerStat), nil
	}
}

func (b brokerStatDo) Take() (*model.BrokerStat, error) {
	if result, err := b.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.BrokerStat), nil
	}
}

func (b brokerStatDo) Last() (*model.BrokerStat, error) {
	if result, err := b.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.BrokerStat), nil
	}
}

func (b brokerStatDo) Find() ([]*model.BrokerStat, error) {
	result, err := b.DO.Find()
	return result.([]*model.BrokerStat), err
}

func (b brokerStatDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.BrokerStat, err error) {
	buf := make([]*model.BrokerStat, 0, batchSize)
	err = b.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (b brokerStatDo) FindInBatches(result *[]*model.BrokerStat, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return b.DO.FindInBatches(result, batchSize, fc)
}

func (b brokerStatDo) Attrs(attrs ...field.AssignExpr) *brokerStatDo {
	return b.withDO(b.DO.Attrs(attrs...))
}

func (b brokerStatDo) Assign(attrs ...field.AssignExpr) *brokerStatDo {
	return b.withDO(b.DO.Assign(attrs...))
}

func (b brokerStatDo) Joins(fields ...field.RelationField) *brokerStatDo {
	for _, _f := range fields {
		b = *b.withDO(b.DO.Joins(_f))
	}
	return &b
}

func (b brokerStatDo) Preload(fields ...field.RelationField) *brokerStatDo {
	for _, _f := range fields {
		b = *b.withDO(b.DO.Preload(_f))
	}
	return &b
}

func (b brokerStatDo) FirstOrInit() (*model.BrokerStat, error) {
	if result, err := b.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.BrokerStat), nil
	}
}

func (b brokerStatDo) FirstOrCreate() (*model.BrokerStat, error) {
	if result, err := b.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.BrokerStat), nil
	}
}

func (b brokerStatDo) FindByPage(offset int, limit int) (result []*model.BrokerStat, count int64, err error) {
	result, err = b.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = b.Offset(-1).Limit(-1).Count()
	return
}

func (b brokerStatDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = b.Count()
	if err != nil {
		return
	}

	err = b.Offset(offset).Limit(limit).Scan(result)
	return
}

func (b brokerStatDo) Scan(result interface{}) (err error) {
	return b.DO.Scan(result)
}

func (b brokerStatDo) Delete(models ...*model.BrokerStat) (result gen.ResultInfo, err error) {
	return b.DO.Delete(models)
}

func (b *brokerStatDo) withDO(do gen.Dao) *brokerStatDo {
	b.DO = *do.(*gen.DO)
	return b
}
