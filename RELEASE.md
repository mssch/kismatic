# Release
[![Build Status](https://snap-ci.com/On8xdVQV0xY5VXICf0Fx0Vq7fVMDUAfU6JFc8Wtt94A/build_image)](https://snap-ci.com/apprenda/kismatic-platform/branch/master)

This document contains the release details for [Kismatic Enterprise Toolkit](https://github.com/apprenda/kismatic).

Kismatic leverages a CI/CD pipeline to build, test and release its software. The pipeline extensively tests the software before any release is published and a build must pass all unit and integration tests to be promoted by the pipeline.

Kismatic depends on [Kuberang](https://github.com/apprenda/kuberang) and the [Kismatic Distro Packages](https://github.com/apprenda/kismatic-distro-packages), along with all other components listed in the [README](https://github.com/apprenda/kismatic).
All components are developed and versioned independently, however a released version of the Kismatic Enterprise Toolkit will have specified versions of the underlying components that have been previously tested for compatibility.   
