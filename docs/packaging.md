# Packaging Amazon ECR Credential Helper

Amazon maintains packages for the Amazon ECR Credential Helper on the following
operating systems:

* Amazon Linux 2
  ([source](https://github.com/awslabs/amazon-ecr-credential-helper/tree/amazonlinux),
  [packaging documentation](https://github.com/awslabs/amazon-ecr-credential-helper/blob/amazonlinux/docs/packaging-amazon-linux.md))
* Debian
  ([source](https://github.com/awslabs/amazon-ecr-credential-helper/tree/debian))
  (note: the packages in derivatives of Debian like Devuan, Ubuntu 19.04 Disco
  Dingo, and PureOS are derived from the Debian package with no additional
  modifications)
  
There are community-maintained packages for the Amazon ECR Credential Helper on
the following operating systems:

* Mac OS X (with the Homebrew package manager)
  ([source](https://github.com/Homebrew/homebrew-core/blob/master/Formula/d/docker-credential-helper-ecr.rb))
* NixOS (and the Nix package manager)
  ([source](https://github.com/NixOS/nixpkgs/blob/master/pkgs/by-name/am/amazon-ecr-credential-helper/package.nix))
* Arch Linux (in the Arch User Repository)
  ([source](https://aur.archlinux.org/packages/amazon-ecr-credential-helper))

If you are interested in packaging the Amazon ECR Credential Helper, please get
in touch!  We can list your community-maintained packaging here and include
installation instructions in our [README.md](../README.md).
