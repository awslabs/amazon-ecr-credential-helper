# Amazon ECR Docker Credential Helper

The Amazon ECR Docker Credential Helper is a 
[credential helper](https://github.com/docker/docker-credential-helpers)
for the Docker daemon that makes it easier to use Amazon EC2 Container
Registry.

## Installation

You must have at least Docker 1.11 installed.

Place the `docker-credential-ecr-login` binary on your `PATH` and set the contents
of your `~/.docker/config.json` file to be:

```json
{
	"credsStore": "ecr-login"
}
```

## Usage

`docker push 123457689012.dkr.ecr.us-west-2.amazonaws.com/my-repository:my-tag`

`docker push 123457689012.dkr.ecr.us-west-2.amazonaws.com/my-repository:my-tag`

There is no need to `docker login` or `docker logout`.

## Building

You'll need Go 1.4 or greater, git, and make installed.

Check out this repository into your existing `GOPATH`, then run `make`.  A
binary will be produced in `bin/local/docker-credential-ecr-login`.

## Troubleshooting

Logs are placed in `~/.ecr/log`.

## License

The Amazon ECR Docker Credential Helper is licensed under the Apache 2.0 License.