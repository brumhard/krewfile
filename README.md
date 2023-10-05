<img src="docs/krewfile.png" width="480" alt="Krewfile logo"/>

# krewfile

krewfile is just like [Brewfile (brew bundle)](https://github.com/Homebrew/homebrew-bundle) or [Gemfiles](https://bundler.io/man/gemfile.5.html) but for the [krew](https://github.com/kubernetes-sigs/krew) kubernetes plugin manager.

## Demo

![krewfile demo](docs/term-animation.svg)

## Usage

Define a krewfile like the following at `~/.krewfile`:

```krewfile
explore
modify-secret
neat
oidc-login
pv-migrate
stern
krew
```

Now run `krewfile` and it will install all the plugins defined in the file and also remove all the plugins that are not in the file.

You can also put your krewfile at any other location and point to that using the `-file` CLI parameter.

Lastly, you can use the `-command` flag to overwrite the binary to call. By default, `krew` is used, but you might as well use `-command "kubectl krew"` to use the kubectl plugin instead.

## Installation

### From source

If you have Go 1.16+, you can directly install by running:

```bash
go install github.com/brumhard/krewfile@latest
```

> Based on your go configuration the `krewfile` binary can be found in `$GOPATH/bin` or `$HOME/go/bin` in case `$GOPATH` is not set.
> Make sure to add the respective directory to your `$PATH`.
> [For more information see go docs for further information](https://golang.org/ref/mod#go-install). Run `go env` to view your current configuration.

### nix

This repo contains a [`flake.nix`](flake.nix) file which you can use for example with

```shell
nix run github:brumhard/krewfile# -- -help
```
