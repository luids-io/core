// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package event_test

import (
	"strings"
	"testing"

	"github.com/luids-io/core/event"
)

func TestEvent(t *testing.T) {
	//try event creation
	e := event.New(event.Security, 1234, event.Low)
	if e.ID != "" {
		t.Error("bad id")
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

func TestRegistry(t *testing.T) {
	//try register code
	event.RegisterCode(1234, "prueba1", "esto es una prueba", []event.FieldMetadata{})
	e := event.New(event.Security, 1234, event.Info)
	if e.Codename() != "prueba1" {
		t.Error("codename misstmatch")
	}
	if e.Desc() != "esto es una prueba" {
		t.Error("description missmatch")
	}
	// check field undefined
	e.Set("message", "la cagaste burlancaster")
	err := event.Validate(e)
	if err == nil {
		t.Error("expected error")
	}
	if !strings.Contains(err.Error(), "undefined") {
		t.Errorf("unexpected error: %v", err)
	}
	event.RegisterCode(1234, "prueba1", "este es el mensaje: [data.message]",
		[]event.FieldMetadata{
			{Name: "message", Type: "string", Required: true},
		})
	err = event.Validate(e)
	if err != nil {
		t.Errorf("unexecter error: %v", err)
	}
	if e.Desc() != "este es el mensaje: la cagaste burlancaster" {
		t.Errorf("unexpected desc: %s", e.Desc())
	}
	// check required field
	event.RegisterCode(1234, "prueba1", "este es el mensaje: [data.message]",
		[]event.FieldMetadata{
			{Name: "message", Type: "string", Required: true},
			{Name: "score", Type: "int", Required: true},
		})
	err = event.Validate(e)
	if err == nil {
		t.Error("expected error")
	}
	if !strings.Contains(err.Error(), "required") {
		t.Errorf("unexpected error: %v", err)
	}
	// test bad int
	e.Set("score", "malvalor")
	err = event.Validate(e)
	if err == nil {
		t.Error("expected error")
	}
	if !strings.Contains(err.Error(), "valid int") {
		t.Errorf("unexpected error: %v", err)
	}
	// test int value
	e.Set("score", 100)
	err = event.Validate(e)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// check don't required field
	event.RegisterCode(1234, "prueba1", "[data.score] [data.prob]",
		[]event.FieldMetadata{
			{Name: "message", Type: "string", Required: true},
			{Name: "score", Type: "int", Required: true},
			{Name: "prob", Type: "float"},
		})
	err = event.Validate(e)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	e.Set("prob", 0.2)
	err = event.Validate(e)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if e.Desc() != "100 0.2" {
		t.Errorf("unexpected desc: %s", e.Desc())
	}
}
