package tini

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func sectionName(line string) string {
	line = strings.Replace(line, "[", "", 1)
	line = strings.Replace(line, "]", "", 1)
	return strings.TrimSpace(line)
}

func parseEntry(line string) (string, string) {
	fragments := strings.Split(line, "=")
	if len(fragments) != 2 {
		panic("parsing entry but splitting at = didn't give two segments: " + line)
	}

	return strings.TrimSpace(fragments[0]), strings.TrimSpace(fragments[1])
}

func read(in io.Reader) (Data, error) {
	scanner := bufio.NewScanner(in)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return Data{}, fmt.Errorf("error reading INI contents: %v", err)
	}

	sectionRegex := regexp.MustCompile(`\[\s*\w+\s*\]\s*`)
	entryRegex := regexp.MustCompile(`\s*\w+\s*=\s*\w*\s*`)

	res := Data{}
	s := Section{}
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], ";") {
			continue
		}

		if sectionRegex.MatchString(lines[i]) {
			if s.Name != "" {
				res.addSection(s)
				s = Section{}
			}
			s.Name = sectionName(lines[i])
			continue
		}

		if entryRegex.MatchString(lines[i]) {
			k, v := parseEntry(lines[i])
			s.addEntry(entry{k, v})
		}
	}
	if s.Name != "" {
		res.addSection(s)
	}

	return res, nil
}

func applyTo(f reflect.Value, s Section) {
	if f.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < f.NumField(); i++ {
		key, ok := fieldName(f.Type().Field(i))
		if !ok {
			continue
		}
		field := f.Field(i)
		if val, ok := s.getValue(key); ok {
			switch f.Type().Field(i).Type.Kind() {
			case reflect.String:
				field.Set(reflect.ValueOf(val))
			case reflect.Int:
				num, err := strconv.Atoi(val)
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(num))
			case reflect.Int8:
				num, err := strconv.ParseInt(val, 10, 8)
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(int8(num)))
			case reflect.Int16:
				num, err := strconv.ParseInt(val, 10, 16)
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(int16(num)))
			case reflect.Int32:
				num, err := strconv.ParseInt(val, 10, 32)
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(int32(num)))
			case reflect.Int64:
				num, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(int64(num)))
			case reflect.Uint:
				num, err := strconv.Atoi(val)
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(uint(num)))
			case reflect.Uint8:
				num, err := strconv.ParseUint(val, 10, 8)
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(uint8(num)))
			case reflect.Uint16:
				num, err := strconv.ParseUint(val, 10, 16)
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(uint16(num)))
			case reflect.Uint32:
				num, err := strconv.ParseUint(val, 10, 32)
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(uint32(num)))
			case reflect.Uint64:
				num, err := strconv.ParseUint(val, 10, 64)
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(uint64(num)))
			case reflect.Bool:
				b, err := strconv.ParseBool(val)
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(b))
			}
		}
	}
}

func construct(v interface{}, res Data) error {
	x := reflect.ValueOf(v)
	if x.Kind() != reflect.Ptr {
		return fmt.Errorf("receiver must be a struct pointer; instead received: %T", v)
	}

	x = x.Elem()
	if x.Kind() != reflect.Struct {
		return fmt.Errorf("receiver must be a struct pointer; instead received: %T", v)
	}

	for i := 0; i < x.NumField(); i++ {
		sname, ok := fieldName(x.Type().Field(i))
		if !ok {
			continue
		}
		if s, ok := res.getSection(sname); ok {
			applyTo(x.Field(i), s)
		}
	}
	return nil
}
