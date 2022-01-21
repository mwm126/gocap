# GOCAP CAP client

CAP Client reimplemented in Go, with the [Fyne Toolkit](https://fyne.io)


## Dependencies

### Mac

``` shell
brew install go gh
```

### Ubuntu

``` shell
sudo apt install golang gcc libgl1-mesa-dev xorg-dev
```

### Windows

[Scoop](https://scoop.sh) is the recommended way to install windows dependencies.

```
scoop install curl
scoop install gcc
scoop install go
scoop install innoextract
scoop install innounp
scoop install yubikey-personalization
```

You need to set environment variables for the location of the Yubikey library when building:

```
notepad $PROFILE    # to set env variables permanently; or run the following each time opening the terminal
$env:CGO_CFLAGS="-I$(scoop prefix yubikey-personalization)/include -I$(scoop prefix yubikey-personalization)/include/ykpers-1"
$env:CGO_LDFLAGS="-L$(scoop prefix yubikey-personalization)/lib"
```

## Generate

Run `go generate` to download embedded dependencies (Yubikey Personalization, Putty.exe, etc.)

```
go generate ./...
```

This only needs to be done once (unless the files are deleted, such as with `git clean`)

## Run

To build and run for development:

```
go run .
```


## Run (demo)

You can also run a demo version of the client, which runs a local SSH server for
testing with dummy data (Lorem Ipsum Etcetera) in order to demonstrate the
interface. This does not require a Yubikey.

To build and run the interface demo:

```
go run demo/main.go
```



## Release

Install [Fyne Cross](https://github.com/fyne-io/fyne-cross) for cross compiling Fyne apps using Docker.

``` shell
go install github.com/fyne-io/fyne-cross
```

By default this installs the `fyne-cross` command to `~/go/bin`, so add `~/go/bin` to your `PATH`.

Read the Fyne Cross documentation about building a docker image for OSX/Darwin/Apple with `fyne-cross darwin-image`.

Edit the `release.sh` shell script to set TAG to the new version and NOTES to the release notes.

``` shell
./release.sh
```
