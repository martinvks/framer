# httptestrunner

An HTTP client for sending (possibly malformed) HTTP/2 and HTTP/3 requests.  
Based on [http2smugl](https://github.com/neex/http2smugl) written by Emil Lerner.

## Installation

```
go install github.com/Martinvks/httptestrunner@latest
```

## Usage

For information about available flags, run:

```
httptestrunner [command] --help
```

### Single command

Send a single request to the target URL and print the response to console

```
$ httptestrunner single -f ./request.json https://martinvks.no/index.js
:status: 200
last-modified: Sat, 26 Nov 2022 15:15:56 GMT
content-type: application/javascript
content-length: 25
date: Mon, 12 Dec 2022 11:16:32 GMT
age: 0
server: ATS/10.0.0

console.log("index.js!");
```

### Multi command

Send multiple requests to the target URL and print the response status code and body length or error to console

```
$ httptestrunner multi -d ./requests https://martinvks.no
FILE                                    STATUS  LENGTH  ERROR                                  
get.json                                200     222                                            
head.json                               200     0                                              
multiple_authority_pseudo_headers.json  400     0                                              
multiple_method_pseudo_headers.json                     RST_STREAM: error code PROTOCOL_ERROR 
```

### Poison command

Send multiple requests to the target and check for cache poisoning

```
$ httptestrunner poison -d ./requests https://martinvks.no
FILE                         STATUS  LENGTH  RETRY  POISONED  ERROR  
get.json                     200     222                             
x-forwarded-host.json        200     222                             
x-http-method-override.json  405     0                               
x-http-method-override.json  405     0       true   true             
```

The poison command will:

1. Fetch the target resource with a normal GET request
2. For each of the json request files, send the (possibly malformed) request with a unique id query param
3. If the response is cacheable and different from the normal GET request,
   retry the request with a normal GET request and the same id query param
4. If the response is still different, log it as poisoned

It compares status code, location header and response body length.
Since dynamically created resources can have varying response length this command might produce a lot of false
positives.

## JSON request files

Requests are defined in JSON files

```json
{
  "headers": {
    ":method": "POST",
    "content-type": "application/x-www-form-urlencoded"
  },
  "body": "param=hello"
}
```

### Fields

* **addDefaultHeaders**  `boolean` (default: `true`)  
  Add the following pseudo-headers to the first HEADERS frame:  
  `:authority` hostname from target URL  
  `:method` GET  
  `:path` path and query part of target URL  
  `:scheme` https  
  Default values can be replaced by adding a header field with the pseudo-header name in \"headers\"
* **headers** `Headers`  
  Header fields sent in the first HEADERS frame
* **continuation** `Headers`  
  Header fields sent in a CONTINUATION frame  
  Only works when using `h2` protocol
* **body** `string`  
  The request body
* **trailer** `Headers`  
  Trailer fields sent in the last HEADERS frame

### Headers

* **[string]** `string | array<string>`

### Examples

When `addDefaultHeaders` is true, the default value can be replaced by adding a header field with the pseudo-header
name to `headers`.

```json
{
  "headers": {
    ":authority": "evil.com"
  }
}
```

Environment variables can be used with the `${ENVIRONMENT_VARIABLE_KEY}` syntax

```json
{
  "headers": {
    ":method": "${REQUEST_METHOD}"
  }
}
```

Control characters can be added to header field names, header field values and the body by escaping them.  
See the [JSON RFC](https://www.rfc-editor.org/rfc/rfc8259.html#section-7) for more details.

```json
{
  "headers": {
    "x-smuggle-header": "foo\r\nx-another-header: bar",
    "x-null-byte": "ab\u0000c"
  }
}
```

To send multiple header fields with the same header field name, specify the values as an array

```json
{
  "headers": {
    "cookie": [
      "a=b",
      "c=d",
      "e=f"
    ]
  }
}
```

A trailer section can be added to the request with the `trailer` field

```json
{
  "headers": {
    ":method": "POST",
    "trailer": "Foo"
  },
  "body": "some data",
  "trailer": {
    "foo": "bar"
  }
}
```