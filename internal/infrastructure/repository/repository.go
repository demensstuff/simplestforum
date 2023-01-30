package repository

import (
	"context"
	"reflect"
	"simplestforum/internal/domain"
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/domain/service"
	"simplestforum/internal/infrastructure/repository/helpers"
	"time"

	"github.com/gocraft/dbr"
)

// Gateway is an interface which contains basic functions to interact with the database.
type Gateway interface {
	DeleteFrom(table string) *dbr.DeleteStmt
	DeleteBySql(query string, value ...interface{}) *dbr.DeleteStmt
	InsertInto(table string) *dbr.InsertStmt
	InsertBySql(query string, value ...interface{}) *dbr.InsertStmt
	Select(column ...string) *dbr.SelectStmt
	SelectBySql(query string, value ...interface{}) *dbr.SelectStmt
	Update(table string) *dbr.UpdateStmt
	UpdateBySql(query string, value ...interface{}) *dbr.UpdateStmt
}

// DBConn wraps an abstract database connection.
type DBConn struct {
	*dbr.Connection
}

// NewTransaction creates an abstract transaction.
func (r *DBConn) NewTransaction(ctx context.Context) (entity.AbstractTransaction, error) {
	return r.NewSession(nil).BeginTx(ctx, nil)
}

// Wrap wraps the callback into a database transaction.
func (r *DBConn) Wrap(sess entity.Session, f func(tx Gateway) error) (err error) {
	tx, ok := sess.Transaction.(*dbr.Tx)

	if ok && tx != nil {
		err = f(tx)
	} else {
		err = f(r.NewSession(nil))
	}

	if err != nil {
		return domain.NewDBErrorWrap(err)
	}

	return nil
}

// NewRepository creates a list of all Storages.
func NewRepository(db *dbr.Connection) *service.Storages {
	base := &DBConn{db}

	return &service.Storages{
		User:         NewUserRepository(base),
		Section:      NewSectionRepository(base),
		Topic:        NewTopicRepository(base),
		Post:         NewPostRepository(base),
		Notification: NewNotificationRepository(base),
	}
}

// insertNotNil iterates over all fields of insertStruct and adds non-nil values to the InsertStmt.
func insertNotNil(stmt *dbr.InsertStmt, insertStruct interface{}) {
	helpers.ProcessExportedNonEmptyFields(insertStruct, func(value reflect.Value, field reflect.StructField) {
		// If there is no `db` tag or if it is an ID, skip
		column, ok := field.Tag.Lookup("db")
		if !ok || column == "id" {
			return
		}

		// If there is `insert:"false"`, skip
		insert, ok := field.Tag.Lookup("insert")
		if ok && insert == "false" {
			return
		}

		// Otherwise add to the InsertStmt if the value is of a known type
		switch v := value.Interface().(type) {
		case string, int64, bool, time.Time, *string, *int64, *bool, *time.Time:
			stmt.Pair(column, v)
		}
	})
}

// updateNotNil iterates over all fields of updateStruct and adds non-nil values to the UpdateStmt.
func updateNotNil(stmt *dbr.UpdateStmt, updateStruct interface{}) {
	helpers.ProcessExportedNonEmptyFields(updateStruct, func(value reflect.Value, field reflect.StructField) {
		// If there is no `db` tag, skip
		column, ok := field.Tag.Lookup("db")
		if !ok {
			return
		}

		// Otherwise add to the UpdateStmt
		stmt.Set(column, value.Interface())
	})
}

// applyFilters applies relevant filters from df.
func applyFilters(df interface{}) []dbr.Builder {
	var res []dbr.Builder

	helpers.ProcessExportedNonEmptyFields(df, func(value reflect.Value, field reflect.StructField) {
		// If there is no `db` tag, skip
		column, ok := field.Tag.Lookup("db")
		if !ok {
			return
		}

		// If there is no `sign` tag, skip
		sign, ok := field.Tag.Lookup("sign")
		if !ok {
			return
		}

		v := value.Interface()

		var expr dbr.Builder

		switch sign {
		case "=", "==":
			expr = dbr.Eq(column, v)
		case "<>", "!=":
			expr = dbr.Neq(column, v)
		case ">":
			expr = dbr.Gt(column, v)
		case ">=":
			expr = dbr.Gte(column, v)
		case "<":
			expr = dbr.Lt(column, v)
		case "<=":
			expr = dbr.Lte(column, v)
		default:
			panic("applyFilters: unknown sign " + sign)
		}

		res = append(res, expr)
	})

	return res
}
