package cqlmapper

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var NotPointerToStructError = errors.New("Instance should be pointer to struct")

type InstanceMapper struct {
	*MapperBuilder
	elem reflect.Value
}

type TableNameable interface {
	TableName() string
}

func quote(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

func (builder *MapperBuilder) NewInstanceMapper(instance interface{}) (*InstanceMapper, error) {
	value := reflect.ValueOf(instance)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return nil, NotPointerToStructError
	}

	elem := value.Elem()
	return &InstanceMapper{
		MapperBuilder: builder,
		elem:          elem,
	}, nil
}

func (mapper *InstanceMapper) typeName() string {
	return mapper.elem.Type().Name()
}

func typeFieldNames(valueType reflect.Type) []string {
	fieldNames := make([]string, 0)
	for fieldIndex := 0; fieldIndex < valueType.NumField(); fieldIndex++ {
		field := valueType.Field(fieldIndex)
		if field.Type.Name() == field.Name {
			fieldNames = append(fieldNames, typeFieldNames(field.Type)...)
		} else {
			fieldNames = append(fieldNames, field.Name)
		}
	}

	return fieldNames
}

func (mapper *InstanceMapper) fieldNames() []string {
	valueType := mapper.elem.Type()
	return typeFieldNames(valueType)
}

func (mapper *InstanceMapper) TableName() string {
	var tableName string

	if tableNameable, ok := mapper.elem.Addr().Interface().(TableNameable); ok {
		tableName = tableNameable.TableName()
	} else {
		tableName = mapper.TableNameConverter(mapper.typeName())
	}

	return quote(tableName)
}

func (mapper *InstanceMapper) ColumnNames() []string {
	fieldNames := mapper.fieldNames()
	columnNames := make([]string, len(fieldNames))

	valueType := mapper.elem.Type()

	for fieldIndex, fieldName := range fieldNames {
		field := valueType.Field(fieldIndex)
		columnName := field.Tag.Get(mapper.ColumnNameTag)
		if "" == columnName {
			columnName = mapper.ColumnNameConverter(fieldName)
		}
		columnNames[fieldIndex] = quote(columnName)
	}

	return columnNames
}

func (mapper *InstanceMapper) FieldValues() []interface{} {
	fieldNames := mapper.fieldNames()
	fieldValues := make([]interface{}, len(fieldNames))

	for fieldIndex, fieldName := range fieldNames {
		field := mapper.elem.FieldByName(fieldName)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			fieldValues[fieldIndex] = nil
		} else {
			fieldValues[fieldIndex] = field.Interface()
		}
	}

	return fieldValues
}

func (mapper *InstanceMapper) FieldAddresses() []interface{} {
	fieldNames := mapper.fieldNames()
	fieldAddresses := make([]interface{}, len(fieldNames))

	for fieldIndex, fieldName := range fieldNames {
		field := mapper.elem.FieldByName(fieldName)
		fieldAddresses[fieldIndex] = field.Addr().Interface()
	}

	return fieldAddresses
}

func (mapper *InstanceMapper) SelectQuery() string {
	return fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(mapper.ColumnNames(), ", "),
		mapper.TableName(),
	)
}

func (mapper *InstanceMapper) InsertQuery() string {
	placeholders := make([]string, len(mapper.ColumnNames()))
	for i, _ := range placeholders {
		placeholders[i] = "?"
	}
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES(%s)",
		mapper.TableName(),
		strings.Join(mapper.ColumnNames(), ", "),
		strings.Join(placeholders, ", "),
	)
}

func (mapper *InstanceMapper) DeleteQuery(pkColumnNames ...string) string {
	conditions := make([]string, len(pkColumnNames))
	for index, pkColumnName := range pkColumnNames {
		conditions[index] = fmt.Sprintf("%s = ?", quote(pkColumnName))
	}

	return fmt.Sprintf(
		`DELETE FROM %s WHERE %s`,
		mapper.TableName(),
		strings.Join(conditions, " AND "),
	)
}
