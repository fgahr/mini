package ini

import (
	"fmt"
	"io"
)

func Write(out io.Writer, v interface{}) error {
	res, err := deconstruct(v)
	if err != nil {
		return fmt.Errorf("construction of INI document failed: %v", err)
	}
	return write(out, res)
}

func Read(in io.Reader, v interface{}) error {
	res, err := read(in)
	if err != nil {
		return fmt.Errorf("reading INI input failed: %v", err)
	}
	return construct(v, res)
}

type ini struct {
	Sections []section
}

func (i *ini) getSection(name string) (section, bool) {
	for _, s := range i.Sections {
		if s.Name == name {
			return s, true
		}
	}
	return section{}, false
}

type section struct {
	Name    string
	Entries []entry
}

func (s *section) getValue(name string) (string, bool) {
	for _, e := range s.Entries {
		if e.Key == name {
			return e.Value, true
		}
	}
	return "", false
}

type entry struct {
	Key   string
	Value string
}
