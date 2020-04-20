// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

// Package event includes components to implement a simple security event
// notification system.
//
// This package is a work in progress and makes no API stability promises.
package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"
)

// Event stores event info.
type Event struct {
	// info completed on construction
	Code    Code                   `json:"code"`
	Level   Level                  `json:"level"`
	Created time.Time              `json:"created"`
	Source  Source                 `json:"source"`
	Data    map[string]interface{} `json:"data,omitempty"`
	// info completed by Notifier
	ID          string        `json:"id" bson:"_id"`
	Type        Type          `json:"type"`
	Codename    string        `json:"codename"`
	Description string        `json:"description"`
	Processors  []ProcessInfo `json:"processors,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
}

// Code defines the code of the event
type Code int32

// Level defines the level of event
type Level int8

// Level possible values
const (
	Info Level = iota
	Low
	Medium
	High
	Critical
)

// Type defines the type of event
type Type int8

// Type possible values
const (
	Undefined Type = iota
	Security
)

// ProcessInfo stores event processing info
type ProcessInfo struct {
	Received  time.Time `json:"received"`
	Processor Source    `json:"processor"`
}

// Source stores event source information.
type Source struct {
	Hostname string `json:"hostname"`
	Program  string `json:"program"`
	Instance string `json:"instance"`
}

func (s Source) String() string {
	return fmt.Sprintf("%s.%s[%s]", s.Hostname, s.Program, s.Instance)
}

// Equals returns true if sources are equals
func (s Source) Equals(o Source) bool {
	if s.Hostname != o.Hostname || s.Program != o.Program || s.Instance != o.Instance {
		return false
	}
	return true
}

//New event
func New(c Code, l Level) Event {
	return Event{
		Code:    c,
		Level:   l,
		Created: time.Now(),
		Source:  defaultSource,
		Data:    make(map[string]interface{}),
	}
}

// String method to return string of IType
func (i Type) String() string {
	switch i {
	case Undefined:
		return "undefined"
	case Security:
		return "security"
	default:
		return fmt.Sprintf("%#v(unknown)", i)
	}
}

// MarshalJSON implements interface for struct marshalling
func (i Type) MarshalJSON() ([]byte, error) {
	s := ""
	switch i {
	case Undefined:
		s = "undefined"
	case Security:
		s = "security"
	default:
		return nil, fmt.Errorf("invalid value '%v' for security", i)
	}
	return json.Marshal(s)
}

// UnmarshalJSON implements interface for struct unmarshalling
func (i *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "undefined":
		*i = Undefined
		return nil
	case "security":
		*i = Security
		return nil
	default:
		return fmt.Errorf("cannot unmarshal type '%s'", s)
	}
}

// String method to return string of ILevel
func (l Level) String() string {
	switch l {
	case Info:
		return "info"
	case Low:
		return "low"
	case Medium:
		return "medium"
	case High:
		return "high"
	case Critical:
		return "critical"
	}
	return "unknown"
}

// MarshalJSON implements interface for struct marshalling
func (l Level) MarshalJSON() ([]byte, error) {
	s := ""
	switch l {
	case Info:
		s = "info"
	case Low:
		s = "low"
	case Medium:
		s = "medium"
	case High:
		s = "high"
	case Critical:
		s = "critical"
	default:
		return nil, fmt.Errorf("invalid value '%v' for level", l)
	}
	return json.Marshal(s)
}

// UnmarshalJSON implements interface for struct unmarshalling
func (l *Level) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "info":
		*l = Info
		return nil
	case "low":
		*l = Low
		return nil
	case "medium":
		*l = Medium
		return nil
	case "high":
		*l = High
		return nil
	case "critical":
		*l = Critical
		return nil
	default:
		return fmt.Errorf("cannot unmarshal level '%s'", s)
	}
}

// ToEventLevel returns event level from string, raise if ok and error if bad value
func ToEventLevel(s string) (level Level, raise bool, err error) {
	switch s {
	case "":
		raise = false
	case "none":
		raise = false
	case "info":
		raise = true
		level = Info
	case "low":
		raise = true
		level = Low
	case "medium":
		raise = true
		level = Medium
	case "high":
		raise = true
		level = High
	case "critical":
		raise = true
		level = Critical
	default:
		err = errors.New("invalid event level value")
	}
	return
}

// Set appends data to an event
func (e *Event) Set(field string, v interface{}) error {
	if !fieldRegExp.MatchString(field) {
		return errors.New("invalid field")
	}
	e.Data[field] = v
	return nil
}

// Get gets data from an event
func (e *Event) Get(field string) (v interface{}, ok bool) {
	v, ok = e.Data[field]
	return
}

// Fields returns the fields in the event. It returns always a sorted list.
func (e *Event) Fields() []string {
	fields := make([]string, 0, len(e.Data))
	for k := range e.Data {
		fields = append(fields, k)
	}
	sort.Strings(fields)
	return fields
}

// PrintFields returns an string
func (e *Event) PrintFields() string {
	s := ""
	first := true
	for _, field := range e.Fields() {
		if first {
			first = false
		} else {
			s = s + ";"
		}
		value := e.Data[field]
		s = s + fmt.Sprintf("%s=%v", field, value)
	}
	return s
}

// SetDefaultSource allows change default notify events source
func SetDefaultSource(s Source) {
	defaultSource = s
}

// GetDefaultSource returns default notify events source
func GetDefaultSource() Source {
	return defaultSource
}

var defaultSource Source

var fieldRegExp, _ = regexp.Compile(`^[A-Za-z][A-Za-z0-9_\.]*$`)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	defaultSource = Source{
		Hostname: hostname,
		Program:  filepath.Base(os.Args[0]),
		Instance: strconv.Itoa(os.Getpid()),
	}
}
