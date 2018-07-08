package wire

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
