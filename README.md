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
FILE                     STATUS  LENGTH  RETRY  POISONED  URL                                                                                        ERROR  
get.json                 200     1038                                                                                                                       
x-forwarded-host.json    404     0                                                                                                                       
x-forwarded-host.json    404     0       true   true      https://martinvks.no?id=dcbeef4b-8c08-4ef4-a5c5-d8da1eec9604         
x-forwarded-scheme.json  200     1038                     
```

The poison command will:
1. Fetch the target resource with a normal GET request
2. For each of the json request files, send the (possibly malformed) request with a unique id query param
3. If the status code or the length of the response body is different from the normal GET request, 
retry the request with a normal GET request and the same id query param
4. If the status code or the length of the response body is still different, log the url as poisoned

Since it only compares the status code and the length of the response body this command will produce
a lot of false positives if the content of the target resource is dynamically created.

## JSON request files

Requests are defined in JSON files

```json
{
  "headers": [
    {
      "name": ":method",
      "value": "POST"
    },
    {
      "name": "content-type",
      "value": "application/x-www-form-urlencoded"
    }
  ],
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
* **headers** `array<header>`  
Header fields sent in the first HEADERS frame  
Default pseudo-header values can be replaced by adding a header field with the pseudo-header name
* **continuation** `array<header>`  
Header fields sent in a CONTINUATION frame  
Only works when using `h2` protocol
* **body** `string`  
The request body
* **trailer** `array<header>`  
Trailer fields sent in the last HEADERS frame

### Header fields

* **name** `string`  
Header field name
* **value** `string`  
Header field value

### Examples

Environment variables can be used with the `"${ENVIRONMENT_VARIABLE_KEY}"` syntax
```json
{
  "headers": [
    {
      "name": ":method",
      "value": "${REQUEST_METHOD}"
    }
  ]
}
```
Control characters can be added to header field names, header field values and the body by escaping them.  
See the [JSON RFC](https://www.rfc-editor.org/rfc/rfc8259.html#section-7) for more details.
```json
{
  "headers": [
    {
      "name": "x-smuggle-header",
      "value": "foo\r\nx-another-header: bar"
    },
    {
      "name": "x-null-byte",
      "value": "ab\u0000c"
    }
  ]
}
```
When `addDefaultHeaders` is true, default values can be replaced by adding a header field with the pseudo-header name to `headers`.
Any extra pseudo-headers will be added to the HEADERS frame.
For example sending two `:path` header fields can be done with:
```json
{
  "headers": [
    {
      "name": ":path",
      "value": "/"
    },
    {
      "name": ":path",
      "value": "/admin"
    }
  ]
}
```