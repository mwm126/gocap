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

Run `go generate` to download embedded dependencies (Yubikey Personalization, Putty.exe, etc.)

```
go generate ./...
```

This only needs to be done once (unless the files in `cap/embeds` are deleted)

## Run

To build and run for development:

```
go run main.go
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
