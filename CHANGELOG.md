# Changelog

All notable changes to this project will be documented in this file.

## 11/23/2023

### Added

* Golang implant! ([0639f87](https://github.com/cbrnrd/maliketh/commit/0639f8797838469a068d91f095e3307d2d73ecc4))
* Add AV disable ([0aeec4c](https://github.com/cbrnrd/maliketh/commit/0aeec4c4be8f1efaeaf15ee3d289507036c691df))

### Fixed

* Fix a few minor bugs with the client ([0639f87](https://github.com/cbrnrd/maliketh/commit/0639f8797838469a068d91f095e3307d2d73ecc4))
* Fix `chdir` command not working in client ([0639f87](https://github.com/cbrnrd/maliketh/commit/0639f8797838469a068d91f095e3307d2d73ecc4))

## 5/21/2023

### Added

* `build` command to frontend ([b116ad3](https://github.com/cbrnrd/maliketh/commit/b116ad390120eee94d87ec313d1728a2f2244289))
* Functionality to build the implant from the server ([917d514](https://github.com/cbrnrd/maliketh/commit/917d514fc6075cc15d0e45b4a1a546e6217e4139))

### Changed

* Tweaked the operator Dockerfile to include the implant source code ([917d514](https://github.com/cbrnrd/maliketh/commit/917d514fc6075cc15d0e45b4a1a546e6217e4139))

### Fixed

* Fix Github Action for building dockerfiles to reflect changes above ([894e14d](https://github.com/cbrnrd/maliketh/commit/894e14d79e4b826e6fcdece988f21d35daf09dd6))
* Link implant to `psapi` ([917d514](https://github.com/cbrnrd/maliketh/commit/917d514fc6075cc15d0e45b4a1a546e6217e4139))
* Remove `Lmcons.h` include ([917d514](https://github.com/cbrnrd/maliketh/commit/917d514fc6075cc15d0e45b4a1a546e6217e4139))

## 5/8/2023

### Added

* Basic OS functions built in ([91a40f2](https://github.com/cbrnrd/maliketh/commit/91a40f2ba1cded5a025004a6143578fa84baec66))
  * `cd` - Change directory
  * `ls` - List directory contents
  * `pwd` - Print working directory
  * `getenv` - Get all environment variables
  * `ps` - List running processes
  * `whoami` - Get current username

### Fixed

* Fixed a bug when decoding binary task results ([91a40f2](https://github.com/cbrnrd/maliketh/commit/91a40f2ba1cded5a025004a6143578fa84baec66))
  * This was causing reading downloaded files to fail because the file was being decoded as a UTF-8 string instead of binary
* Fixed a logic bug when SELFDESTRUCT tasks are read ([91a40f2](https://github.com/cbrnrd/maliketh/commit/91a40f2ba1cded5a025004a6143578fa84baec66))
  * The C2 listener was returning before removing the implant from the database, resulting in stale implants
* Fixed a bug in the implant where the initial `FreeConsole` and `Sleep` calls were never running ([91a40f2](https://github.com/cbrnrd/maliketh/commit/91a40f2ba1cded5a025004a6143578fa84baec66))
