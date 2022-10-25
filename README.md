# GTERM - 💕Share your terminal as a web application💕

[![Go Version](http://img.shields.io/github/go-mod/go-version/iotxfoundry/gterm)][gomod]
[![GitHub release](http://img.shields.io/github/release/iotxfoundry/gterm.svg?style=flat-square)][release]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[gomod]: https://github.com/iotxfoundry/gterm/blob/master/go.md
[release]: https://github.com/iotxfoundry/gterm/releases
[license]: https://github.com/iotxfoundry/gterm/blob/master/LICENSE

GTERM is a tool to share your terminal as a web applications.

Unlike [gotty](https://github.com/yudai/gotty), we use the same terminal, so we can see the commands from others.

![gif](./docs/gterm.gif)

# Installation

## `go install`

GTerm requires go1.16 or later.

```sh
$ go install github.com/iotxfoundry/gterm/cmd/gterm@latest
```

## Usage

```sh
Usage:
  gterm [OPTIONS]

Application Options:
  -p, --port= http port (default: 8080)

Help Options:
  -h, --help  Show this help message
```

## License

The MIT License
