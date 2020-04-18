// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package event

import (
	"fmt"
	"regexp"
	"strings"
)

//global registry mantained by the package
var registry = make(map[Code]regItem)

type regItem struct {
	code   Code
	name   string
	desc   string
	fields []FieldMetadata
}

// FieldMetadata stores metadata of fields
type FieldMetadata struct {
	Name     string
	Type     string
	Required bool
}

// RegisterCode register an event code and its message creator
func RegisterCode(code Code, name string, desc string, fields []FieldMetadata) {
	r := regItem{
		code: code,
		name: name,
		desc: desc,
	}
	if len(fields) > 0 {
		r.fields = make([]FieldMetadata, len(fields), len(fields))
		copy(r.fields, fields)
	}
	registry[code] = r
}

func (i regItem) getFields() map[string]FieldMetadata {
	fields := make(map[string]FieldMetadata, len(i.fields))
	for _, f := range i.fields {
		fields[f.Name] = f
	}
	return fields
}

// Codename returns name of associated code
func (e *Event) Codename() string {
	r, ok := registry[e.Code]
	if !ok {
		message := fmt.Sprintf("unregistered(%v)", e.Code)
		return message
	}
	return r.name
}

// Desc returns description
func (e *Event) Desc() string {
	r, ok := registry[e.Code]
	if !ok {
		return fmt.Sprintf("can't get description for event (%v)", e.Code)
	}
	if reBetweenBrackets.MatchString(r.desc) {
		return reBetweenBrackets.ReplaceAllStringFunc(r.desc, e.replacer)
	}
	return r.desc
}

var reBetweenBrackets = regexp.MustCompile(`\[([^\[\]]*)\]`)

func (e *Event) replacer(s string) string {
	element := strings.Trim(s, "[")
	element = strings.Trim(element, "]")
	if strings.HasPrefix(element, "data.") {
		field := strings.TrimPrefix(element, "data.")
		value, ok := e.Data[field]
		if !ok {
			return s
		}
		return fmt.Sprintf("%v", value)
	}
	return s
}
