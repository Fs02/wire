package wire

import (
	"errors"
	"reflect"
)

// TODO:
// - check duplicate component
// - check for ambiguous interface

const tag = "wire"

type component struct {
	name         string
	complete     bool
	value        reflect.Value
	dependencies []dependency
}

type dependency struct {
	name  string
	index int
	typ   reflect.Type
}

type group []component

func (gr group) find(name string) (component, bool) {
	for _, c := range gr {
		if c.name == name {
			return c, true
		}
	}

	return component{}, false
}

func (gr group) get(name string) component {
	if c, ok := gr.find(name); ok {
		return c
	}

	panic("wire: no " + gr[0].value.Type().Name() + " identified using \"" + name + "\" found")
}

var components map[reflect.Type]group

func init() {
	Reset()
}

func Reset() {
	components = make(map[reflect.Type]group)
}

func Connect(val interface{}, name ...string) interface{} {
	ptr := false
	rv := reflect.ValueOf(val)
	rt := rv.Type()
	nam := ""

	if len(name) > 0 {
		nam = name[0]
	}

	comp := component{
		name:     nam,
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
		components[rt] = append(components[rt], comp)
		return val
	}

	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		if tval, ok := sf.Tag.Lookup(tag); ok {
			depRt := sf.Type

			if depRt.Kind() == reflect.Ptr {
				depRt = depRt.Elem()
			}

			comp.complete = false
			comp.dependencies = append(comp.dependencies, dependency{
				name:  tval,
				typ:   depRt,
				index: i,
			})
		}
	}

	if !comp.complete && !ptr {
		panic(errors.New("wire: trying to connect incompleted component as a value, use a reference instead"))
	}

	components[rt] = append(components[rt], comp)
	return val
}

func Get(strct interface{}, name ...string) interface{} {
	typ := reflect.TypeOf(strct)
	nam := ""

	if len(name) > 0 {
		nam = name[0]
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if gr, exist := components[typ]; exist {
		if c := gr.get(nam); c.value.CanAddr() {
			return c.value.Addr().Interface()
		}
	}

	panic("wire: addressable value for " + typ.Name() + " not found, try to connect the component using reference instead of pointer")
}

func Apply() {
	for _, gr := range components {
		for _, comp := range gr {
			fill(comp)
		}
	}

	// free memory
	components = nil
}

func fill(c component) {
	if c.complete {
		return
	}

	for i := range c.dependencies {
		dep := c.dependencies[i]
		cdep := component{}

		if gr, exist := components[dep.typ]; exist {
			cdep = gr.get(dep.name)
		} else {
			// scan if it's interface
			matches := 0

			if dep.typ.Kind() == reflect.Interface {
				for _, gr := range components {
					if gr[0].value.Type().Implements(dep.typ) {
						if fcedp, ok := gr.find(dep.name); ok {
							cdep = fcedp
							matches++
						}
					}
				}
			}

			if matches == 0 {
				panic(errors.New("wire: " + c.value.Type().Name() + " requires " + dep.typ.Name() + " identified using \"" + dep.name + "\", but none was found"))
			} else if matches > 1 {
				panic(errors.New("wire: ambiguous connection found on " + c.value.Type().Name() + ", multiple components satisfy " + dep.typ.Name() + " interface, consider using named component"))
			}
		}

		fill(cdep)

		fv := c.value.Field(dep.index)
		if fv.Kind() == reflect.Ptr {
			if !cdep.value.CanAddr() {
				panic(errors.New("wire: " + c.value.Type().Name() + " requires " + dep.typ.Name() + " as pointer, wire as a reference instead of a value"))
			}

			fv.Set(cdep.value.Addr())
		} else {
			fv.Set(cdep.value)
		}
	}
}
