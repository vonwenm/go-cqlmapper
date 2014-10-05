package cqlm

var rawNameConverter = func(name string) string {
	return name
}

var DefaultMapper *Mapper = NewMapper(rawNameConverter, rawNameConverter, "cqlm")
