package logger

import (
	"errors"
	"reflect"

	"github.com/fatih/structs"
)

type Fields map[string]interface{}

func NewFields() Fields {
	return make(Fields)
}

func (f Fields) Get(k string) interface{} {
	return f[k]
}

func (f Fields) Set(k string, v interface{}) error {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Struct:
		f[k] = toFields(structs.Map(v))
	case reflect.Ptr:
		if rv.Elem().Type().Kind() == reflect.Struct {
			f[k] = toFields(structs.Map(v))
		}
	case reflect.Chan:
		return errors.New("channel is unacceptable value")
	default:
		f[k] = v
	}
	return nil
}

func (f Fields) Merge(nf Fields, ow bool) {
	for k := range nf {
		if isFields(f[k]) && isFields(nf[k]) {
			f[k].(Fields).Merge(nf[k].(Fields), ow)
		}
		if f[k] != nil && nf[k] != nil && !ow {
			continue
		}
		if f[k] != nil && nf[k] == nil {
			continue
		}
		f[k] = nf[k]
	}
}

func (f Fields) MergeWithWarning(nf Fields, ow bool) (Fields, error) {
	n := f
	for k := range nf {
		if isFields(f[k]) && isFields(nf[k]) {
			n, err := f[k].(Fields).MergeWithWarning(nf[k].(Fields), ow)
			if err != nil {
				return nil, err
			}
			return n, nil
		}
		if f[k] != nil && nf[k] != nil && !ow {
			return nil, errors.New("merged field already has valued")
		}
		if f[k] != nil && nf[k] == nil {
			n[k] = f[k]
		}
		n[k] = nf[k]
	}
	return n, nil
}

func isFields(v interface{}) bool {
	switch v.(type) {
	case Fields:
		return true
	default:
		return false
	}
}

func toFields(m map[string]interface{}) Fields {
	f := NewFields()
	for k, v := range m {
		switch v := v.(type) {
		case map[string]interface{}:
			f[k] = toFields(v)
		default:
			f[k] = v
		}
	}
	return f
}
