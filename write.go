package ini

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"text/template"
)

var temp *template.Template

func (i *ini) addSection(s section) {
	i.Sections = append(i.Sections, s)
}

func (s *section) addEntry(e entry) {
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

func asEntry(f reflect.StructField, v reflect.Value) (entry, bool) {
	e := entry{}
	name, ok := fieldName(f)
	if !ok {
		return e, false
	}
	e.Key = name
	switch v.Kind() {
	case reflect.String:
		e.Value = v.String()
		return e, true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		e.Value = strconv.FormatInt(v.Int(), 10)
		return e, true
		// TODO: handle other types
	default:
		return e, false
	}
}

func asSection(f reflect.StructField, v reflect.Value) (section, bool) {
	s := section{}
	name, ok := fieldName(f)
	if !ok {
		return s, false
	}
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

func deconstruct(v interface{}) (ini, error) {
	x := reflect.ValueOf(v)
	switch x.Kind() {
	case reflect.Struct:
		res := ini{}
		for i := 0; i < x.NumField(); i++ {
			if s, ok := asSection(x.Type().Field(i), x.Field(i)); ok {
				res.addSection(s)
			}
		}
		return res, nil
	default:
		return ini{}, fmt.Errorf("invalid argument: %v", v)
	}
}

func write(out io.Writer, content ini) error {
	return temp.Execute(out, content)
}

func init() {
	temp = template.Must(template.New("ini").
		Parse("{{range .Sections}}[{{.Name}}]\r\n" +
			"{{range .Entries}}{{.Key}} = {{.Value}}\r\n" +
			"{{end}}{{end}}"))
}
