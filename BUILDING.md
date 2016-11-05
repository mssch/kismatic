# Kismatic
[![Build Status](https://snap-ci.com/ulSrRsof30gMr7eaXZ_eufLs7XQtmS6Lw4eYwkmATn4/build_image)](https://snap-ci.com/apprenda/kismatic/branch/master)

### Pre-requisites
- Darwin (OSX) or Linux x86-64
- Go installed
- Docker (required for building)

### Build using make
We use `make` to clean, build, and produce our distribution package. Take a look at the Makefile for more details.

In order to build the Go binaries (e.g. Kismatic CLI):
```
make build
```

In order to clean:
```
make clean
```

In order to produce the distribution package:
```
make dist
```
This will produce an `./out` directory, which contains the bits, and a tarball.

You may pass build options as necessary:
```
GOOS=linux make build
```
### Continuous Integration
Kismatic uses Snap for CI.

https://snap-ci.com/apprenda/kismatic/
