# Changelog
All notable changes to this project will be documented in this file.

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