# go-udev
Go bindings for libudev

## Documentation
Documentation is on Godoc.
[![GoDoc](https://godoc.org/github.com/jochenvg/go-udev?status.svg)](https://godoc.org/github.com/jochenvg/go-udev)

## Build
No cgo required: libudev (`libudev.so.1`) is loaded at runtime with
[purego](https://github.com/ebitengine/purego), so `CGO_ENABLED=0` builds work
and `libudev-dev` is not needed. If libudev is missing on the host, the `Udev`
constructors return nil instead of the process failing to start.
