package form

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// TODO this functionality could be moved to a separate repository

// IEncoder converts custom structure to map[string]string representation. This
// format could be useful when sending requests with multipart/form-data or
// application/x-www-form-urlencoded content types.
type IEncoder interface {

	// Values returns a map[string]string representation of initial data.
	Values() (map[string]string, error)

	// AddValue adds new value and its path to the result set.
	AddValue(path []string, value string, omitEmpty bool)

	// HandleField analyzes the value and processes it according to its type kind.
	HandleValue(path []string, v reflect.Value, omitEmpty bool)
}

// ValueCustomProcessor is a type for custom value processors.
type ValueCustomProcessor func(encoder IEncoder, path []string, v reflect.Value, omitEmpty bool) error

// NewEncoder creates *Encoder instance and makes necessary initializations.
func NewEncoder(data interface{}) *Encoder {
	return &Encoder{
		PathToStringConverter: PathToStringConverter,
		ValueCustomProcessors: make(map[reflect.Kind]ValueCustomProcessor),

		initialData: data,
		formData:    make(map[string]string),
	}
}

// Encoder converts custom structure to map[string]string representation. This
// format could be useful when sending requests with multipart/form-data or
// application/x-www-form-urlencoded content types.
type Encoder struct {

	// Is used to convert value's path to a single string.
	// By default it adds surrounding square brackets for all path
	// parts except the first one.
	PathToStringConverter func([]string) string

	// Contains custom value processors for value kinds.
	ValueCustomProcessors map[reflect.Kind]ValueCustomProcessor

	// Tag name that should be used to take field's name. By default is empty,
	// meaning none tag will be used.
	FieldTag string

	initialData interface{}
	initialPath []string
	err         error
	formData    map[string]string
}

// Values returns a map[string]string representation of initial data.
func (e *Encoder) Values() (map[string]string, error) {
	e.HandleValue(e.initialPath, reflect.ValueOf(e.initialData), false)

	return e.formData, e.err
}

// AddValue adds new value and its path to the result set.
func (e *Encoder) AddValue(path []string, value string, omitEmpty bool) {
	if !omitEmpty || value != "" {
		e.formData[e.PathToStringConverter(path)] = value
	}
}

// Gets struct field information. Returns three values:
// - name: is a field name
// - skip: whether or not the field must be skipped
// - omit: whether or not the field must be skipped if its value is empty
func (e *Encoder) fieldInfo(field reflect.StructField) (name string, skip, omit bool) {
	name = field.Name
	if e.FieldTag != "" {
		tagName := field.Tag.Get(e.FieldTag)
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

// HandleValue analyzes the value and processes it according to its type kind.
// Custom value kind processor has a higher priority than a builtin. To set a
// desired behaviour for any value kind you need to add a custom processor to
// ValueCustomProcessors map with value.Kind() as a key.
func (e *Encoder) HandleValue(path []string, v reflect.Value, omitEmpty bool) {
	if customProcessor, found := e.ValueCustomProcessors[v.Kind()]; found {
		if err := customProcessor(e, path, v, omitEmpty); err != nil {
			e.err = err
		}
	} else {
		e.builtinValueProcessor(path, v, omitEmpty)
	}
}

// Default value processor.
func (e *Encoder) builtinValueProcessor(path []string, v reflect.Value, omitEmpty bool) {
	switch v.Kind() {

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			e.HandleValue(append(path, strconv.Itoa(i)), v.Index(i), false)
		}

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if name, skip, omit := e.fieldInfo(v.Type().Field(i)); !skip {
				if v.Type().Field(i).Anonymous {
					e.HandleValue(path, v.Field(i), omit)
				} else {
					e.HandleValue(append(path, name), v.Field(i), omit)
				}
			}
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			e.HandleValue(append(path, fmt.Sprint(key.Interface())), v.MapIndex(key), false)
		}

	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			e.AddValue(path, "", omitEmpty)
		} else {
			e.HandleValue(path, v.Elem(), omitEmpty)
		}

	case reflect.Invalid:
		e.err = fmt.Errorf("unsupported type [%s/%s] and kind [%s]",
			v.Type().PkgPath(), v.Type().Name(), v.Kind().String())

	case reflect.Bool:
		if v.Bool() {
			e.AddValue(path, "true", omitEmpty)
		}
		e.AddValue(path, "", omitEmpty)

	default:
		e.AddValue(path, fmt.Sprint(v.Interface()), omitEmpty)
	}
}

// PathToStringConverter is one of the possible implementations of []string
// to string conversions. It takes a slice of strings, representing a value
// path, and converts it to a single string. It adds surrounding square
// brackets for all path parts except the first one. Here is the example how
// it works: []string{"path", "to", "value", "0"} => "path[to][value][0]"
func PathToStringConverter(path []string) string {
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
