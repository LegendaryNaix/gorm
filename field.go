package gorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

type Field struct {
	*StructField
	IsBlank bool
	Field   reflect.Value
}

func (field *Field) Set(value interface{}) (err error) {
	if !field.Field.IsValid() {
		return errors.New("field value not valid")
	}

	if !field.Field.CanAddr() {
		return errors.New("unaddressable value")
	}

	reflectValue, ok := value.(reflect.Value)
	if !ok {
		reflectValue = reflect.ValueOf(value)
	}

	fieldValue := field.Field
	if reflectValue.IsValid() {
		if reflectValue.Type().ConvertibleTo(fieldValue.Type()) {
			fieldValue.Set(reflectValue.Convert(fieldValue.Type()))
		} else {
			if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(field.Struct.Type.Elem()))
				}
				fieldValue = fieldValue.Elem()
			}

			if reflectValue.Type().ConvertibleTo(fieldValue.Type()) {
				fieldValue.Set(reflectValue.Convert(fieldValue.Type()))
			} else if scanner, ok := fieldValue.Addr().Interface().(sql.Scanner); ok {
				err = scanner.Scan(reflectValue.Interface())
			} else {
				err = fmt.Errorf("could not convert argument of field %s from %s to %s", field.Name, reflectValue.Type(), fieldValue.Type())
			}
		}
	} else {
		field.Field.Set(reflect.Zero(field.Field.Type()))
	}

	field.IsBlank = isBlank(field.Field)
	return nil
}

// Fields get value's fields
func (scope *Scope) Fields() map[string]*Field {
	if scope.fields == nil {
		var (
			fields             = map[string]*Field{}
			indirectScopeValue = scope.IndirectValue()
			isStruct           = indirectScopeValue.Kind() == reflect.Struct
		)

		for _, structField := range scope.GetModelStruct().StructFields {
			if field, ok := fields[structField.DBName]; !ok || field.IsIgnored {
				if isStruct {
					fieldValue := indirectScopeValue
					for _, name := range structField.Names {
						fieldValue = reflect.Indirect(fieldValue).FieldByName(name)
					}
					fields[structField.DBName] = &Field{StructField: structField, Field: fieldValue, IsBlank: isBlank(fieldValue)}
				} else {
					fields[structField.DBName] = &Field{StructField: structField, IsBlank: true}
				}
			}
		}

		scope.fields = fields
	}
	return scope.fields
}
