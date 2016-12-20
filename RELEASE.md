# Release
[![Build Status](https://snap-ci.com/On8xdVQV0xY5VXICf0Fx0Vq7fVMDUAfU6JFc8Wtt94A/build_image)](https://snap-ci.com/apprenda/kismatic-platform/branch/master)

This document contains the release details for [Kismatic Enterprise Toolkit](https://github.com/apprenda/kismatic).

Kismatic leverages a CI/CD pipeline to build, test and release its software. The pipeline extensively tests the software before any release is published and a build must pass all unit and integration tests to be promoted by the pipeline.

Kismatic depends on [Kuberang](https://github.com/apprenda/kuberang) and the [Kismatic Distro Packages](https://github.com/apprenda/kismatic-distro-packages), along with all other components listed in the [README](https://github.com/apprenda/kismatic).
All components are developed and versioned independently, however a released version of the Kismatic Enterprise Toolkit will have specified versions of the underlying components that have been previously tested for compatibility.   

### Releasing

Draft the release using GitHub's release feature:
* Include a section that lists the major features.
* Include a section that lists any versions that have been changed in this new release
* Inculude any other information that is relevant to the release

Once the master branch is stable and ready for releasing, upload the binaries to GitHub:
* Grab the latest *linux* artifact from the S3 bucket, and rename it to kismatic-vX.Y.Z-linux-amd64.tar.gz
* Grab the latest *darwin* artifact from the S3 bucket, and rename it to kismatic-vX.Y.Z-darwin-amd64.tar.gz
* Upload both to GitHub

