# Packaging for Debian

Packaging sources for Debian are contained in the
[`debian` branch](https://github.com/awslabs/amazon-ecr-credential-helper/tree/amazonlinux).

This branch is updated every time a release is prepared for Debian.

### Testing with local, unreleased sources

To build a package locally for testing from unreleased sources, you can use the
[`debian/docker-deb` script](docker-deb).  This script will make a new
temporary clone of this local repository, check out a tag, create a filtered
source tarball (all the sources, minus the vendor folder), and then run a
Docker container to perform the build from the source tarball.  Because the
script operates from a temporary clone, it is important to commit your changes
before invoking the script.

To test locally, you can use the following workflow:

* Make changes in a branch derived from upstream
* Update the [`VERSION`](/VERSION) file
* Tag the branch with the same version indicated in the file
* Merge the branch into this branch
* Delete any files that have been added to the
  [`vendor` folder](/ecr-login/vendor)
* Update the [`debian/changelog`](changelog) file with the new version (or
  release) number and with appropriate changes listed
* Run `./debian/docker-deb test`

### Prepare a release

To prepare a release, you can use the following workflow:

* Merge the release tag into this branch
* Delete any files that have been added to the
  [`vendor` folder](/ecr-login/vendor)
* Update the [`debian/ORIG_URL`](ORIG_URL) and
  [`debian/ORIG_SHA256`](ORIG_SHA256) files with the upstream
  "release-novendor.tar.gz" file and SHA256 sum.
* Update the [`debian/changelog`](changelog) file with the new version (or
  release) number and with appropriate changes listed

To build a package ready for submission, you can use the same
[`debian/docker-deb` script](docker-deb), but you want to ensure that you have
the final commits in place and a location where you can preserve the generated
files.  When run without the "test" argument, the script will download the
upstream tarball from the location specified in the
[`debian/ORIG_URL`](ORIG_URL) file and verify it with the sha256 sum specified
in the [`debian/ORIG_SHA256`](ORIG_SHA256) file.  After the script runs, all
the appropriate files will be placed in the root directory of this repository,
including the "orig.tar.gz" file, the "debian.tar.gz" file, the ".changes"
file, and the ".dsc" file.

To submit a package, it must be signed with your GPG key (using the
`debsign(1)` utility) and can then be uploaded to
[mentors.debian.org](https://mentors.debian.org) using `dput(1)`.  Instructions
for configuring `dput(1)` are
[available on the mentors website](https://mentors.debian.net/intro-maintainers).
Once uploaded, a Debian Developer (or Debian Maintainer) must sign the package
with their key and upload it to Debian.
