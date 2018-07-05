package wire

import (
	"errors"
	"reflect"
)

// TODO:
// - check duplicate component
// - named component
// - check for ambiguous interface

const tag = "wire"

type component struct {
	complete     bool
	value        reflect.Value
	dependencies []dependency
}

type dependency struct {
	fieldType  reflect.Type
	fieldIndex int
}

var components map[reflect.Type]component

func init() {
	Reset()
}

func Reset() {
	components = make(map[reflect.Type]component)
}

func Connect(val interface{}) interface{} {
	ptr := false
	rv := reflect.ValueOf(val)
	rt := rv.Type()

	comp := component{
		value:    rv,
		complete: true,
	}

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
		comp.value = rv
		ptr = true
	}

	if rt.Kind() != reflect.Struct {
		components[rt] = comp
		return val
	}

	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		if _, ok := sf.Tag.Lookup(tag); ok {
			depRt := sf.Type

			if depRt.Kind() == reflect.Ptr {
				depRt = depRt.Elem()
			}

			comp.complete = false
			comp.dependencies = append(comp.dependencies, dependency{
				fieldType:  depRt,
				fieldIndex: i,
			})
		}
	}

	if !comp.complete && !ptr {
		panic(errors.New("wire: trying to connect incompleted component as a value, use a reference instead"))
	}

	components[rt] = comp
	return val
}

func Get(strct interface{}) interface{} {
	typ := reflect.TypeOf(strct)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if c, exist := components[typ]; exist && c.value.CanAddr() {
		return c.value.Addr().Interface()
	}

	panic("wire: addressable value for " + typ.Name() + " not found, try to connect the component using reference instead of pointer")
}

func Apply() {
	for _, comp := range components {
		fill(comp)
	}

	// free memory
	components = nil
}

func fill(c component) {
	if c.complete {
		return
	}

	for i := range c.dependencies {
		dep := &c.dependencies[i]

		cdep, exist := components[dep.fieldType]

		if !exist {
			if dep.fieldType.Kind() == reflect.Interface {
				for _, lc := range components {
					if lc.value.Type().Implements(dep.fieldType) {
						cdep = lc
						goto Fill
					}
				}
			}

			panic(errors.New("wire: " + c.value.Type().Name() + " requires " + dep.fieldType.Name() + ", but none was found"))
		}

	Fill:
		fill(cdep)

		fv := c.value.Field(dep.fieldIndex)
		if fv.Kind() == reflect.Ptr {
			if !cdep.value.CanAddr() {
				panic(errors.New("wire: " + c.value.Type().Name() + " requires " + dep.fieldType.Name() + " as pointer, wire as a reference instead of a value"))
			}

			fv.Set(cdep.value.Addr())
		} else {
			fv.Set(cdep.value)
		}
	}
}
