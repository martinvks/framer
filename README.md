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
httptestrunner --help
```

### Single command

Sends a single request to the target URL and prints the response to console

```
httptestrunner single -f ./request.json https://martinvks.no
```

### Multi command

Sends multiple requests to the target URL and prints the response status code or error to console

```
httptestrunner multi -d ./pseudo-headers https://martinvks.no
```

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