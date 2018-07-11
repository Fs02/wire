package wire

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getComponent() component {
	return component{
		value:      reflect.ValueOf(0),
		declaredAt: "/somefile.go:1",
	}
}

func getDependency() dependency {
	return dependency{
		id:   "a",
		name: "A",
		typ:  reflect.TypeOf(0),
	}
}

func TestIdNotFoundError(t *testing.T) {
	assert.Equal(t, "wire: no int identified using \"a\" found",
		idNotFoundError{id: "a", component: getComponent()}.Error())
}

func TestDuplicateError(t *testing.T) {
	assert.Equal(t, "wire: trying to connect component with same type and id. previosly declared here:\n\t/somefile.go:1",
		duplicateError{previous: getComponent()}.Error())
}

func TestTagMissingError(t *testing.T) {
	assert.Equal(t, "wire: field with nil interface or pointer without wire detected for A with type int, perhaps you forgot? to ignore add `wire:\"-\"`",
		tagMissingError{field: reflect.StructField{Name: "A", Type: reflect.TypeOf(0)}}.Error())
}

func TestIncompletedError(t *testing.T) {
	assert.Equal(t, "wire: trying to connect incompleted component as a value, use a reference instead",
		incompletedError{}.Error())
}

func TestResolveParamError(t *testing.T) {
	assert.Equal(t, "wire: resolve parameter must be a pointer",
		resolveParamError{}.Error())
}

func TestNotAddressableError(t *testing.T) {
	assert.Equal(t, "wire: component with type int identified by \"a\" is not addressable, connect component using reference instead of value. declared here:\n\t/somefile.go:1",
		notAddressableError{id: "a", paramType: reflect.TypeOf(0), component: getComponent()}.Error())
}

func TestTypeNotFoundError(t *testing.T) {
	assert.Equal(t, "wire: no component with type int found",
		typeNotFoundError{paramType: reflect.TypeOf(0)}.Error())
}

func TestDependencyNotFoundError(t *testing.T) {
	assert.Equal(t, "wire: field A of int requires int identified using \"a\", but none was found. declared here:\n\t/somefile.go:1",
		dependencyNotFound{id: "a", component: getComponent(), dependency: getDependency()}.Error())
}

func TestAmbiguousError(t *testing.T) {
	assert.Equal(t, "wire: ambiguous connection found on field A of int, multiple components satisfy int interface, consider using id. declared here:\n\t/somefile.go:1",
		ambiguousError{component: getComponent(), dependency: getDependency()}.Error())
}

func TestRequiresPointerError(t *testing.T) {
	assert.Equal(t, "wire: field A of int requires int as pointer, connect int as a reference instead of a value. declared here:\n\t/somefile.go:1\n\t/somefile.go:1",
		requiresPointerError{component: getComponent(), dependency: getDependency(), depComponent: getComponent()}.Error())
}
