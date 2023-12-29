# Maliketh Golang Implant

This directory contains the source code for the Golang implementation of the Maliketh implant.
Most functions are supported on Windows, macOS, and Linux.

Some functions were adapted from [Coldfire](https://github.com/redcode-labs/Coldfire).

## Building

There are a few values you need to change before building the implant.

In `pkg/config/config.go` you need to change the following values to match your server configuration:
```go
// !!!!!!!!!! CHANGE THESE !!!!!!!!! //
const C2_DOMAIN = "localhost"
const C2_PORT = 80
const C2_USE_TLS = false
const C2_REGISTER_PASSWORD = "SWh5bHhGOENYQWF1TW9KR3VTb0YwVkVWbDRud1RFaHc="
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! //
```

To build the development implant, simply run

```bash
$ make
```

This builds the debug version of the implant for your native OS and architecture.

To build for Mac, Linux, and Windows, run:

```bash
$ make most
```

This will cover most if the use cases.

To build for *all* supported OSes and architectures, run

```bash
$ make all
```

To build a stripped version of the implant, simply set `DEBUG=0` when running make. Example:

```bash
$ DEBUG=0 make most
```

You can also set the usual `GOOS` and `GOARCH` environment variables to build for a specific OS and architecture using

```bash
$ GOOS=linux GOARCH=arm make build
```

## Differences between this and the C++ implant
The C++ implant is a bit more optimized for real world use. The golang implant **does not** have:

* Compiletime string obfuscation
* A small memory footprint
* A small binary size
* Scheduled task persistence


## TODO

* [x] Register retries
* [ ] Persistence
* [ ] Do something "normal" when sandbox is detected
* [x] Auto self destruct
* [x] Jitter
* [ ] Builder sript


