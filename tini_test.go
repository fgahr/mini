package tini

import (
	"os"
	"testing"
	"time"
)

func TestWriteTagged(t *testing.T) {
	type nested struct {
		Foo string
		Bar time.Duration
	}
	obj := struct {
		Global nested
	}{
		Global: nested{
			Foo: "foo",
			Bar: 8 * time.Millisecond,
		},
	}

	if err := Write(os.Stdout, obj); err != nil {
		t.Errorf("writing failed: %v", err)
	}
}

func TestWriteINI(t *testing.T) {
	obj := Data{
		Sections: []Section{
			{
				Name: "foo",
				Entries: []Entry{
					{Key: "bar", Value: "baz"},
					{Key: "a", Value: "b"},
				},
			},
			{
				Name: "foo2",
				Entries: []Entry{
					{Key: "oh_no", Value: "yay"},
				},
			},
		},
	}

	if err := write(os.Stdout, obj); err != nil {
		t.Errorf("writing failed: %v", err)
	}
}
