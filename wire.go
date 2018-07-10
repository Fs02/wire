// Package wire is runtime depedency injection/wiring for golang.
// It's designed to be strict to avoid your go application running without proper dependency injected.
package wire

var global Container

func init() {
	global = New()
	global.callerSkip = 1
}

// Connect a component, optionally identified by name.
//
// This will panic if:
//   1. Duplicate component found.
//   2. Possible forgotten `wire` tag found on pointer and interface field.
//   3. Component that need wiring is passed using value.
func Connect(val interface{}, name ...string) {
	global.Connect(val, name...)
}

// Resolve a component optionally identified by name.
//
// This should be called only after wiring applied.
// Resolving component multiple times should be avoided, consider caching the component if you need.
// For example, if you are running a web server, Resolve should only be done before the server start listening for request.
//
// This will panic if no matching component is found.
func Resolve(out interface{}, name ...string) {
	global.Resolve(out, name...)
}

// Apply wiring to all components.
//
// This will panic if:
//   1. There are missing component.
//   2. Ambiguous field found, usually field with type interface that can satisfy more than one component.
func Apply() {
	global.Apply()
}
