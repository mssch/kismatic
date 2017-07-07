# Kismatic

### Pre-requisites
- make
- Docker

### Build using make
We use `make` to clean, build, and produce our distribution package. Take a look at the Makefile for more details.

Build and test phases happen inside docker containers.

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