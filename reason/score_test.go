// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package reason_test

import (
	"testing"

	"github.com/luids-io/core/reason"
)

func TestScoreClean(t *testing.T) {
	var tests = []struct {
		in   string
		want string
	}{
		{"[score][/score]razon", "razon"},
		{"[score][/score][/score]razon", "[/score]razon"},
		{"[score][/score][score]razon", "[score]razon"},
		{"ra[score][/score]zon", "razon"},
		{"ra[SCORE][/score]zon", "razon"},
		{"[SCORE]232[/score]razon[score][/score]", "razon"},
		{"[SCORE]232[/score]razon[score][/score][/score]", "razon[/score]"},
		{"[score]10[/score][/score]", "[/score]"},
	}
	for _, test := range tests {
		if got := reason.Clean(test.in); got != test.want {
			t.Errorf("Clean(%v) = %v", test.in, got)
		}
	}
}

func TestScoreFromString(t *testing.T) {
	var tests = []struct {
		inStr string
		want  int
	}{
		{"[score]10[/score]", 10},
		{"[score]0[/score]", 0},
		{"[score][/score]", 0},
		{"[score]-10[/score]", -10},
		{"[score]10[/score][/score]", 10},
	}
	for _, test := range tests {
		got, _, err := reason.ExtractScore(test.inStr)
		if err != nil {
			t.Fatalf("ExtractScore(%s): %v", test.inStr, err)
		}
		if got != test.want {
			t.Errorf("want='%v', got='%v'", test.want, got)
		}
	}
}

func TestScoreFromStringErr(t *testing.T) {
	var tests = []struct {
		in      string
		wantErr bool
	}{
		{"[score]10[/score]", false},
		{"[score]10", false},
		{"[score]10[/SCORE]", false},
		{"[score]1.5[/score]", true},
		{"[score]1 0[/score]", true},
		{"[score]    10[/score]", true},
		{"[score]10[/score][/score]", false},
		{"[score][score]10[/score]", true},
		{"[/score][score]10[/score]", false},
		{"[score]10[/score][/score]", false},
	}
	for _, test := range tests {
		_, _, err := reason.ExtractScore(test.in)
		if err != nil && !test.wantErr {
			t.Errorf("ExtractScore(%s) unexpected err: %v", test.in, err)
		} else if err == nil && test.wantErr {
			t.Errorf("ExtractScore(%s) expected err", test.in)
		}
	}
}

func TestExtractScore(t *testing.T) {
	var tests = []struct {
		inStr      string
		wantErr    bool
		wantValue  int
		wantReason string
	}{
		{"raz[score]10[/score]ón", false, 10, "razón"},
		{"raz[score][/score]ón", false, 0, "razón"},
		{"razón", false, 0, "razón"},
		{"raz[score]10.1[/score]ón", true, 0, "razón"},
		{"[score]-1[/score]", false, -1, ""},
		{"r[score]1[/score]a[score]2[/score]zón[/score]", false, 3, "razón[/score]"},
		{"[score]1[/score]ra[score]2[/score]zón[/score]", false, 3, "razón[/score]"},
		{"[/score]mal[/score]ra[score]2[/score]zón[/score]", false, 2, "[/score]mal[/score]razón[/score]"},
		{"r[score][score]1[/score]azón", true, 0, "razón"},
	}
	for _, test := range tests {
		got, reason, err := reason.ExtractScore(test.inStr)
		if err != nil && !test.wantErr {
			t.Fatalf("goterr='%v'", err)
		}
		if test.wantErr && err == nil {
			t.Fatalf("expected error")
		}
		if got != test.wantValue {
			t.Errorf("want='%v', got='%v'", test.wantValue, got)
		}
		if got := reason; got != test.wantReason {
			t.Errorf("want='%v', got='%v'", test.wantReason, got)
		}
	}
}
