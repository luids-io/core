// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package xlist

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Format stores the type of available formats for lists
type Format int

// List of formats
const (
	Plain Format = iota
	CIDR
	Sub
)

func (f Format) string() string {
	switch f {
	case Plain:
		return "plain"
	case CIDR:
		return "cidr"
	case Sub:
		return "sub"
	default:
		return ""
	}
}

// String implements stringer interface
func (f Format) String() string {
	s := f.string()
	if s == "" {
		return fmt.Sprintf("unkown(%d)", f)
	}
	return s
}

// ToFormat returns the format type from its string representation
func ToFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "plain":
		return Plain, nil
	case "cidr":
		return CIDR, nil
	case "sub":
		return Sub, nil
	default:
		return Format(-1), fmt.Errorf("invalid format %s", s)
	}
}

// MarshalJSON implements interface for struct marshalling
func (f Format) MarshalJSON() ([]byte, error) {
	s := f.string()
	if s == "" {
		return nil, fmt.Errorf("invalid value %v for format", f)
	}
	return json.Marshal(s)
}

// UnmarshalJSON implements interface for struct unmarshalling
func (f *Format) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "plain":
		*f = Plain
		return nil
	case "cidr":
		*f = CIDR
		return nil
	case "sub":
		*f = Sub
		return nil
	default:
		return fmt.Errorf("cannot unmarshal format %s", s)
	}
}
