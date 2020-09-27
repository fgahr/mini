package tini

import (
	"os"
	"testing"
)

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

	if err := write(os.Stderr, obj); err != nil {
		t.Errorf("writing failed: %v", err)
	}
}
