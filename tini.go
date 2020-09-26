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
	res, err := read(in)
	if err != nil {
		return fmt.Errorf("reading INI input failed: %v", err)
	}
	return construct(v, res)
}

// ReadRaw reads INI data from `in` and returns a raw representation of it.
func ReadRaw(in io.Reader) (map[string]map[string]string, error) {
	data, err := read(in)
	if err != nil {
		return nil, fmt.Errorf("reading INI input failed: %v", err)
	}

	res := make(map[string]map[string]string)
	for _, sec := range data.Sections {
		rawSec := make(map[string]string)
		for _, e := range sec.Entries {
			rawSec[e.Key] = e.Value
		}
		res[sec.Name] = rawSec
	}
	return res, nil
}

type data struct {
	Sections []section
}

func (i *data) getSection(name string) (section, bool) {
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
