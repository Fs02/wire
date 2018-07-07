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

	vstring := wire.Connect("LGTM!").(string)
	vbool := wire.Connect(true).(bool)
	vsint := wire.Connect([]int{1, 2, 3}).([]int)
	componentA := wire.Connect(&ComponentA{Value1: "Hi!", Value2: 10}).(*ComponentA)
	componentB := wire.Connect(&ComponentB{Value1: []int{1}, Value3: "Hello!"}).(*ComponentB)
	componentC := wire.Connect(&ComponentC{Value3: false}).(*ComponentC)
	componentD := wire.Connect(&ComponentD{}, "component_d").(*ComponentD)
	componentE := wire.Connect(&ComponentE{}).(*ComponentE)
	wire.Connect(&ComponentA{Value1: "Hello!", Value2: 20}, "hello")

	wire.Apply()

	cloneA := ComponentA{}
	cloneB := ComponentB{}
	cloneC := ComponentC{}
	cloneD := ComponentD{}
	cloneE := ComponentE{}

	wire.Resolve(&cloneA)
	wire.Resolve(&cloneB)
	wire.Resolve(&cloneC)
	wire.Resolve(&cloneD, "component_d")
	wire.Resolve(&cloneE)

	assert.Equal(t, "LGTM!", vstring)
	assert.Equal(t, true, vbool)
	assert.Equal(t, []int{1, 2, 3}, vsint)

	assert.Equal(t, &ComponentA{Value1: "Hi!", Value2: 10}, componentA)

	assert.Equal(t, &ComponentB{
		Value1: []int{1},
		Value2: ComponentA{Value1: "Hello!", Value2: 20},
		Value3: "Hello!",
		Value4: true,
	}, componentB)

	assert.Equal(t, &ComponentC{
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

	assert.Equal(t, &ComponentD{Value1: "LGTM!"}, componentD)

	assert.Equal(t, &ComponentE{Value1: ComponentA{Value1: "Hi!", Value2: 10}}, componentE)

	assert.Equal(t, *componentA, cloneA)
	assert.Equal(t, *componentB, cloneB)
	assert.Equal(t, *componentC, cloneC)
	assert.Equal(t, *componentD, cloneD)
	assert.Equal(t, *componentE, cloneE)
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
