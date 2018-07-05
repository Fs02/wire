package wire

import (
	"errors"
	"reflect"
)

// TODO:
// - check duplicate component
// - named component

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
	Init()
}

func Init() {
	components = make(map[reflect.Type]component)
}

func Connect(val interface{}) {
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
		return
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
		panic(errors.New("trying to connect incompleted component as a value, use a reference instead"))
	}

	components[rt] = comp
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

		if cdep, exist := components[dep.fieldType]; exist {
			fill(cdep)

			fv := c.value.Field(dep.fieldIndex)
			if fv.Kind() == reflect.Ptr {
				if !cdep.value.CanAddr() {
					panic(errors.New(c.value.Type().Name() + " requires " + dep.fieldType.Name() + " as pointer, wire as a reference instead of a value"))
				}

				fv.Set(cdep.value.Addr())
			} else {
				fv.Set(cdep.value)
			}
		} else {
			panic(errors.New(c.value.Type().Name() + " requires " + dep.fieldType.Name() + ", but none was found"))
		}
	}
}
