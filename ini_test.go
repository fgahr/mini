package ini

import (
	"os"
	"testing"
)

func TestWriteINI(t *testing.T) {
	obj := ini{
		Sections: []section{
			{
				Name: "foo",
				Entries: []entry{
					{Key: "bar", Value: "baz"},
					{Key: "a", Value: "b"},
				},
			},
			{
				Name: "foo2",
				Entries: []entry{
					{Key: "oh_no", Value: "yay"},
				},
			},
		},
	}

	if err := write(os.Stderr, obj); err != nil {
		t.Errorf("writing failed: %v", err)
	}
}
