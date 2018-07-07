package wire_test

import (
	"testing"

	"github.com/Fs02/wire"
	"github.com/stretchr/testify/assert"
)

type Valuer interface {
	Value() string
}

type ComponentA struct {
	Value1 string
	Value2 int
}

func (c ComponentA) Value() string {
	return c.Value1
}

type ComponentB struct {
	Value1 []int
	Value2 ComponentA `wire:"hello"`
	Value3 string
	Value4 bool `wire:""`
}

type ComponentC struct {
	Value1 *ComponentA `wire:""`
	Value2 *ComponentB `wire:""`
	Value3 bool
	Value4 []int  `wire:""`
	Value5 Valuer `wire:"component_d,ComponentD"`
}

type ComponentD struct {
	Value1 string `wire:",ComponentA"`
}

func (c ComponentD) Value() string {
	return c.Value1
}

type ComponentE struct {
	Value1 Valuer `wire:""`
}

func TestWire(t *testing.T) {
	wire.Clear()

	componentA := ComponentA{Value1: "Hi!", Value2: 10}
	componentB := ComponentB{Value1: []int{1}, Value3: "Hello!"}
	componentC := ComponentC{Value3: false}
	componentD := ComponentD{}
	componentE := ComponentE{}

	wire.Connect("LGTM!")
	wire.Connect(true)
	wire.Connect([]int{1, 2, 3})
	wire.Connect(&componentA)
	wire.Connect(&componentB)
	wire.Connect(&componentC)
	wire.Connect(&componentD, "component_d")
	wire.Connect(&componentE)
	wire.Connect(&ComponentA{Value1: "Hello!", Value2: 20}, "hello")

	wire.Apply()

	var resolvedString string
	var resolvedBool bool
	var resolvedSint []int
	var resolvedA ComponentA
	var resolvedB ComponentB
	var resolvedC ComponentC
	var resolvedD ComponentD
	var resolvedE ComponentE

	wire.Resolve(&resolvedString)
	wire.Resolve(&resolvedBool)
	wire.Resolve(&resolvedSint)
	wire.Resolve(&resolvedA)
	wire.Resolve(&resolvedB)
	wire.Resolve(&resolvedC)
	wire.Resolve(&resolvedD, "component_d")
	wire.Resolve(&resolvedE)

	assert.Equal(t, ComponentA{Value1: "Hi!", Value2: 10}, componentA)

	assert.Equal(t, ComponentB{
		Value1: []int{1},
		Value2: ComponentA{Value1: "Hello!", Value2: 20},
		Value3: "Hello!",
		Value4: true,
	}, componentB)

	assert.Equal(t, ComponentC{
		Value1: &ComponentA{Value1: "Hi!", Value2: 10},
		Value2: &ComponentB{
			Value1: []int{1},
			Value2: ComponentA{Value1: "Hello!", Value2: 20},
			Value3: "Hello!",
			Value4: true,
		},
		Value3: false,
		Value4: []int{1, 2, 3},
		Value5: ComponentD{Value1: "LGTM!"},
	}, componentC)

	assert.Equal(t, ComponentD{Value1: "LGTM!"}, componentD)

	assert.Equal(t, ComponentE{Value1: ComponentA{Value1: "Hi!", Value2: 10}}, componentE)

	assert.Equal(t, "LGTM!", resolvedString)
	assert.Equal(t, true, resolvedBool)
	assert.Equal(t, []int{1, 2, 3}, resolvedSint)
	assert.Equal(t, componentA, resolvedA)
	assert.Equal(t, componentB, resolvedB)
	assert.Equal(t, componentC, resolvedC)
	assert.Equal(t, componentD, resolvedD)
	assert.Equal(t, componentE, resolvedE)
}

func TestWire_ambiguousConnection(t *testing.T) {
	wire.Clear()

	wire.Connect("LGTM!")
	wire.Connect(&ComponentA{Value1: "Hi!", Value2: 10})
	wire.Connect(&ComponentD{})
	wire.Connect(&ComponentE{})

	assert.Panics(t, func() {
		wire.Apply()
	})
}

func TestWire_requiresReferenceInsteadOfValue(t *testing.T) {
	vstring := "LGTM!"
	vbool := true
	componentB := ComponentB{Value1: []int{1}, Value3: "Hello!"}
	componentC := ComponentC{Value3: false}

	wire.Clear()

	wire.Connect(vstring)
	wire.Connect(vbool)
	wire.Connect(ComponentA{Value1: "Hi!", Value2: 10})
	wire.Connect(ComponentA{Value1: "Hello!", Value2: 10}, "hello")
	wire.Connect(&componentB)
	wire.Connect(&componentC)

	assert.Panics(t, func() {
		wire.Apply()
	})
}

func TestWire_missingDependency(t *testing.T) {
	componentD := ComponentD{}

	wire.Clear()
	wire.Connect(&componentD)

	assert.Panics(t, func() {
		wire.Apply()
	})
}

func TestWire_duplicateDependency(t *testing.T) {
	componentD := ComponentD{}

	wire.Clear()
	wire.Connect(&componentD)

	assert.Panics(t, func() {
		wire.Connect(&componentD)
	})
}

func TestWire_cannotWireComponent(t *testing.T) {
	componentD := ComponentD{}

	wire.Clear()

	assert.Panics(t, func() {
		wire.Connect(componentD)
	})
}

func TestWire_Resolve_mustPointer(t *testing.T) {
	assert.Panics(t, func() {
		wire.Resolve(ComponentA{})
	})
}

func TestWire_Resolve_typeNotFound(t *testing.T) {
	wire.Clear()

	assert.Panics(t, func() {
		wire.Resolve(&ComponentA{}, "notexist")
	})
}

func TestWire_Resolve_nameNotFound(t *testing.T) {
	componentA := ComponentA{}

	wire.Clear()
	wire.Connect(&componentA)

	assert.Panics(t, func() {
		wire.Resolve(&ComponentA{}, "notexist")
	})
}
