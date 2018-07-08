package wire

import (
	"reflect"
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

var global Container

func init() {
	global = New()
}

// Connect a component, optionally identified by name.
func Connect(val interface{}, name ...string) {
	global.Connect(val, name...)
}

// Resolve a component with identified name.
// This should be called only after wiring applied.
// Resolving component multiple times should be avoided, consider caching the component if you need.
func Resolve(out interface{}, name ...string) {
	global.Resolve(out, name...)
}

// Apply wiring to all components.
func Apply() {
	global.Apply()
}
