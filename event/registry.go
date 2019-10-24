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
	code Code
	name string
	desc string
}

// RegisterCode register an event code and its message creator
func RegisterCode(code Code, name string, desc string) {
	r := regItem{
		code: code,
		name: name,
		desc: desc,
	}
	registry[code] = r
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
