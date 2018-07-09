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
	Value3 *int   `wire:"-"`
}

func (c ComponentD) Value() string {
	return c.Value1
}

type ComponentE struct {
	Value1 Valuer `wire:""`
}

func TestWire(t *testing.T) {
	componentA := ComponentA{Value1: "Hi!", Value2: 10}
	componentB := ComponentB{Value1: []int{1}, Value3: "Hello!"}
	componentC := ComponentC{Value3: false}
	componentD := ComponentD{}
	componentE := ComponentE{}

	app := wire.New()
	app.Connect("LGTM!")
	app.Connect(true)
	app.Connect([]int{1, 2, 3})
	app.Connect(&componentA)
	app.Connect(&componentB)
	app.Connect(&componentC)
	app.Connect(&componentD, "component_d")
	app.Connect(&componentE)
	app.Connect(&ComponentA{Value1: "Hello!", Value2: 20}, "hello")

	app.Apply()

	var resolvedString string
	var resolvedBool bool
	var resolvedSint []int
	var resolvedA ComponentA
	var resolvedB ComponentB
	var resolvedC ComponentC
	var resolvedD *ComponentD
	var resolvedE *ComponentE

	app.Resolve(&resolvedString)
	app.Resolve(&resolvedBool)
	app.Resolve(&resolvedSint)
	app.Resolve(&resolvedA)
	app.Resolve(&resolvedB)
	app.Resolve(&resolvedC)
	app.Resolve(&resolvedD, "component_d")
	app.Resolve(&resolvedE)

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
	assert.Equal(t, componentD, *resolvedD)
	assert.Equal(t, componentE, *resolvedE)
}
