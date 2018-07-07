# wire

[![GoDoc](https://godoc.org/github.com/Fs02/wire?status.svg)](https://godoc.org/github.com/Fs02/wire) [![Build Status](https://travis-ci.org/Fs02/wire.svg?branch=master)](https://travis-ci.org/Fs02/wire) [![Go Report Card](https://goreportcard.com/badge/github.com/Fs02/wire)](https://goreportcard.com/report/github.com/Fs02/wire)

Wire is runtime depedency injection/wiring for golang. It's designed to be strict to avoid your go application running without proper dependency injected.

Features:

- Easily connect and resolve object anywhere.
- Strictly validates dependency and prevents missing or ambiguous dependency.
- Annotates ambiguous interface type using connection name or implementation name.

## Install

```bash
go get github.com/Fs02/wire
```

## License

Released under the [MIT License](https://github.com/Fs02/grimoire/blob/master/LICENSE)
