package cqlmapper

import (
	"fmt"
	"regexp"
	"strings"
)

type NameConverter func(string) string

type MapperBuilder struct {
	TableNameConverter  NameConverter
	ColumnNameConverter NameConverter
	ColumnNameTag       string
}

func NewMapperBuilder(tableNameConverter, columnNameConverter NameConverter, columnNameTag string) *MapperBuilder {
	return &MapperBuilder{
		TableNameConverter:  tableNameConverter,
		ColumnNameConverter: columnNameConverter,
		ColumnNameTag:       columnNameTag,
	}
}

var rawNameConverter = func(name string) string {
	return name
}

var Default *MapperBuilder = NewMapperBuilder(rawNameConverter, rawNameConverter, "cqlm")

var underscoreRegexp = regexp.MustCompile("\\p{Lu}")

var underscoreConverter = func(name string) string {
	name = strings.ToLower(name[0:1]) + name[1:]
	return underscoreRegexp.ReplaceAllStringFunc(name, func(upperChar string) string {
		return fmt.Sprintf("_%s", strings.ToLower(upperChar))
	})
}

var Underscore *MapperBuilder = NewMapperBuilder(underscoreConverter, underscoreConverter, "")
