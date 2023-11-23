# Maliketh Golang Implant

This directory contains the source code for the Golang implementation of the Maliketh implant.
Most functions are supported on Windows, macOS, and Linux.

Some functions were adapted from [Coldfire](https://github.com/redcode-labs/Coldfire).

## Differences between this and the C++ implant
The C++ implant is a bit more optimized for real world use. The golang implant **does not** have:

* Compiletime string obfuscation
* A small memory footprint
* A small binary size
* Scheduled task persistence


## TODO

* [] Register retries
* [] Persistence
* [] Do something "normal" when sandbox is detected


