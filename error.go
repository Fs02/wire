package wire

import (
	"reflect"
)

type idNotFoundError struct {
	id        string
	component component
}

func (err idNotFoundError) Error() string {
	return "wire: no " + err.component.value.Type().String() +
		" identified using \"" + err.id + "\" found"
}

type duplicateError struct {
	previous component
}

func (err duplicateError) Error() string {
	return "wire: trying to connect component with same type and id. previosly declared here:\n\t" +
		err.previous.declaredAt
}

type tagMissingError struct {
	field reflect.StructField
}

func (err tagMissingError) Error() string {
	return "wire: field with nil interface or pointer without wire detected for " +
		err.field.Name + " with type " + err.field.Type.String() +
		", perhaps you forgot? to ignore add `wire:\"-\"`"
}

type tagForgottenError struct {
	field reflect.StructField
}

func (err tagForgottenError) Error() string {
	return "wire: tag is missing for already connected type on " + err.field.Name + " with type " + err.field.Type.String() +
		", perhaps you forgot? to ignore add `wire:\"-\"`"
}

type incompletedError struct{}

func (err incompletedError) Error() string {
	return "wire: trying to connect incompleted component as a value, use a reference instead"
}

type resolveParamError struct{}

func (err resolveParamError) Error() string {
	return "wire: resolve parameter must be a pointer"
}

type notAddressableError struct {
	id        string
	paramType reflect.Type
	component component
}

func (err notAddressableError) Error() string {
	return "wire: component with type " + err.paramType.String() + " identified by \"" + err.id +
		"\" is not addressable, connect component using reference instead of value. declared here:\n\t" +
		err.component.declaredAt
}

type typeNotFoundError struct {
	paramType reflect.Type
}

func (err typeNotFoundError) Error() string {
	return "wire: no component with type " + err.paramType.String() + " found"
}

type dependencyNotFound struct {
	id         string
	component  component
	dependency dependency
}

func (err dependencyNotFound) Error() string {
	return "wire: field " + err.dependency.name + " of " + err.component.value.Type().String() +
		" requires " + err.dependency.typ.String() + " identified using \"" + err.id +
		"\", but none was found. declared here:\n\t" + err.component.declaredAt
}

type ambiguousError struct {
	component  component
	dependency dependency
}

func (err ambiguousError) Error() string {
	return "wire: ambiguous connection found on field " + err.dependency.name + " of " +
		err.component.value.Type().String() + ", multiple components satisfy " + err.dependency.typ.String() +
		" interface, consider using id. declared here:\n\t" + err.component.declaredAt
}

type requiresPointerError struct {
	component    component
	dependency   dependency
	depComponent component
}

func (err requiresPointerError) Error() string {
	return "wire: field " + err.dependency.name + " of " + err.component.value.Type().String() +
		" requires " + err.dependency.typ.String() + " as pointer, connect " + err.dependency.typ.String() +
		" as a reference instead of a value. declared here:\n\t" + err.component.declaredAt +
		"\n\t" + err.depComponent.declaredAt
}
