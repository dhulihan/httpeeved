# httpeeved

A simple fake webserver that cycles through good or bad response codes.

*Make your HTTP clients angry.*

## Usage

```
-v, --verbose                                 Show verbose debug information: -v for debug, -vv for trace.
-a, --addr <addr>                             Address to bind to (default: :8080)
-c, --codes <code>[, <code>, ...]             Repsonse status codes. Can be specified many times. (default: 200, 202, 206, 400, 401, 403, 404, 409, 500, 502)
-s, --selection-strategy=[round-robin|random] response code selection strategy (default: round-robin)
-r, --responses <msg>                         Use this to set a custom response message
-x, --proxy                                   Run as proxy. This will forward requests to destination and modify the status code of the original response.
```


## Install

```sh
go get -u github.com/dhulihan/httpeeved
```

## Examples

* Round robin between several response codes (default):
	```sh
	httpeeved
	```
	```sh
	$ curl -I localhost:8080
	HTTP/1.1 200 OK

	$ curl -I localhost:8080
	HTTP/1.1 400 Bad Request

	$ curl -I localhost:8080
	HTTP/1.1 404 Not Found

	$ curl -I localhost:8080 -d '{"foo": "bar"}'
	HTTP/1.1 502 Bad Gateway
	Content-Type: application/json; charset=utf-8
	Date: Mon, 19 Oct 2020 23:25:09 GMT
	Content-Length: 58

	{"body":"{ \"foo\": \"bar\"}","code":"502","method":"PUT"}%
	```
* Round robin between 200 and 502 responses only
	```sh
	httpeeved --codes=200 --codes=502 --selection-strategy=round-robin
	```
* Randomly respond with different response codes
	```sh
	httpeeved --codes=200 --codes=301 --codes=500 --selection-strategy=random
	```
* Inspect a request containing `multipart/form-data`
	```sh
	curl --location --request POST 'localhost:8000/foo' \
		--form 'some-file=@"foo.txt"' \
		--form 'foo="bar"' \
		--form 'foo="baz"' | jq .

	{
	  "Code": "200",
	  "Content-Type": "multipart/form-data; boundary=------------------------c58c7266e69ad9f9",
	  "Method": "POST",
	  "URL": "/foo",
	  "form": {
	    "foo": [
	      "bar",
	      "baz"
	    ]
	  },
	  "multipart-form-file-some-file-0": "foo.txt map[Content-Disposition:[form-data; name=\"some-file\"; filename=\"foo.txt\"] Content-Type:[text/plain]] 4",
	  "multipart-form-value-foo-0": "bar",
	  "multipart-form-value-foo-1": "baz"
	}
	```
* Run as proxy, which makes legit requests to a destination and modifies response status code
	```sh
	httpeeved -x
	```
	```sh
	$ curl -I -x localhost:8080/ google.com/
	HTTP/1.1 500 Internal Server Error
	```
