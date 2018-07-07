package wire

import (
	"reflect"
	"strings"
)

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
	impl  string
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

	panic("wire: no " + gr[0].value.Type().String() + " identified using \"" + name + "\" found")
}

var components map[reflect.Type]group

func init() {
	Clear()
}

// Clear cached components.
func Clear() {
	components = make(map[reflect.Type]group)
}

// Connect a component, optionally identified by name.
func Connect(val interface{}, name ...string) {
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

	if gr, ok := components[rt]; ok {
		if _, ok := gr.find(nam); ok {
			panic("wire: trying to connect component with same type and name")
		}
	}

	if rt.Kind() != reflect.Struct {
		components[rt] = append(components[rt], comp)
		return
	}

	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		if tval, ok := sf.Tag.Lookup(tag); ok {

			if tval == "-" {
				continue
			}

			depRt := sf.Type

			if depRt.Kind() == reflect.Ptr {
				depRt = depRt.Elem()
			}

			nameAndImpl := strings.Split(tval, ",")
			name := nameAndImpl[0]
			impl := ""

			if len(nameAndImpl) > 1 {
				impl = nameAndImpl[1]
			}

			comp.complete = false
			comp.dependencies = append(comp.dependencies, dependency{
				name:  name,
				index: i,
				typ:   depRt,
				impl:  impl,
			})
		} else if (sf.Type.Kind() == reflect.Ptr || sf.Type.Kind() == reflect.Interface) && rv.Field(i).IsNil() {
			panic("wire: nil interface or pointer without wire detected for " + sf.Type.String() + ", to ignore add `wire:\"-\"`")
		}
	}

	if !comp.complete && !ptr {
		panic("wire: trying to connect incompleted component as a value, use a reference instead")
	}

	components[rt] = append(components[rt], comp)
}

// Resolve a component with identified name.
// This should be called only after wiring applied.
// Resolving component multiple times should be avoided, consider caching the component if you need.
func Resolve(out interface{}, name ...string) {
	rv := reflect.ValueOf(out)

	if rv.Type().Kind() != reflect.Ptr {
		panic("wire: resolve parameter must be a pointer")
	}

	rv = rv.Elem()
	rt := rv.Type()

	nam := ""
	if len(name) > 0 {
		nam = name[0]
	}

	if gr, ok := components[rt]; ok {
		rv.Set(gr.get(nam).value)
		return
	}

	panic("wire: no component with type " + rt.String() + " found")
}

// Apply wiring to all components.
func Apply() {
	for _, gr := range components {
		for _, comp := range gr {
			fill(comp)
		}
	}
}

func fill(c component) {
	if c.complete {
		return
	}

	for i := range c.dependencies {
		dep := c.dependencies[i]
		cdep := component{}

		if gr, ok := components[dep.typ]; ok {
			cdep = gr.get(dep.name)
		} else {
			// scan if it's interface
			matches := 0

			if dep.typ.Kind() == reflect.Interface {
				for _, gr := range components {
					ctyp := gr[0].value.Type()

					if ctyp.Implements(dep.typ) && (dep.impl == "" || dep.impl == ctyp.Name()) {
						if fcedp, ok := gr.find(dep.name); ok {
							cdep = fcedp
							matches++
						}
					}
				}
			}

			if matches == 0 {
				panic("wire: " + c.value.Type().String() + " requires " + dep.typ.String() + " identified using \"" + dep.name + "\", but none was found")
			} else if matches > 1 {
				panic("wire: ambiguous connection found on " + c.value.Type().String() + ", multiple components satisfy " + dep.typ.String() + " interface, consider using named component")
			}
		}

		fill(cdep)

		fv := c.value.Field(dep.index)
		if fv.Kind() == reflect.Ptr {
			if !cdep.value.CanAddr() {
				panic("wire: " + c.value.Type().String() + " requires " + dep.typ.String() + " as pointer, connect as a reference instead of a value")
			}

			fv.Set(cdep.value.Addr())
		} else {
			fv.Set(cdep.value)
		}
	}

	c.complete = true
	c.dependencies = nil
}
