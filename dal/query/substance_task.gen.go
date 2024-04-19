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

func newSubstanceTask(db *gorm.DB, opts ...gen.DOOption) substanceTask {
	_substanceTask := substanceTask{}

	_substanceTask.substanceTaskDo.UseDB(db, opts...)
	_substanceTask.substanceTaskDo.UseModel(&model.SubstanceTask{})

	tableName := _substanceTask.substanceTaskDo.TableName()
	_substanceTask.ALL = field.NewAsterisk(tableName)
	_substanceTask.ID = field.NewInt64(tableName, "id")
	_substanceTask.TaskID = field.NewInt64(tableName, "task_id")
	_substanceTask.MinionID = field.NewInt64(tableName, "minion_id")
	_substanceTask.Inet = field.NewString(tableName, "inet")
	_substanceTask.BrokerID = field.NewInt64(tableName, "broker_id")
	_substanceTask.BrokerName = field.NewString(tableName, "broker_name")
	_substanceTask.Failed = field.NewBool(tableName, "failed")
	_substanceTask.Reason = field.NewString(tableName, "reason")
	_substanceTask.Executed = field.NewBool(tableName, "executed")
	_substanceTask.CreatedAt = field.NewTime(tableName, "created_at")
	_substanceTask.UpdatedAt = field.NewTime(tableName, "updated_at")

	_substanceTask.fillFieldMap()

	return _substanceTask
}

type substanceTask struct {
	substanceTaskDo substanceTaskDo

	ALL        field.Asterisk
	ID         field.Int64
	TaskID     field.Int64
	MinionID   field.Int64
	Inet       field.String
	BrokerID   field.Int64
	BrokerName field.String
	Failed     field.Bool
	Reason     field.String
	Executed   field.Bool
	CreatedAt  field.Time
	UpdatedAt  field.Time

	fieldMap map[string]field.Expr
}

func (s substanceTask) Table(newTableName string) *substanceTask {
	s.substanceTaskDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s substanceTask) As(alias string) *substanceTask {
	s.substanceTaskDo.DO = *(s.substanceTaskDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *substanceTask) updateTableName(table string) *substanceTask {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewInt64(table, "id")
	s.TaskID = field.NewInt64(table, "task_id")
	s.MinionID = field.NewInt64(table, "minion_id")
	s.Inet = field.NewString(table, "inet")
	s.BrokerID = field.NewInt64(table, "broker_id")
	s.BrokerName = field.NewString(table, "broker_name")
	s.Failed = field.NewBool(table, "failed")
	s.Reason = field.NewString(table, "reason")
	s.Executed = field.NewBool(table, "executed")
	s.CreatedAt = field.NewTime(table, "created_at")
	s.UpdatedAt = field.NewTime(table, "updated_at")

	s.fillFieldMap()

	return s
}

func (s *substanceTask) WithContext(ctx context.Context) *substanceTaskDo {
	return s.substanceTaskDo.WithContext(ctx)
}

func (s substanceTask) TableName() string { return s.substanceTaskDo.TableName() }

func (s substanceTask) Alias() string { return s.substanceTaskDo.Alias() }

func (s substanceTask) Columns(cols ...field.Expr) gen.Columns {
	return s.substanceTaskDo.Columns(cols...)
}

func (s *substanceTask) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *substanceTask) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 11)
	s.fieldMap["id"] = s.ID
	s.fieldMap["task_id"] = s.TaskID
	s.fieldMap["minion_id"] = s.MinionID
	s.fieldMap["inet"] = s.Inet
	s.fieldMap["broker_id"] = s.BrokerID
	s.fieldMap["broker_name"] = s.BrokerName
	s.fieldMap["failed"] = s.Failed
	s.fieldMap["reason"] = s.Reason
	s.fieldMap["executed"] = s.Executed
	s.fieldMap["created_at"] = s.CreatedAt
	s.fieldMap["updated_at"] = s.UpdatedAt
}

func (s substanceTask) clone(db *gorm.DB) substanceTask {
	s.substanceTaskDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s substanceTask) replaceDB(db *gorm.DB) substanceTask {
	s.substanceTaskDo.ReplaceDB(db)
	return s
}

type substanceTaskDo struct{ gen.DO }

func (s substanceTaskDo) Debug() *substanceTaskDo {
	return s.withDO(s.DO.Debug())
}

func (s substanceTaskDo) WithContext(ctx context.Context) *substanceTaskDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s substanceTaskDo) ReadDB() *substanceTaskDo {
	return s.Clauses(dbresolver.Read)
}

func (s substanceTaskDo) WriteDB() *substanceTaskDo {
	return s.Clauses(dbresolver.Write)
}

func (s substanceTaskDo) Session(config *gorm.Session) *substanceTaskDo {
	return s.withDO(s.DO.Session(config))
}

func (s substanceTaskDo) Clauses(conds ...clause.Expression) *substanceTaskDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s substanceTaskDo) Returning(value interface{}, columns ...string) *substanceTaskDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s substanceTaskDo) Not(conds ...gen.Condition) *substanceTaskDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s substanceTaskDo) Or(conds ...gen.Condition) *substanceTaskDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s substanceTaskDo) Select(conds ...field.Expr) *substanceTaskDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s substanceTaskDo) Where(conds ...gen.Condition) *substanceTaskDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s substanceTaskDo) Order(conds ...field.Expr) *substanceTaskDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s substanceTaskDo) Distinct(cols ...field.Expr) *substanceTaskDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s substanceTaskDo) Omit(cols ...field.Expr) *substanceTaskDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s substanceTaskDo) Join(table schema.Tabler, on ...field.Expr) *substanceTaskDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s substanceTaskDo) LeftJoin(table schema.Tabler, on ...field.Expr) *substanceTaskDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s substanceTaskDo) RightJoin(table schema.Tabler, on ...field.Expr) *substanceTaskDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s substanceTaskDo) Group(cols ...field.Expr) *substanceTaskDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s substanceTaskDo) Having(conds ...gen.Condition) *substanceTaskDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s substanceTaskDo) Limit(limit int) *substanceTaskDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s substanceTaskDo) Offset(offset int) *substanceTaskDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s substanceTaskDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *substanceTaskDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s substanceTaskDo) Unscoped() *substanceTaskDo {
	return s.withDO(s.DO.Unscoped())
}

func (s substanceTaskDo) Create(values ...*model.SubstanceTask) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s substanceTaskDo) CreateInBatches(values []*model.SubstanceTask, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s substanceTaskDo) Save(values ...*model.SubstanceTask) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s substanceTaskDo) First() (*model.SubstanceTask, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.SubstanceTask), nil
	}
}

func (s substanceTaskDo) Take() (*model.SubstanceTask, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.SubstanceTask), nil
	}
}

func (s substanceTaskDo) Last() (*model.SubstanceTask, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.SubstanceTask), nil
	}
}

func (s substanceTaskDo) Find() ([]*model.SubstanceTask, error) {
	result, err := s.DO.Find()
	return result.([]*model.SubstanceTask), err
}

func (s substanceTaskDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.SubstanceTask, err error) {
	buf := make([]*model.SubstanceTask, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s substanceTaskDo) FindInBatches(result *[]*model.SubstanceTask, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s substanceTaskDo) Attrs(attrs ...field.AssignExpr) *substanceTaskDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s substanceTaskDo) Assign(attrs ...field.AssignExpr) *substanceTaskDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s substanceTaskDo) Joins(fields ...field.RelationField) *substanceTaskDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s substanceTaskDo) Preload(fields ...field.RelationField) *substanceTaskDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s substanceTaskDo) FirstOrInit() (*model.SubstanceTask, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.SubstanceTask), nil
	}
}

func (s substanceTaskDo) FirstOrCreate() (*model.SubstanceTask, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.SubstanceTask), nil
	}
}

func (s substanceTaskDo) FindByPage(offset int, limit int) (result []*model.SubstanceTask, count int64, err error) {
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

func (s substanceTaskDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s substanceTaskDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s substanceTaskDo) Delete(models ...*model.SubstanceTask) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *substanceTaskDo) withDO(do gen.Dao) *substanceTaskDo {
	s.DO = *do.(*gen.DO)
	return s
}
