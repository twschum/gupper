package version

import (
	"testing"
)

func TestParse(t *testing.T) {
	tables := []struct {
		in  string
		out Version
		err bool
	}{
		{"", Version{0, 0, 0}, true},
		{"1", Version{1, 0, 0}, false},
		{"v", Version{0, 0, 0}, true},
		{"1.", Version{0, 0, 0}, true},
		{"1.2", Version{1, 2, 0}, false},
		{"1.0.2", Version{1, 0, 2}, false},
		{"1.0.2-pre", Version{0, 0, 0}, true},
		{"1.x.2", Version{0, 0, 0}, true},
		{"1..2", Version{0, 0, 0}, true},
	}
	for _, table := range tables {
		v, err := Parse(table.in)
		if table.err && err == nil {
			t.Errorf("Expected error, got version %v", v)
		}
		if !table.err && err != nil {
			t.Errorf("Unexpected error, got error %v", err)
		} else if !table.err && table.out != v {
			t.Errorf("Wrong result, got %v, want %v", v, table.out)
		}
	}
}

func TestCompare(t *testing.T) {
	const (
		L int = -1
		E int = 0
		G int = 1
	)
	tables := []struct {
		a      Version
		b      Version
		result int
	}{
		{Version{}, Version{}, E},
		{Version{0, 0, 0}, Version{0, 0, 0}, E},
		{Version{0, 0, 1}, Version{0, 0, 1}, E},
		{Version{0, 0, 1}, Version{0, 0, 2}, L},
		{Version{0, 0, 2}, Version{0, 0, 1}, G},
		{Version{0, 1, 1}, Version{0, 1, 1}, E},
		{Version{0, 1, 2}, Version{0, 2, 1}, L},
		{Version{0, 2, 1}, Version{0, 1, 2}, G},
		{Version{1, 9, 1}, Version{1, 9, 1}, E},
		{Version{2, 1, 2}, Version{11, 2, 1}, L},
		{Version{20, 2, 1}, Version{3, 1, 2}, G},
	}
	for _, table := range tables {
		result := table.a.Compare(table.b)
		if result != table.result {
			t.Errorf("Wrong comparison, got %v, want %v", result, table.result)
		}
	}
}
