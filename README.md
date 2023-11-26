<!-- Centered logo -->
<p align="center">
  <img src="./data/maliketh_logo.png" alt="Maliketh logo" width="900" height="300"/>
</p>

<p align="center">
  A multi-user, customizable C2 framework.
  <br>
  <br>
  <img alt="GitHub License" src="https://img.shields.io/github/license/cbrnrd/maliketh?style=flat-square&color=ff5733">
  <img alt="Latest Release" src="https://img.shields.io/github/v/release/cbrnrd/maliketh?logo=github&style=flat-square&color=ff5733&link=https%3A%2F%2Fgithub.com%2Fcbrnrd%2Fmaliketh%2Freleases">
  <img alt="GitHub Workflow Status (with event)" src="https://img.shields.io/github/actions/workflow/status/cbrnrd/maliketh/docker-image.yml?style=flat-square">
  <img alt="GitHub Release Date - Published_At" src="https://img.shields.io/github/release-date/cbrnrd/maliketh?style=flat-square&color=ff5733&link=https%3A%2F%2Fgithub.com%2Fcbrnrd%2Fmaliketh%2Freleases">
  <img alt="GitHub contributors" src="https://img.shields.io/github/contributors/cbrnrd/maliketh?style=flat-square&color=ff5733">
</p>

---

The goal of Maliketh is to provide a flexible, easy to use C2 framework that can be customized to fit the needs of the operator. The poster used in the initial presentation is located [here](./data/Maliketh%20C2%20Poster.png).

## Implant features

The initial implant was written in C++ and targeted for Windows, but a Golang implant has also been implemented and supports all major platforms.

The main feature of the implant is its ability to change its behavior based on the configuration file it receives from the server. This allows the operator to customize the implant to fit their needs. The implant also has the following features (see [here](./design/opcodes.md) for more info):

* File upload/download
* Command execution
* Shellcode injection
* Update configuration
* Send system information
* Self-destruct
* Sleep
* Basic Anti-debugging
* *Very* Basic Anti-VM
* Sleep skipping detection

## Future work

- [x] Implement Golang client ([0639f87](https://github.com/cbrnrd/maliketh/commit/0639f8797838469a068d91f095e3307d2d73ecc4))
* [x] Per-operator builder in-server ([917d514](https://github.com/cbrnrd/maliketh/commit/917d514fc6075cc15d0e45b4a1a546e6217e4139))
* [ ] Stealer/basic looter
* [x] AV Disable ([0aeec4c](https://github.com/cbrnrd/maliketh/commit/0aeec4c4be8f1efaeaf15ee3d289507036c691df))
* [ ] Change design of config to be protocol agnostic.
  * ie Define an HTTPS layer/adapter and separate out the code better.
* [ ] Keylogger
* [ ] Allow implant aliasing/renaming
  * This shouldn't change the actual ID, just create a mapping table
* [ ] More fine grained backend roles and actions (blocking users, % bot allocation)
* [ ] Add ability to send command to every bot
* [ ] Floods
* [ ] Route RabbitMQ traffic through Admin listener instead of directly connecting
* [ ] Improved anti-vm (check BIOS information)
  * [x] Not bad in golang implant
* [x] More stable file uploads/downloads ([91a40f2](https://github.com/cbrnrd/maliketh/commit/91a40f2ba1cded5a025004a6143578fa84baec66))
* [x] Basic OS functions built in ([91a40f2](https://github.com/cbrnrd/maliketh/commit/91a40f2ba1cded5a025004a6143578fa84baec66))
* [x] Situational Awareness ([91a40f2](https://github.com/cbrnrd/maliketh/commit/91a40f2ba1cded5a025004a6143578fa84baec66))

## Star History

<a href="https://star-history.com/#cbrnrd/maliketh&Date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=cbrnrd/maliketh&type=Date&theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=cbrnrd/maliketh&type=Date" />
    <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=cbrnrd/maliketh&type=Date" />
  </picture>
</a>

<p align="center"><a href="https://github.com/cbrnrd/maliketh#"><img src="http://randojs.com/images/backToTopButtonTransparentBackground.png" alt="Back to top" height="29"/></a></p>