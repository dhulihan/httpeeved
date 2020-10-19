# httpeeved

A simple, configurable mock webserver that cycles through good or bad responses.

## Install

```sh
go get -u github.com/dhulihan/httpeeved
```

## Getting Started

Run this:

```
$ httpeeved
```

Then, in a separate shell, run this:

```sh
$ curl -I localhost:8080
HTTP/1.1 200 OK

$ curl -I localhost:8080
HTTP/1.1 400 Bad Request

$ curl -I localhost:8080
HTTP/1.1 404 Not Found

$ curl -i -X PUT localhost:8080 -d '{ "foo": "bar"}'
HTTP/1.1 502 Bad Gateway
Content-Type: application/json; charset=utf-8
Date: Mon, 19 Oct 2020 23:25:09 GMT
Content-Length: 58

{"body":"{ \"foo\": \"bar\"}","code":"502","method":"PUT"}%
```

## Examples

```sh
# default configuration: round robin between several response codes (200, 206, 400, 404, 500, 502)
httpeeved

# round robin between 200 and 502 responses
httpeeved --codes=200 --codes=502 --selection-strategy=round-robin
```

## Usage

```
httpeeved [OPTIONS]

Application Options:
  -v, --verbose                                 Show verbose debug information: -v for debug, -vv for trace.
  -a, --addr=                                   Address to bind too (default: :8080)
  -c, --codes=                                  Repsonse status codes. Can be specified many times. (default: 200, 206, 400, 404, 500, 502)
  -s, --selection-strategy=[round-robin|random] response code selection strategy (default: round-robin)
  -r, --responses=                              use this to set a custom response message

Help Options:
  -h, --help                                    Show this help message
```

