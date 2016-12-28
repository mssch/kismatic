# Release
[![Build Status](https://snap-ci.com/On8xdVQV0xY5VXICf0Fx0Vq7fVMDUAfU6JFc8Wtt94A/build_image)](https://snap-ci.com/apprenda/kismatic-platform/branch/master)

This document contains the release details for [Kismatic Enterprise Toolkit](https://github.com/apprenda/kismatic).

Kismatic leverages a CI/CD pipeline to build, test and release its software. The pipeline extensively tests the software before any release is published and a build must pass all unit and integration tests to be promoted by the pipeline.

Kismatic depends on [Kuberang](https://github.com/apprenda/kuberang) and the [Kismatic Distro Packages](https://github.com/apprenda/kismatic-distro-packages), along with all other components listed in the [README](https://github.com/apprenda/kismatic).
All components are developed and versioned independently, however a released version of the Kismatic Enterprise Toolkit will have specified versions of the underlying components that have been previously tested for compatibility.   

### Releasing
Once the master branch is ready for a release, it must be tagged and the tag pushed
to the remote repository.
```
git tag v1.1.0
git push origin v1.1.0
```

In our current CI setup on Snap, pushing a tag to GitHub will not trigger the
build pipeline (not supported by Snap). For this reason, the build must be
triggered manually on Snap once the tag is pushed.

The release stage of the build pipeline is triggered when the commit that is being built
is equal to the commit of the latest tag. The release stage uses the `release.go` script
to draft a release on GitHub and upload the built artifacts.

Once the build is complete, go to `https://github.com/apprenda/kismatic/releases`
to edit the release draft:
* Include a section that lists the major features.
* Include a section that lists any versions that have been changed in this new release
* Include any other information that is relevant to the release

### Release script usage
Before being able to run the release script, you need the following:
* GitHub API token (can be obtained from https://github.com/settings/tokens)
* Built linux artifact, located at `./artifact/linux/kismatic.tar.gz`
* Built darwin artifact, located at `./artifact/darwin/kismatic.tar.gz`
* Git tag to be used (e.g. v1.1.0) for the release

```
GITHUB_TOKEN=theTokenHere go run release.go -tag v1.1.0
```

### Useful Git commands
Get latest tag name:
```
git describe --abbrev=0 --tags
```

Get commit hash for a given tag
```
git rev-list -n 1 $TAG
```
