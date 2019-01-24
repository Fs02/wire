package wire

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const tag = "wire"

type component struct {
	id           string
	value        reflect.Value
	dependencies []dependency
	declaredAt   string
}

type dependency struct {
	id    string
	name  string
	index int
	typ   reflect.Type
	impl  string
}

type group []component

func (gr group) find(id string) (component, bool) {
	for _, c := range gr {
		if c.id == id {
			return c, true
		}
	}

	return component{}, false
}

func (gr group) get(id string) component {
	if c, ok := gr.find(id); ok {
		return c
	}

	panic(idNotFoundError{id: id, component: gr[0]})
}

// Container provides an isolated container for DI.
type Container struct {
	components map[reflect.Type]group
	callerSkip int
}

// New create new isolated DI container.
func New() Container {
	return Container{
		components: make(map[reflect.Type]group),
	}
}

// Connect a component, optionally identified by id.
func (container Container) Connect(val interface{}, id ...string) {
	ptr := false
	rv := reflect.ValueOf(val)
	rt := rv.Type()
	nam := ""
	_, file, no, _ := runtime.Caller(container.callerSkip + 1)

	if len(id) > 0 {
		nam = id[0]
	}

	comp := component{
		id:         nam,
		value:      rv,
		declaredAt: file + ":" + strconv.Itoa(no),
	}

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
		comp.value = rv
		ptr = true
	}

	if gr, ok := container.components[rt]; ok {
		if comp, ok := gr.find(nam); ok {
			panic(duplicateError{previous: comp})
		}
	}

	if rt.Kind() != reflect.Struct {
		container.components[rt] = append(container.components[rt], comp)
		return
	}

	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)

		// skip unexported field
		if sf.Name[0] >= 'a' && sf.Name[0] <= 'z' {
			continue
		}

		if tval, ok := sf.Tag.Lookup(tag); ok {

			if tval == "-" {
				continue
			}

			depRt := sf.Type

			if depRt.Kind() == reflect.Ptr {
				depRt = depRt.Elem()
			}

			idAndImpl := strings.Split(tval, ",")
			id := idAndImpl[0]
			impl := ""

			if len(idAndImpl) > 1 {
				impl = idAndImpl[1]
			}

			comp.dependencies = append(comp.dependencies, dependency{
				id:    id,
				name:  sf.Name,
				index: i,
				typ:   depRt,
				impl:  impl,
			})
		} else if (sf.Type.Kind() == reflect.Ptr || sf.Type.Kind() == reflect.Interface) && rv.Field(i).IsNil() {
			panic(tagMissingError{field: sf})
		} else if sf.Type.Kind() == reflect.Struct {
			// check forgotten tag only for struct.
			if _, exist := container.components[sf.Type]; exist {
				panic(tagForgottenError{field: sf})
			}
		}
	}

	if len(comp.dependencies) != 0 && !ptr {
		panic(incompletedError{})
	}

	container.components[rt] = append(container.components[rt], comp)
}

// Resolve a component with identified id.
func (container Container) Resolve(out interface{}, id ...string) {
	rv := reflect.ValueOf(out)

	if rv.Type().Kind() != reflect.Ptr {
		panic(resolveParamError{})
	}

	rv = rv.Elem()
	rt := rv.Type()

	nam := ""
	if len(id) > 0 {
		nam = id[0]
	}

	if rt.Kind() == reflect.Ptr {
		// pointer inside pointer
		if gr, ok := container.components[rt.Elem()]; ok {
			comp := gr.get(nam)
			if comp.value.CanAddr() {
				rv.Set(comp.value.Addr())
				return
			}

			panic(notAddressableError{id: nam, paramType: rt, component: comp})
		}
	} else {
		if gr, ok := container.components[rt]; ok {
			rv.Set(gr.get(nam).value)
			return
		}
	}

	panic(typeNotFoundError{paramType: rt})
}

// Apply wiring to all components.
func (container Container) Apply() {
	for _, gr := range container.components {
		for _, comp := range gr {
			container.fill(comp)
		}
	}
}

func (container Container) fill(c component) {
	if len(c.dependencies) == 0 {
		return
	}

	for i := range c.dependencies {
		dep := c.dependencies[i]
		cdep := component{}

		if gr, ok := container.components[dep.typ]; ok {
			cdep = gr.get(dep.id)
		} else {
			// scan if it's interface
			matches := 0

			if dep.typ.Kind() == reflect.Interface {
				for _, gr := range container.components {
					ctyp := gr[0].value.Type()

					if ctyp.Implements(dep.typ) && (dep.impl == "" || dep.impl == ctyp.Name()) {
						if fcedp, ok := gr.find(dep.id); ok {
							cdep = fcedp
							matches++
						}
					}
				}
			}

			if matches == 0 {
				panic(dependencyNotFound{id: dep.id, component: c, dependency: dep})
			} else if matches > 1 {
				panic(ambiguousError{component: c, dependency: dep})
			}
		}

		container.fill(cdep)

		fv := c.value.Field(dep.index)
		if fv.Kind() == reflect.Ptr {
			if !cdep.value.CanAddr() {
				panic(requiresPointerError{component: c, dependency: dep, depComponent: cdep})
			}

			fv.Set(cdep.value.Addr())
		} else {
			fv.Set(cdep.value)
		}
	}

	c.dependencies = nil
}
