# Maliketh Golang Implant

This directory contains the source code for the Golang implementation of the Maliketh implant.
Most functions are supported on Windows, macOS, and Linux.

Some functions were adapted from [Coldfire](https://github.com/redcode-labs/Coldfire).

## Building

To build the development implant, simply run

```bash
$ make
```

This builds the debug version of the implant for your native OS and architecture.

To build for all supported OSes and architectures, run

```bash
$ make all
```

To build a stripped version of the implant, simply set `DEBUG=0` when running make. Example:

```bash
$ DEBUG=0 make all
```

You can also set the usual `GOOS` and `GOARCH` environment variables to build for a specific OS and architecture using

```bash
$ make build
```

## Differences between this and the C++ implant
The C++ implant is a bit more optimized for real world use. The golang implant **does not** have:

* Compiletime string obfuscation
* A small memory footprint
* A small binary size
* Scheduled task persistence


## TODO

* [ ] Register retries
* [ ] Persistence
* [ ] Do something "normal" when sandbox is detected
* [ ] Auto self destruct
* [ ] Jitter


