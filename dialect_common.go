package gorm

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type commonDialect struct{}

func (commonDialect) BindVar(i int) string {
	return "$$" // ?
}

func (commonDialect) Quote(key string) string {
	return fmt.Sprintf(`"%s"`, key)
}

func (commonDialect) DataTypeOf(dataValue reflect.Value, tagSettings map[string]string) string {
	var size int
	if num, ok := tagSettings["SIZE"]; ok {
		size, _ = strconv.Atoi(num)
	}

	switch dataValue.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		if _, ok := tagSettings["AUTO_INCREMENT"]; ok {
			return "INTEGER AUTO_INCREMENT"
		}
		return "INTEGER"
	case reflect.Int64, reflect.Uint64:
		if _, ok := tagSettings["AUTO_INCREMENT"]; ok {
			return "BIGINT AUTO_INCREMENT"
		}
		return "BIGINT"
	case reflect.Float32, reflect.Float64:
		return "FLOAT"
	case reflect.String:
		if size > 0 && size < 65532 {
			return fmt.Sprintf("VARCHAR(%d)", size)
		}
		return "VARCHAR(65532)"
	case reflect.Struct:
		if _, ok := dataValue.Interface().(time.Time); ok {
			return "TIMESTAMP"
		}
	default:
		if _, ok := dataValue.Interface().([]byte); ok {
			if size > 0 && size < 65532 {
				return fmt.Sprintf("BINARY(%d)", size)
			}
			return "BINARY(65532)"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s) for commonDialect", dataValue.Type().Name(), dataValue.Kind().String()))
}

func (c commonDialect) HasIndex(scope *Scope, tableName string, indexName string) bool {
	var (
		count        int
		databaseName = c.currentDatabase(scope)
	)
	c.RawScanInt(scope, &count, "SELECT count(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE table_schema = ? AND table_name = ? AND index_name = ?", databaseName, tableName, indexName)
	return count > 0
}

func (commonDialect) RemoveIndex(scope *Scope, indexName string) {
	scope.Err(scope.NewDB().Exec(fmt.Sprintf("DROP INDEX %v ON %v", indexName, scope.QuotedTableName())).Error)
}

func (c commonDialect) HasTable(scope *Scope, tableName string) bool {
	var (
		count        int
		databaseName = c.currentDatabase(scope)
	)
	c.RawScanInt(scope, &count, "SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ? AND table_name = ?", databaseName, tableName)
	return count > 0
}

func (c commonDialect) HasColumn(scope *Scope, tableName string, columnName string) bool {
	var (
		count        int
		databaseName = c.currentDatabase(scope)
	)
	c.RawScanInt(scope, &count, "SELECT count(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = ? AND table_name = ? AND column_name = ?", databaseName, tableName, columnName)
	return count > 0
}

// RawScanInt scans the first column of the first row into the `scan' int pointer.
// This function captures raw query errors and propagates them to the original scope.
func (commonDialect) RawScanInt(scope *Scope, scanPtr *int, query string, args ...interface{}) {
	scope.Err(scope.NewDB().Raw(query, args...).Row().Scan(scanPtr))
}

// RawScanString scans the first column of the first row into the `scan' string pointer.
// This function captures raw query errors and propagates them to the original scope.
func (commonDialect) RawScanString(scope *Scope, scanPtr *string, query string, args ...interface{}) {
	scope.Err(scope.NewDB().Raw(query, args...).Row().Scan(scanPtr))
}

func (commonDialect) currentDatabase(scope *Scope) (name string) {
	scope.Err(scope.NewDB().Raw("SELECT DATABASE()").Row().Scan(&name))
	return
}

func (commonDialect) LimitAndOffsetSQL(limit, offset int) (sql string) {
	if limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset >= 0 {
		sql += fmt.Sprintf(" OFFSET %d", offset)
	}
	return
}

func (commonDialect) SelectFromDummyTable() string {
	return ""
}

func (commonDialect) LastInsertIdReturningSuffix(tableName, columnName string) string {
	return ""
}
