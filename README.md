# Gupper

A golang `app` that will auto-update itself from a `server` package repository compose this project.

## Design

The application is a single golang binary. On startup, the app will attempt to query a remote package repository, download and install a new binary if needed, then restart. It may be run as a mock daemon as well, periodically checking for updates.

A single package server is accessible via http requests. It can tell a client what packages it has, and allows downloads of any package it has.

### App

The app checks for an update on startup by querying the package server for what versions are available. If the app's current version is behind an available package, it will then download the appropriate package for its OS and architecture. The "package" is installed via an atomic mv to replace the currently running binary. To restart, the application leverages the exec system call to start running the newly downloaded binary.

The `update` module and `main` are the only two parts that know there's a specific http server, and `update` encapsulates this s.t. a replacement package repository is easily implemented. It gets the url from `main` which can be specified as a command-line parameter.

The app queries the server for all packages it has available, without any filtering or other criteria. If this were to grow to support multiple applications and lots of versions, some basic parameters on the list query would be a good idea.

There is no authentication from the app to the server. For now, it is as simple as possible, but this needs to be addressed before integrating into a production stack. Additionally, there is no checksum or other additional method to verify the downloaded package, beyond what http includes. Since the install wipes the current app this can be a problem. A backup and recovery mechanism is recommended, but considered out of the scope at this point.

### Server

The server is designed to be a simple model of any sort of remote or local file system, providing the ability to list and download files. It is a dead simple implementation to meet the minimum requirements for the updating app. This could easily be replaced by any number of production ready file servers or storage systems, like s3 or hdfs.

The server currently responds to `GET /list` JSON-encoded array of the available packages on the server's file system. It keeps no metadata, with all the package information contained in the package name, and scanning the directory with packages every time. That minimizes race conditions and reduces complexity at the cost of more filesystem accesses. Files are downloaded via normal http static file serving, via `GET /download/package-name`.

The server doesn't have an upload over http to avoid worrying about authentication and a host of security issues that would add. Instead, it is assumed access via scp or other means to the host are available.

Currently, https is not the default for the server, it is assumed that integration with a production stack would include the appropriate security measures.

### Why Golang vs Python

#### Pros
* Built-in networking facilities in golang, made the http server & requests really clean and simple without having to rely on something like Flask or bigger in a similar python app.
* Can be done with zero dependencies.
* Can encode a version in the app at compile time, making it really easy to create different "versions" for testing, binaries always know what version they are.

#### Cons
* I am much more familiar with python than golang, so some of the time was spent learning and not doing things right. Also less familiar with good, idomatic go.
* Binary needs to be cross-compiled for different operating systems.
* Application could likely have been multiple files to be a well structured program, and as a result would have had to implement multi-file package download/install via archives

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

* golang >=1.10
* bash
* curl (is it?)
* docker, for containerized deployment and testing (usable without on a [unix system](https://media.giphy.com/media/zWXSCGdE2nDYQ/giphy.gif))

### Installing

Clone the repo or use `go get`
```
go get -d -t github.com/twschum/gupper
cd $GOPATH/src/github.com/twschum/gupper/
```
Build the app and server

### Running

The default `app` snd `server` can be run right from the project directory.
Starting the server, then running the app from another shell will look like this
```
$ server
2019/09/12 19:25:14 Listening on port :8080
2019/09/12 19:28:55 GET /list from [::1]:50593
2019/09/12 19:28:55 GET /download/app-1.9.3-darwin-amd64 from [::1]:50593
2019/09/12 19:28:55 GET /list from [::1]:50596
```

```
$ app
2019/09/12 19:28:55 appetizer version 0.0.0 darwin/amd64
2019/09/12 19:28:55 Checking for updates
2019/09/12 19:28:55 Downloading latest package version 1.9.3 from
http://localhost:8080/download/app-1.9.3-darwin-amd64
2019/09/12 19:28:55 Installing app-1.9.3-darwin-amd64 to app
2019/09/12 19:28:55 Restarting...
2019/09/12 19:28:55 app version 1.9.3 darwin/amd64
2019/09/12 19:28:55 Checking for updates
2019/09/12 19:28:55 Up to date
doing useful work now...
```

The server can be run directly, or in the docker container, provided a package directory and external port
`docker run -d -p 8080:8080 -v $(pwd)/packages:/var/packages gupper-server` will run the server in the background

The app can be run in the project dir as a one shot or with `--daemon` which will keep it running and checking for updates every 5 seconds.

If the app is running in daemon mode, try adding some updated packages to the package dir
```
build_packages.sh 1.3
```

## Testing

The bash script `testing.sh` is provided, and will by default run go unit tests (there are not that many of these yet). It then runs an integration test as follows:
* builds some packages and puts them in a directory that will be mounted into the server
* runs the `gupper-server` docker container built by `setup_docker.sh`
* builds and starts an "old" version `app` (latest code, the app is just told it is out of date)
* runs the app with `--version` and checks that the app is expected latest version
* (cleanup)

## Deployment

The application can be compiled for anything in the [go architecture and os list](https://github.com/golang/go/blob/master/src/go/build/syslist.go), and has been tested on darwin/amd64, TODO

The server and testing requires and environment that can compile and run golang, run a docker container, and has unix-like bash. The server assumes a unix-like environment.

## Built With

The golang standard lib and zero dependencies

## Authors

* **Tim Schumacher** - *Core Author* - [twschum](https://github.com/twschum)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
