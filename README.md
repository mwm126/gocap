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

TODO

Run development build in MSYS2 environment

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
