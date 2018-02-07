package form

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// NewEncoder creates new *Encoder instance and makes necessary initializations.
func NewEncoder(data interface{}) *Encoder {
	return &Encoder{
		initialData: data,
		formData:    make(map[string]string),

		PathToString: PathToString,
	}
}

// Encoder converts custom structure to map[string]string representation.
type Encoder struct {
	PathToString func([]string) string
	Tag          string

	initialData interface{}
	initialPath []string
	err         error
	formData    map[string]string
}

// Values returns a map[string]string representation of initial data.
func (e *Encoder) Values() (map[string]string, error) {
	e.processType(e.initialPath, reflect.ValueOf(e.initialData), false)

	return e.formData, e.err
}

// Gets struct field information. Returns three values:
//    name: is a field name
//    skip: whether or not the value must be skipped
//    omit: whether or not the value must be skipped if it is empty
func (e *Encoder) fieldInfo(field reflect.StructField) (name string, skip, omit bool) {
	name = field.Name
	if e.Tag != "" {
		tagName := field.Tag.Get(e.Tag)
		if tagName == "-" {
			skip = true
		} else if tagName != "" {
			name = tagName
			if i1 := strings.Index(tagName, ","); i1 != -1 {
				name = tagName[:i1]
				if i2 := strings.Index(tagName[i1:], "omitempty"); i2 != -1 {
					omit = true
				}
			}
		}
	}

	return name, skip, omit
}

// Analyzes the data and processes it according to its type.
func (e *Encoder) processType(path []string, v reflect.Value, omitEmpty bool) {
	switch v.Kind() {

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			e.processType(append(path, strconv.Itoa(i)), v.Index(i), false)
		}

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if name, skip, omit := e.fieldInfo(v.Type().Field(i)); !skip {
				e.processType(append(path, name), v.Field(i), omit)
			}
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			e.processType(append(path, e.atomValue(key)), v.MapIndex(key), false)
		}

	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			if !omitEmpty {
				e.formData[e.PathToString(path)] = ""
			}
		} else {
			e.processType(path, v.Elem(), omitEmpty)
		}

	case reflect.Invalid:
		e.err = fmt.Errorf("unsupported type [%s/%s] and kind [%s]",
			v.Type().PkgPath(), v.Type().Name(), v.Kind().String())

	default:
		if value := e.atomValue(v); !omitEmpty || value != "" {
			e.formData[e.PathToString(path)] = value
		}
	}
}

// This function converts exact values to its string representation.
// It handles all types aside from following list:
// Slice, Array, Struct, Map, Ptr, Interface, Invalid.
func (e *Encoder) atomValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}

		return ""

	default:
		return fmt.Sprint(v.Interface())
	}
}

// PathToString is one of the possible implementations of []string to string
// conversion. It takes a slice of strings, representing a field path, and
// converts it to a single string. It adds surrounding brackets for all path
// parts except the first one. Here is the example how it works:
// []string{"path", "to", "value"} => "path[to][value]"
func PathToString(path []string) string {
	if len(path) == 0 {
		return ""
	}

	b := new(bytes.Buffer)
	b.WriteString(path[0])
	for i := 1; i < len(path); i++ {
		b.WriteRune('[')
		b.WriteString(path[i])
		b.WriteRune(']')
	}

	return b.String()
}
