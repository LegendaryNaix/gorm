package gorm

import "reflect"

func (scope *Scope) getColumnAsArray(columns []string) (results [][]interface{}) {
	indirectScopeValue := scope.IndirectValue()
	switch indirectScopeValue.Kind() {
	case reflect.Slice:
		for i := 0; i < indirectScopeValue.Len(); i++ {
			var result []interface{}
			var object = reflect.Indirect(indirectScopeValue.Index(i))
			for _, column := range columns {
				result = append(result, object.FieldByName(column).Interface())
			}
			results = append(results, result)
		}
	case reflect.Struct:
		var result []interface{}
		for _, column := range columns {
			result = append(result, indirectScopeValue.FieldByName(column).Interface())
		}
		return [][]interface{}{result}
	}
	return
}

func (scope *Scope) getColumnAsScope(column string) *Scope {
	indirectScopeValue := scope.IndirectValue()

	switch indirectScopeValue.Kind() {
	case reflect.Slice:
		if fieldStruct, ok := scope.GetModelStruct().ModelType.FieldByName(column); ok {
			fieldType := fieldStruct.Type
			if fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			results := reflect.New(reflect.SliceOf(reflect.PtrTo(fieldType))).Elem()

			for i := 0; i < indirectScopeValue.Len(); i++ {
				result := reflect.Indirect(reflect.Indirect(indirectScopeValue.Index(i)).FieldByName(column))

				if result.Kind() == reflect.Slice {
					for j := 0; j < result.Len(); j++ {
						if elem := result.Index(j); elem.CanAddr() {
							results = reflect.Append(results, elem.Addr())
						}
					}
				} else if result.CanAddr() {
					results = reflect.Append(results, result.Addr())
				}
			}
			return scope.New(results.Interface())
		}
	case reflect.Struct:
		if field := indirectScopeValue.FieldByName(column); field.CanAddr() {
			return scope.New(field.Addr().Interface())
		}
	}
	return nil
}
