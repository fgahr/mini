package mini

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

func (d *Data) addSection(s Section) {
	d.Sections = append(d.Sections, s)
}

func (s *Section) addEntry(e Entry) {
	s.Entries = append(s.Entries, e)
}

func fieldName(f reflect.StructField) (string, bool) {
	iniTag := f.Tag.Get("ini")
	switch iniTag {
	case "-":
		return "", false
	case "":
		return f.Name, true
	default:
		return iniTag, true
	}
}

func asEntry(f reflect.StructField, v reflect.Value) (Entry, bool) {
	e := Entry{}
	name, ok := fieldName(f)
	if !ok {
		return e, false
	}
	e.Comment = f.Tag.Get("inicomment")

	e.Key = name
	if m, ok := v.Type().MethodByName("ToINI"); ok {
		res := m.Func.Call([]reflect.Value{v})
		e.Value = res[0].String()
		return e, true
	}

	if m, ok := v.Type().MethodByName("String"); ok {
		res := m.Func.Call([]reflect.Value{v})
		e.Value = res[0].String()
		return e, true
	}

	switch v.Kind() {
	case reflect.String:
		e.Value = v.String()
		return e, true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		e.Value = strconv.FormatInt(v.Int(), 10)
		return e, true
	case reflect.Float32, reflect.Float64:
		e.Value = strconv.FormatFloat(v.Float(), 'g', 12, 64)
		return e, true
	case reflect.Bool:
		e.Value = strconv.FormatBool(v.Bool())
		return e, true
	default:
		return e, false
	}
}

func asSection(f reflect.StructField, v reflect.Value) (Section, bool) {
	s := Section{}
	name, ok := fieldName(f)
	if !ok {
		return s, false
	}
	s.Comment = f.Tag.Get("inicomment")

	s.Name = name
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if e, ok := asEntry(v.Type().Field(i), v.Field(i)); ok {
				s.addEntry(e)
			}
		}
	default:
		return s, false
	}
	return s, true
}

func deconstruct(v interface{}) (Data, error) {
	x := reflect.ValueOf(v)
	// dereference pointer if necessary
	if x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	switch x.Kind() {
	case reflect.Struct:
		res := Data{}
		for i := 0; i < x.NumField(); i++ {
			if s, ok := asSection(x.Type().Field(i), x.Field(i)); ok {
				res.addSection(s)
			}
		}
		return res, nil
	default:
		return Data{}, fmt.Errorf("invalid argument type: %T", v)
	}
}

type presenter struct {
	out io.Writer
	err error
}

func (p *presenter) printf(format string, v ...interface{}) {
	if p.err != nil {
		return
	}

	_, p.err = fmt.Fprintf(p.out, format, v...)
}

func write(out io.Writer, content Data) error {
	p := presenter{out, nil}
	for _, section := range content.Sections {
		if section.Comment != "" {
			for _, line := range strings.Split(section.Comment, "\n") {
				p.printf("; %s\n", line)
			}
		}
		p.printf("[ %s ]\n", section.Name)

		for _, ent := range section.Entries {
			if ent.Comment != "" {
				for _, line := range strings.Split(ent.Comment, "\n") {
					p.printf("; %s\n", line)
				}
			}
			p.printf("%s = %s\n", ent.Key, ent.Value)
		}
	}

	return p.err
}
