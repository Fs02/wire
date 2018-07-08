# wire

[![GoDoc](https://godoc.org/github.com/Fs02/wire?status.svg)](https://godoc.org/github.com/Fs02/wire) [![Build Status](https://travis-ci.org/Fs02/wire.svg?branch=master)](https://travis-ci.org/Fs02/wire) [![Go Report Card](https://goreportcard.com/badge/github.com/Fs02/wire)](https://goreportcard.com/report/github.com/Fs02/wire) [![Maintainability](https://api.codeclimate.com/v1/badges/7957f2fe0d2c6fd5d72c/maintainability)](https://codeclimate.com/github/Fs02/wire/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/7957f2fe0d2c6fd5d72c/test_coverage)](https://codeclimate.com/github/Fs02/wire/test_coverage)

Wire is runtime depedency injection/wiring for golang. It's designed to be strict to avoid your go application running without proper dependency injected.

Features:

- Strictly validates dependency and prevents missing or ambiguous dependency.
- Easily connect and resolve object anywhere.
- Annotates ambiguous interface type using connection name or implementation name.

## Install

```bash
go get github.com/Fs02/wire
```

## License

Released under the [MIT License](https://github.com/Fs02/grimoire/blob/master/LICENSE)
