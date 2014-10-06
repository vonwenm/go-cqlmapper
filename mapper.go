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
		println(field.Type.Name())
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
	return mapper.TableNameConverter(mapper.typeName())
}

func (mapper *InstanceMapper) ColumnNames() []string {
	fieldNames := mapper.fieldNames()
	columnNames := make([]string, len(fieldNames))

	valueType := mapper.elem.Type()

	for fieldIndex, fieldName := range fieldNames {
		field := valueType.Field(fieldIndex)
		if columnName := field.Tag.Get(mapper.ColumnNameTag); "" != columnName {
			columnNames[fieldIndex] = columnName
		} else {
			columnNames[fieldIndex] = mapper.ColumnNameConverter(fieldName)
		}
	}

	return columnNames
}

func (mapper *InstanceMapper) FieldPointers() []interface{} {
	fieldNames := mapper.fieldNames()
	fieldInterfaces := make([]interface{}, len(fieldNames))

	for fieldIndex, fieldName := range fieldNames {
		field := mapper.elem.FieldByName(fieldName)
		fieldInterfaces[fieldIndex] = field.Addr().Interface()
	}

	return fieldInterfaces
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
