peco-docker
===========

These are containers for testing (and, releasing, in the future) [peco](https://github.com/peco/peco)

Rationale: I'm sick of relying on third party CI tools and what not. I have docker. Let's just use that to release from my own machine.

## Usage

First, you need to build the images. Simply run:

```
go run build-images.go
```

Once you have the images, you can easily run tests, release files, etc.

## Testing

Simply run:

```
./script/test.sh
```

## Releasing

You need a file containing your github token somewhere

```
GITHUB_TOKEN_FILE=/path/to/your_token PECO_VERSION=vX.Y.Z ./script/release.sh
```

## TODO

* Make this work seemlessly from the peco work directory
* Make release.sh update homebrew-peco as well