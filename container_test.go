package wire_test

import (
	"testing"

	"github.com/Fs02/wire"
	"github.com/stretchr/testify/assert"
)

func TestContainer_Apply_ambiguousConnection(t *testing.T) {
	app := wire.New()

	app.Connect("LGTM!")
	app.Connect(&ComponentA{Value1: "Hi!", Value2: 10})
	app.Connect(&ComponentD{})
	app.Connect(&ComponentE{})

	assert.Panics(t, func() {
		app.Apply()
	})
}

func TestContainer_Apply_requiresPointerInsteadOfValue(t *testing.T) {
	vstring := "LGTM!"
	vbool := true
	componentB := ComponentB{Value1: []int{1}, Value3: "Hello!"}
	componentC := ComponentC{Value3: false}

	app := wire.New()

	app.Connect(vstring)
	app.Connect(vbool)
	app.Connect(ComponentA{Value1: "Hi!", Value2: 10})
	app.Connect(ComponentA{Value1: "Hello!", Value2: 10}, "hello")
	app.Connect(&componentB)
	app.Connect(&componentC)

	assert.Panics(t, func() {
		app.Apply()
	})
}

func TestContainer_Apply_missingDependency(t *testing.T) {
	componentD := ComponentD{}

	app := wire.New()
	app.Connect(&componentD)

	assert.Panics(t, func() {
		app.Apply()
	})
}

func TestContainer_Connect_duplicateDependency(t *testing.T) {
	componentD := ComponentD{}

	app := wire.New()
	app.Connect(&componentD)

	assert.Panics(t, func() {
		app.Connect(&componentD)
	})
}

func TestContainer_Connect_cannotWireComponent(t *testing.T) {
	componentD := ComponentD{}

	app := wire.New()

	assert.Panics(t, func() {
		app.Connect(componentD)
	})
}

func TestContainer_Connect_nilInterface(t *testing.T) {
	var a struct {
		Value Valuer
	}

	app := wire.New()

	assert.Panics(t, func() {
		app.Connect(&a)
	})
}

func TestContainer_Connect_nilPointer(t *testing.T) {
	var a struct {
		Value *int
	}

	app := wire.New()

	assert.Panics(t, func() {
		app.Connect(&a)
	})
}

func TestContainer_Connect_tagForgotten(t *testing.T) {
	var a struct {
		Value ComponentA
	}

	app := wire.New()
	app.Connect(ComponentA{})

	assert.Panics(t, func() {
		app.Connect(&a)
	})
}

func TestContainer_Resolve_mustPointer(t *testing.T) {
	app := wire.New()

	assert.Panics(t, func() {
		app.Resolve(ComponentA{})
	})
}

func TestContainer_Resolve_typeNotFound(t *testing.T) {
	app := wire.New()

	assert.Panics(t, func() {
		app.Resolve(&ComponentA{}, "notexist")
	})
}

func TestContainer_Resolve_nameNotFound(t *testing.T) {
	componentA := ComponentA{}

	app := wire.New()
	app.Connect(&componentA)
	app.Apply()

	assert.Panics(t, func() {
		app.Resolve(&ComponentA{}, "notexist")
	})
}

func TestContainer_Resolve_valueAsPointer(t *testing.T) {
	app := wire.New()
	app.Connect(ComponentA{})
	app.Apply()

	var resolve *ComponentA

	assert.Panics(t, func() {
		app.Resolve(&resolve)
	})
}
