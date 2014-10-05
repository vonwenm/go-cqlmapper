package cqlm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gocql/gocql"
)

type NameConverter func(string) string

type Mapper struct {
	TableNameConverter  NameConverter
	ColumnNameConverter NameConverter
	ColumnNameTag       string
}

func NewMapper(tableNameConverter, columnNameConverter NameConverter, columnNameTag string) *Mapper {
	return &Mapper{
		TableNameConverter:  tableNameConverter,
		ColumnNameConverter: columnNameConverter,
		ColumnNameTag:       columnNameTag,
	}
}

func (mapper *Mapper) createTableMap(target interface{}) *TableMap {
	tableMap := TableMap{
		Mapper: mapper,
		value:  reflect.ValueOf(target),
	}
	return &tableMap
}

type TableMap struct {
	*Mapper
	value reflect.Value
}

func (tableMap *TableMap) tableName() string {
	return tableMap.value.Elem().Type().Name()
}

func (tableMap *TableMap) fieldNames() []string {
	valueType := tableMap.value.Elem().Type()
	fieldNames := make([]string, valueType.NumField())

	for fieldIndex := 0; fieldIndex < valueType.NumField(); fieldIndex++ {
		field := valueType.Field(fieldIndex)
		fieldNames[fieldIndex] = field.Name
	}

	return fieldNames
}

func (tableMap *TableMap) convertedTableName() string {
	return tableMap.TableNameConverter(tableMap.tableName())
}

func (tableMap *TableMap) columnNames() []string {
	fieldNames := tableMap.fieldNames()
	columnNames := make([]string, len(fieldNames))

	valueType := tableMap.value.Elem().Type()

	for fieldIndex, fieldName := range fieldNames {
		field := valueType.Field(fieldIndex)
		if columnName := field.Tag.Get(tableMap.ColumnNameTag); "" != columnName {
			columnNames[fieldIndex] = columnName
		} else {
			columnNames[fieldIndex] = tableMap.ColumnNameConverter(fieldName)
		}
	}

	return columnNames
}

func (tableMap *TableMap) fieldInterfaces() []interface{} {
	elem := tableMap.value.Elem()
	fieldInterfaces := make([]interface{}, elem.NumField())

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldInterfaces[i] = field.Addr().Interface()
	}

	return fieldInterfaces
}

type Query struct {
	base string
	args []interface{}
}

func (mapper *Mapper) Select(target interface{}) *Query {
	tableMap := mapper.createTableMap(target)
	query := Query{
		base: fmt.Sprintf(
			"SELECT %s FROM %s",
			strings.Join(tableMap.columnNames(), ", "),
			tableMap.convertedTableName(),
		),
		args: tableMap.fieldInterfaces(),
	}
	return &query
}

func (query *Query) Scan(session *gocql.Session) error {
	return session.Query(query.base).Scan(query.args...)
}
