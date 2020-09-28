// Package tini defines functions to write and read simple INI files
// from and to suitable structs.
package mini

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

// Data represents the kind of INI file data supported by this package.
type Data struct {
	Sections []Section
}

// GetSection tries to find the named section in the data.
func (d *Data) GetSection(name string) (Section, bool) {
	for _, s := range d.Sections {
		if s.Name == name {
			return s, true
		}
	}
	return Section{}, false
}

// GetValue tries to find the vallue associated with a key inside a section.
func (d *Data) GetValue(section, key string) (string, bool) {
	if s, ok := d.GetSection(section); ok {
		return s.GetValue(key)
	}
	return "", false
}

// Section represents a named section in an INI document.
type Section struct {
	Name    string
	Entries []Entry
}

// GetValue tries to find the named entry in a section.
func (s *Section) GetValue(name string) (string, bool) {
	for _, e := range s.Entries {
		if e.Key == name {
			return e.Value, true
		}
	}
	return "", false
}

// Entry represents a single key-value pair in an INI section.
type Entry struct {
	Key   string
	Value string
}
