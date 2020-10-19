# httpeeved

A simple, configurable mock webserver that cycles through good or bad responses.

## Installation

```sh
go get -u github.com/dhulihan/httpeeved
```

## Usage

```
Usage:
  main [OPTIONS]

Application Options:
  -v, --verbose                                 Show verbose debug information: -v for debug, -vv for trace.
  -a, --addr=                                   Address to bind too (default: :8080)
  -c, --codes=                                  Repsonse status codes. Can be specified many times. (default: 200, 206, 400, 404, 500, 502)
  -s, --selection-strategy=[round-robin|random] response code selection strategy (default: round-robin)
  -r, --responses=                              use this to set a custom response message

Help Options:
  -h, --help                                    Show this help message

panic: Usage:
  main [OPTIONS]

Application Options:
  -v, --verbose                                 Show verbose debug information: -v for debug, -vv for trace.
  -a, --addr=                                   Address to bind too (default: :8080)
  -c, --codes=                                  Repsonse status codes. Can be specified many times. (default: 200, 206, 400, 404, 500, 502)
  -s, --selection-strategy=[round-robin|random] response code selection strategy (default: round-robin)
  -r, --responses=                              use this to set a custom response message

Help Options:
  -h, --help                                    Show this help message
```

## Examples

in one shell/terminal:

```
$ httpeeved
```

in a separate shell

```sh
$ curl -I localhost:8080
HTTP/1.1 200 OK

$ curl -I localhost:8080
HTTP/1.1 400 Bad Request

$ curl -I localhost:8080
HTTP/1.1 404 Not Found

$ curl -I localhost:8080
HTTP/1.1 502 Bad Gateway
```
