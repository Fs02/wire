# wire

[![GoDoc](https://godoc.org/github.com/Fs02/wire?status.svg)](https://godoc.org/github.com/Fs02/wire) [![Build Status](https://travis-ci.org/Fs02/wire.svg?branch=master)](https://travis-ci.org/Fs02/wire) [![Go Report Card](https://goreportcard.com/badge/github.com/Fs02/wire)](https://goreportcard.com/report/github.com/Fs02/wire) [![Maintainability](https://api.codeclimate.com/v1/badges/7957f2fe0d2c6fd5d72c/maintainability)](https://codeclimate.com/github/Fs02/wire/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/7957f2fe0d2c6fd5d72c/test_coverage)](https://codeclimate.com/github/Fs02/wire/test_coverage)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FFs02%2Fwire.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FFs02%2Fwire?ref=badge_shield)

Wire is runtime depedency injection/wiring for golang. It's designed to be strict to avoid your go application running without proper dependency injected.

Features:

- Strictly validates dependency and prevents missing or ambiguous dependency.
- Check againts possible forgotten `wire` tag.
- Easily connect and resolve object anywhere.
- Annotates ambiguous interface type using connection name or implementation name.

## Install

```bash
go get github.com/Fs02/wire
```

## Example

```golang
package wire_test

import (
	"fmt"

	"github.com/Fs02/wire"
)

type Listener struct{}

func (listener Listener) Next() string {
	return "system"
}

type Printer interface {
	Exec(string) error
}

type SystemPrint struct {
	App string `wire:""`
}

func (systemPrint SystemPrint) Exec(msg string) error {
	fmt.Println("[" + systemPrint.App + "] System: " + msg)
	return nil
}

type UserPrint struct {
	App    string `wire:""`
	Target string
}

func (userPrint UserPrint) Exec(msg string) error {
	fmt.Println("[" + userPrint.App + "]" + userPrint.Target + ": " + msg)
	return nil
}

type Service struct {
	// Each of `wire`` tag below indicate fields to be wired with apporpriate component.
	// value inside `wire` tag indicate the name of the component and optionally it's type.
	// `wire` with empty value will be wired with default value (named using empty string).
	// Ambiguous field can be resolved by adding it's type, name or both (separated using comma) to the `wire` tag.
	// Don't worry if you forgot to add wire tag to an interface or a pointer, wire will warn you if any nil field are found.
	// To ignore wiring on specific field, you can use add `wire:"-"`.
	Listener     Listener `wire:""`
	SystemPrint  Printer  `wire:",SystemPrint"`
	FooUserPrint Printer  `wire:"foo"`
	BooUserPrint Printer  `wire:"boo,UserPrint"`
}

func (service Service) Update() error {
	switch service.Listener.Next() {
	case "system":
		return service.SystemPrint.Exec("hello from system")
	case "user-foo":
		return service.FooUserPrint.Exec("hello from foo")
	case "user-boo":
		return service.BooUserPrint.Exec("hello from boo")
	default:
		return nil
	}
}

func init() {
	// add components to be wired by the library.
	// wire all components only once and as early as possible.
	wire.Connect("CoolApp")
	wire.Connect(Listener{})                       // we don't need to pass by reference here, since it doesn't require any wiring.
	wire.Connect(&SystemPrint{})                   // we need to pass by reference it to allow wiring, wire will panic if we pass by value.
	wire.Connect(&UserPrint{Target: "foo"}, "foo") // wire a UserPrint named by "foo".
	wire.Connect(&UserPrint{Target: "boo"}, "boo") // wire a UserPrint named by "boo", wire will panic if there's duplicate components detected.
	wire.Connect(&Service{})

	// Apply wiring
	wire.Apply()
}

func Example() {
	// Resolve a service component to be used later.
	var service Service
	wire.Resolve(&service)

	service.Update()
	// Output: [CoolApp] System: hello from system
}
```

## License

Released under the [MIT License](https://github.com/Fs02/wire/blob/master/LICENSE)


[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FFs02%2Fwire.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FFs02%2Fwire?ref=badge_large)