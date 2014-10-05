package cqlm

import (
	"fmt"
	"regexp"
	"strings"
)

var underscoreRegexp = regexp.MustCompile("\\p{Lu}")

var underscoreConverter = func(name string) string {
	name = strings.ToLower(name[0:1]) + name[1:]
	return underscoreRegexp.ReplaceAllStringFunc(name, func(upperChar string) string {
		return fmt.Sprintf("_%s", strings.ToLower(upperChar))
	})
}

var UnderscoreMapper = NewMapper(underscoreConverter, underscoreConverter, "")
