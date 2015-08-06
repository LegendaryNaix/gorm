package gorm

import (
	"fmt"
	"reflect"
	"time"
)

type mssql struct {
	commonDialect
}

func (mssql) HasTop() bool {
	return true
}

func (mssql) SqlTag(value reflect.Value, size int, autoIncrease bool) string {
	switch value.Kind() {
	case reflect.Bool:
		return "bit"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		if autoIncrease {
			return "int IDENTITY(1,1)"
		}
		return "int"
	case reflect.Int64, reflect.Uint64:
		if autoIncrease {
			return "bigint IDENTITY(1,1)"
		}
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.String:
		if size > 0 && size < 65532 {
			return fmt.Sprintf("nvarchar(%d)", size)
		}
		return "text"
	case reflect.Struct:
		if _, ok := value.Interface().(time.Time); ok {
			return "datetime2"
		}
	default:
		if _, ok := value.Interface().([]byte); ok {
			if size > 0 && size < 65532 {
				return fmt.Sprintf("varchar(%d)", size)
			}
			return "text"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s) for mssql", value.Type().Name(), value.Kind().String()))
}

func (s mssql) HasTable(scope *Scope, tableName string) bool {
	var (
		count        int
		databaseName string
	)
	s.CurrentDatabase(scope, &databaseName)
	scope.NewDB().Raw("SELECT count(*) FROM INFORMATION_SCHEMA.tables WHERE table_name = ? AND table_catalog = ?", tableName, databaseName).Row().Scan(&count)
	return count > 0
}

func (s mssql) HasColumn(scope *Scope, tableName string, columnName string) bool {
	var (
		count        int
		databaseName string
	)
	s.CurrentDatabase(scope, &databaseName)
	scope.NewDB().Raw("SELECT count(*) FROM information_schema.columns WHERE table_catalog = ? AND table_name = ? AND column_name = ?", databaseName, tableName, columnName).Row().Scan(&count)
	return count > 0
}

func (mssql) HasIndex(scope *Scope, tableName string, indexName string) bool {
	var count int
	scope.NewDB().Raw("SELECT count(*) FROM sys.indexes WHERE name=? AND object_id=OBJECT_ID(?)", indexName, tableName).Row().Scan(&count)
	return count > 0
}

func (mssql) CurrentDatabase(scope *Scope, name *string) {
	scope.Err(scope.NewDB().Raw("SELECT DB_NAME() AS [Current Database]").Row().Scan(name))
}
