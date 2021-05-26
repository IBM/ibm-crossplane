# IBM Crossplane

IBM Crossplane is a fork from an opensource project Crossplane with additonal enhancements and modifications from IBM.

## Requirements

- go 1.14.x or go 1.15.x
- helm 3.x

## Supported platforms

Red Hat OpenShift Container Platform 4.6 or newer installed on one of the following platforms:

- Linux x86_64
- Linux on Power (ppc64le)
- Linux on IBM Z and LinuxONE

## Development

#### Setting up build harness and running initial build

```
# make 
```

#### Building image for local architecture

```
# make build
```

#### Building images for all supported platforms

```
# make build.all
```

#### Running unit tests

```
# make test
```

#### Preparing to submit a PR

Before submitting a PR, all these make tasks must be completed without error.

```
# make reviewable
# make build
# make test
```

#### Publishing images

```
# make images
```

