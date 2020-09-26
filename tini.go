// Package tini defines functions to write and read simple INI files
// from and to suitable structs.
package tini

import (
	"fmt"
	"io"
)

// Write writes an INI representation of v to `out`.
// v needs to be a suitable struct.
func Write(out io.Writer, v interface{}) error {
	res, err := deconstruct(v)
	if err != nil {
		return fmt.Errorf("construction of INI document failed: %v", err)
	}
	return write(out, res)
}

// Read reads an INI representation from `in` and adjusts v accordingly.
// v needs to be a pointer to a suitable struct.
func Read(in io.Reader, v interface{}) error {
	data, err := read(in)
	if err != nil {
		return fmt.Errorf("reading INI input failed: %v", err)
	}
	return construct(v, data)
}

// ReadRaw reads INI data from `in` and returns a raw representation of it.
func ReadRaw(in io.Reader) (Data, error) {
	return read(in)
}

type Data struct {
	Sections []Section
}

func (d *Data) getSection(name string) (Section, bool) {
	for _, s := range d.Sections {
		if s.Name == name {
			return s, true
		}
	}
	return Section{}, false
}

type Section struct {
	Name    string
	Entries []Entry
}

func (s *Section) getValue(name string) (string, bool) {
	for _, e := range s.Entries {
		if e.Key == name {
			return e.Value, true
		}
	}
	return "", false
}

type Entry struct {
	Key   string
	Value string
}
