# Gupper

A golang `app` that will auto-update itself from a `server` package repository compose this project.

## Design

Write a program that is capable of updating itself to a newer version. When a new version is available, deployed programs should be replaced by the newer version. Your solution should be expected to work across common desktop operating systems (Windows, Mac, Linux).

Please write code that you would consider ready to be integrated into a production software stack. If your solution requires non-trivial server components, please write those as well. This problem is under-scoped and very broad. Part of your job is to prevent this task from spiraling into a week of effort.

If you have a question, write it down, and guess what you think my answer would be. Please provide a README describing how your solution works, and what tradeoffs and assumptions you made. Like our own code, we expect an automated build and test process

### Architecture

A single package server is accessible via http requests. It can tell a client what the latest package is, and allows downloads of any package it has.

The application is a single golang binary. On startup, it will query the package server for what the latest version is. If the app is behind that version, it will then download the appropriate version for it's OS and architecture. The "package" is installed via an atomic mv to replace the currently running binary. To restat, the application leverages the exec system call to start running the newly downloaded binary.


#### Server

This could easily be replaced by any number of third-party file servers or storage systems, like s3. This project includes a dead simple implementation to provide a simple way to run the updating app, but is designed such that it would be pretty easy to replace. The logic for picking the latest package lives on the server and would need to either be moved to the client or translated appropriate to the replacement package store.

The server currently responds to `GET /latest` with an application-specific JSON containing the latest package on the server's file system. It keeps no metadata, with all the package information contained in the package name, and scanning the directory with packages every time. That minimizes race conditions and reduces complexity at the cost of more filesystem accesses. Files are downloaded via normal http static file serving.

The server doesn't have an upload over http to avoid worrying about authentication and a host of security issues that would add. Instead, it is assumed access via scp or other means to the host are available.

Were the server to be kept around for the interface, it would also be pretty simple to swap out the local filesystem for another store, be it a database, hdfs, s3, etc.

#### Potential Features

While desiging and working on this, there were definitaly a few features that should/would likely be needed before it's very useful
* https by default
* Packages as actual packages, with multiple files, more setup tasks than a single mv
* Support for more than on "app" on the package server
* Checksums on downloaded packages
* Upload built into the package server instead of just

### Why Golang vs Python

#### Pros
* Built-in networking facilities in golang, made the http server & requests really clean and simple without having to rely on something like Flask or bigger in a similar python app
* Can be done with zero dependencies
* Easy to cross-compile from a single deployment machine to support the app on common desktop operating systems

#### Cons
* I am much more familiar with python than go, so some of the time was spent learning and not doing things right. Also less familiar with good, idomatic go
* Unclear what kind of project this could get integrated into, by python projects are generally more numerous

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

* golang >=1.10
* bash
* docker for containerized deployment and testing (usable without)

### Installing

```
pip install -r requirements.txt
```

### Configuration

Configuration for the flask app is contained in a `config.py` and `instance/config.py`. See the [flask config docs](http://flask.pocoo.org/docs/1.0/config/)

This repo includes an example [instance/config_example.py](instance/config_example.py), to be used for all the secret configuration.

For email notifications with the user system and order notifications, an email account with api access will be required. This is easy to do with gmail, and the base configuration assumes as much.


## Deployment

The original version of this site is running on [PythonAnywhere](pythonanywhere.com), which is the author's recommended deployment solution.

## Built With

The golang standard lib and zero dependencies

## Authors

* **Tim Schumacher** - *Core Author* - [twschum](https://github.com/twschum)

See also the list of [contributors](https://github.com/twschum/mix-mind/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
