// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package event_test

import (
	"strings"
	"testing"

	"github.com/luids-io/core/event"
)

func TestEvent(t *testing.T) {
	//try event creation
	e := event.New(1234, event.Low)
	if e.ID != "" {
		t.Error("id must be empty")
	}
	if e.Type != event.Undefined {
		t.Error("bad type")
	}
	if e.Code != 1234 {
		t.Error("bad code")
	}
	if e.Created.IsZero() {
		t.Error("bad created timestamp")
	}
	if e.Level != event.Low {
		t.Error("bad event level")
	}
	if !event.GetDefaultSource().Equals(e.Source) {
		t.Error("bad event source")
	}
	if len(e.Data) > 0 {
		t.Error("bad event data")
	}
	if len(e.Processors) > 0 {
		t.Error("bad event processors")
	}
	err := e.Set("prueba", "prueba texto")
	if err != nil {
		t.Error("bad set error")
	}
	err = e.Set(":badfield", "prueba texto")
	if err == nil {
		t.Error("allowed bad field")
	}
	txt, ok := e.Get("prueba")
	if !ok {
		t.Error("can't get field")
	}
	if txt != "prueba texto" {
		t.Error("bad event value")
	}
	_, ok = e.Get("noexiste")
	if ok {
		t.Error("bad field return")
	}
	e.Set("prueba", "actualiza valor")
	txt, _ = e.Get("prueba")
	if txt != "actualiza valor" {
		t.Error("bad update of values")
	}
	e.Set("zz", "zz")
	e.Set("score", 100)
	if strings.Join(e.Fields(), ",") != "prueba,score,zz" {
		t.Errorf("unexpected fields: %v", e.Fields())
	}
	if e.PrintFields() != "prueba=actualiza valor;score=100;zz=zz" {
		t.Errorf("unexpected PrintFields: %s", e.PrintFields())
	}
}
