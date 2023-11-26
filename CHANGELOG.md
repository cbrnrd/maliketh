## v0.2.0

[`b85030c`](https://github.com/cbrnrd/maliketh/commit/b85030c) - Bump version to v0.2.1
[`eb6ca00`](https://github.com/cbrnrd/maliketh/commit/eb6ca00) - Fix create_release.sh permissions
[`86bb9d0`](https://github.com/cbrnrd/maliketh/commit/86bb9d0) - Add create_release.sh script
[`b3a87f5`](https://github.com/cbrnrd/maliketh/commit/b3a87f5) - Use poetry instead of requirements.txt
[`01da32c`](https://github.com/cbrnrd/maliketh/commit/01da32c) - Run formatter
[`22ad839`](https://github.com/cbrnrd/maliketh/commit/22ad839) - Reduce code duplication and remove some unused files/functions
[`b0e71a1`](https://github.com/cbrnrd/maliketh/commit/b0e71a1) - Documentation and makefile improvements

View diff [`v0.2.0...v0.2.1`](https://github.com/cbrnrd/maliketh/compare/v0.2.0...v0.2.1)

## v0.1.1

[`e8cf456`](https://github.com/cbrnrd/maliketh/commit/e8cf456) - Update changelog and readme
[`0639f87`](https://github.com/cbrnrd/maliketh/commit/0639f87) - Add go implant
[`0aeec4c`](https://github.com/cbrnrd/maliketh/commit/0aeec4c) - Add Defender disable command
[`546054f`](https://github.com/cbrnrd/maliketh/commit/546054f) - Clean up docs and files
[`cc63766`](https://github.com/cbrnrd/maliketh/commit/cc63766) - Update build instructions in README
[`a130207`](https://github.com/cbrnrd/maliketh/commit/a130207) - Update README and CHANGELOG

View diff [`v0.1.1...v0.2.0`](https://github.com/cbrnrd/maliketh/compare/v0.1.1...v0.2.0)


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
