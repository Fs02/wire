package wire_test

import (
	"reflect"
	"testing"

	"github.com/Fs02/wire"
	"github.com/stretchr/testify/assert"
)

type ComponentA struct {
	Value1 string
	Value2 int
}

type ComponentB struct {
	Value1 []int
	Value2 ComponentA `wire:""`
	Value3 string
	Value4 bool `wire:""`
}

type ComponentC struct {
	Value1 *ComponentA `wire:""`
	Value2 *ComponentB `wire:""`
	Value3 bool
	Value4 []int      `wire:""`
	Value5 ComponentD `wire:""`
}

type ComponentD struct {
	Value1 string `wire:""`
}

func TestWire(t *testing.T) {
	vstring := "LGTM!"
	vbool := true
	vsint := []int{1, 2, 3}
	componentA := ComponentA{Value1: "Hi!", Value2: 10}
	componentB := ComponentB{Value1: []int{1}, Value3: "Hello!"}
	componentC := ComponentC{Value3: false}
	componentD := ComponentD{}

	wire.Init()

	wire.Connect(vstring)
	wire.Connect(vbool)
	wire.Connect(vsint)
	wire.Connect(&componentA)
	wire.Connect(&componentB)
	wire.Connect(&componentC)
	wire.Connect(&componentD)

	cloneA := wire.Get(reflect.TypeOf(ComponentA{})).(*ComponentA)
	cloneB := wire.Get(reflect.TypeOf(ComponentB{})).(*ComponentB)
	cloneC := wire.Get(reflect.TypeOf(&ComponentC{})).(*ComponentC)
	cloneD := wire.Get(reflect.TypeOf(&ComponentD{})).(*ComponentD)

	wire.Apply()

	assert.Equal(t, "LGTM!", vstring)
	assert.Equal(t, true, vbool)
	assert.Equal(t, []int{1, 2, 3}, vsint)

	assert.Equal(t, ComponentA{Value1: "Hi!", Value2: 10}, componentA)

	assert.Equal(t, ComponentB{
		Value1: []int{1},
		Value2: ComponentA{Value1: "Hi!", Value2: 10},
		Value3: "Hello!",
		Value4: true,
	}, componentB)

	assert.Equal(t, ComponentC{
		Value1: &ComponentA{Value1: "Hi!", Value2: 10},
		Value2: &ComponentB{
			Value1: []int{1},
			Value2: ComponentA{Value1: "Hi!", Value2: 10},
			Value3: "Hello!",
			Value4: true,
		},
		Value3: false,
		Value4: []int{1, 2, 3},
		Value5: ComponentD{Value1: "LGTM!"},
	}, componentC)

	assert.Equal(t, ComponentD{Value1: "LGTM!"}, componentD)

	assert.Equal(t, componentA, *cloneA)
	assert.Equal(t, componentB, *cloneB)
	assert.Equal(t, componentC, *cloneC)
	assert.Equal(t, componentD, *cloneD)
}

func TestWire_requiresReferenceInsteadOfValue(t *testing.T) {
	vstring := "LGTM!"
	vbool := true
	vsint := []int{1, 2, 3}
	componentA := ComponentA{Value1: "Hi!", Value2: 10}
	componentB := ComponentB{Value1: []int{1}, Value3: "Hello!"}
	componentC := ComponentC{Value3: false}
	componentD := ComponentD{}

	wire.Init()

	wire.Connect(vstring)
	wire.Connect(vbool)
	wire.Connect(vsint)
	wire.Connect(componentA)
	wire.Connect(&componentB)
	wire.Connect(&componentC)
	wire.Connect(&componentD)

	assert.Panics(t, func() {
		wire.Apply()
	})
}

func TestWire_missingDependency(t *testing.T) {
	componentD := ComponentD{}

	wire.Init()
	wire.Connect(&componentD)

	assert.Panics(t, func() {
		wire.Apply()
	})
}

func TestWire_cannotWireComponent(t *testing.T) {
	componentD := ComponentD{}

	wire.Init()

	assert.Panics(t, func() {
		wire.Connect(componentD)
	})
}

func TestWire_getUnaddressableComponent(t *testing.T) {
	componentA := ComponentA{}

	wire.Init()

	wire.Connect(componentA)

	assert.Panics(t, func() {
		wire.Get(reflect.TypeOf(ComponentA{}))
	})
}
