package form

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func NewEncoder(data interface{}) *Encoder {
	return &Encoder{
		data: data,
		formData: make(map[string]string),
	}
}

// Encoder converts custom structure to map[string]string representation.
// This converter could be useful to prepare form data to pass it to some API.
type Encoder struct {
	data interface{}
	path []string
	err  error
	tag  string
	formData map[string]string
}

// Values returns data map[string]string representation.
func (e *Encoder) Values() (map[string]string, error) {
	e.processType(e.path, reflect.ValueOf(e.data))
	return e.formData, e.err
}

func (e *Encoder) atomName(path []string) string {
	return strings.Join(path, ".")
}

func (e *Encoder) atomValue(v reflect.Value) string {
	return fmt.Sprint(v.Interface())
}

func (e *Encoder) processAtom(path []string, v reflect.Value) {
	e.formData[e.atomName(path)] = e.atomValue(v)
}

func (e *Encoder) fieldName()

func (e *Encoder) processType (path []string, v reflect.Value) {
	switch v.Kind() {

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			e.processType(append(path, strconv.Itoa(i)), v.Index(i))
		}

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			name := v.Type().Field(i).Name
			if e.tag != "" {
				tagName := v.Type().Field(i).Tag.Get(e.tag)
				if tagName != "" {
					name = tagName
				}
			}

			e.processType(append(path, name), v.Field(i))
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			e.processType(append(path, e.atomValue(key)), v.MapIndex(key))
		}

	case reflect.Invalid, reflect.Ptr, reflect.Interface:
		e.err = fmt.Errorf("unsupported type: %s", v.Type())

	default:
		e.processAtom(path, v)
	}
}