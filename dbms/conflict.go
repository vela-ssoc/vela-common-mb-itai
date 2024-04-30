package dbms

import (
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils"
)

func onConflictFunc(c clause.Clause, builder clause.Builder) {
	onConflict, ok := c.Expression.(clause.OnConflict)
	if !ok {
		c.Build(builder)
		return
	}

	builder.WriteString("ON DUPLICATE KEY UPDATE ")
	if len(onConflict.DoUpdates) == 0 {
		if s := builder.(*gorm.Statement).Schema; s != nil {
			var column clause.Column
			onConflict.DoNothing = false

			if s.PrioritizedPrimaryField != nil {
				column = clause.Column{Name: s.PrioritizedPrimaryField.DBName}
			} else if len(s.DBNames) > 0 {
				column = clause.Column{Name: s.DBNames[0]}
			}

			if column.Name != "" {
				onConflict.DoUpdates = []clause.Assignment{{Column: column, Value: column}}
			}

			builder.(*gorm.Statement).AddClause(onConflict)
		}
	}

	for idx, assignment := range onConflict.DoUpdates {
		if idx > 0 {
			builder.WriteByte(',')
		}

		builder.WriteQuoted(assignment.Column)
		builder.WriteByte('=')
		builder.AddVar(builder, assignment.Value)
	}
}

// Create create hook
func Create(config *callbacks.Config) func(db *gorm.DB) {
	supportReturning := utils.Contains(config.CreateClauses, "RETURNING")

	return func(db *gorm.DB) {
		if db.Error != nil {
			return
		}

		if db.Statement.Schema != nil {
			if !db.Statement.Unscoped {
				for _, c := range db.Statement.Schema.CreateClauses {
					db.Statement.AddClause(c)
				}
			}

			if supportReturning && len(db.Statement.Schema.FieldsWithDefaultDBValue) > 0 {
				if _, ok := db.Statement.Clauses["ON CONFLICT"]; !ok {
					if _, ok = db.Statement.Clauses["RETURNING"]; !ok {
						fromColumns := make([]clause.Column, 0, len(db.Statement.Schema.FieldsWithDefaultDBValue))
						for _, field := range db.Statement.Schema.FieldsWithDefaultDBValue {
							fromColumns = append(fromColumns, clause.Column{Name: field.DBName})
						}
						db.Statement.AddClause(clause.Returning{Columns: fromColumns})
					}
				}
			}
		}

		if db.Statement.SQL.Len() == 0 {
			db.Statement.SQL.Grow(180)
			db.Statement.AddClauseIfNotExists(clause.Insert{})
			db.Statement.AddClause(callbacks.ConvertToCreateValues(db.Statement))

			db.Statement.Build(db.Statement.BuildClauses...)
		}

		isDryRun := !db.DryRun && db.Error == nil
		if !isDryRun {
			return
		}

		ok, mode := hasReturning(db, supportReturning)
		if ok {
			if c, ok := db.Statement.Clauses["ON CONFLICT"]; ok {
				if onConflict, _ := c.Expression.(clause.OnConflict); onConflict.DoNothing {
					mode |= gorm.ScanOnConflictDoNothing
				}
			}

			rows, err := db.Statement.ConnPool.QueryContext(
				db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...,
			)
			if db.AddError(err) == nil {
				defer func() {
					db.AddError(rows.Close())
				}()
				gorm.Scan(rows, db, mode)
			}

			return
		}

		result, err := db.Statement.ConnPool.ExecContext(
			db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...,
		)
		if err != nil {
			db.AddError(err)
			return
		}

		db.RowsAffected, _ = result.RowsAffected()
		if db.RowsAffected == 0 {
			return
		}

		var (
			pkField     *schema.Field
			pkFieldName = "@id"
		)

		insertID, err := result.LastInsertId()
		insertOk := err == nil && insertID > 0

		if !insertOk {
			if !supportReturning {
				db.AddError(err)
			}
			return
		}

		if db.Statement.Schema != nil {
			if db.Statement.Schema.PrioritizedPrimaryField == nil || !db.Statement.Schema.PrioritizedPrimaryField.HasDefaultValue {
				return
			}
			pkField = db.Statement.Schema.PrioritizedPrimaryField
			pkFieldName = db.Statement.Schema.PrioritizedPrimaryField.DBName
		}

		// append @id column with value for auto-increment primary key
		// the @id value is correct, when: 1. without setting auto-increment primary key, 2. database AutoIncrementIncrement = 1
		switch values := db.Statement.Dest.(type) {
		case map[string]interface{}:
			values[pkFieldName] = insertID
		case *map[string]interface{}:
			(*values)[pkFieldName] = insertID
		case []map[string]interface{}, *[]map[string]interface{}:
			mapValues, ok := values.([]map[string]interface{})
			if !ok {
				if v, ok := values.(*[]map[string]interface{}); ok {
					if *v != nil {
						mapValues = *v
					}
				}
			}

			if config.LastInsertIDReversed {
				insertID -= int64(len(mapValues)-1) * schema.DefaultAutoIncrementIncrement
			}

			for _, mapValue := range mapValues {
				if mapValue != nil {
					mapValue[pkFieldName] = insertID
				}
				insertID += schema.DefaultAutoIncrementIncrement
			}
		default:
			if pkField == nil {
				return
			}

			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				if config.LastInsertIDReversed {
					for i := db.Statement.ReflectValue.Len() - 1; i >= 0; i-- {
						rv := db.Statement.ReflectValue.Index(i)
						if reflect.Indirect(rv).Kind() != reflect.Struct {
							break
						}

						_, isZero := pkField.ValueOf(db.Statement.Context, rv)
						if isZero {
							db.AddError(pkField.Set(db.Statement.Context, rv, insertID))
							insertID -= pkField.AutoIncrementIncrement
						}
					}
				} else {
					for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
						rv := db.Statement.ReflectValue.Index(i)
						if reflect.Indirect(rv).Kind() != reflect.Struct {
							break
						}

						if _, isZero := pkField.ValueOf(db.Statement.Context, rv); isZero {
							db.AddError(pkField.Set(db.Statement.Context, rv, insertID))
							insertID += pkField.AutoIncrementIncrement
						}
					}
				}
			case reflect.Struct:
				_, isZero := pkField.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
				if isZero {
					db.AddError(pkField.Set(db.Statement.Context, db.Statement.ReflectValue, insertID))
				}
			}
		}
	}
}

func hasReturning(tx *gorm.DB, supportReturning bool) (bool, gorm.ScanMode) {
	if supportReturning {
		if c, ok := tx.Statement.Clauses["RETURNING"]; ok {
			returning, _ := c.Expression.(clause.Returning)
			if len(returning.Columns) == 0 || (len(returning.Columns) == 1 && returning.Columns[0].Name == "*") {
				return true, 0
			}
			return true, gorm.ScanUpdate
		}
	}
	return false, 0
}
