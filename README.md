# httptestrunner

An HTTP client for sending malformed HTTP/2 and HTTP/3 requests.  
Based on [http2smugl](https://github.com/neex/http2smugl) written by Emil Lerner.

## Installation

```
go install github.com/Martinvks/httptestrunner@latest
```

## Usage

For information about usage, run:

```
httptestrunner
```

## JSON request files

Requests are defined in JSON files

```json
{
  "headers": [
    {
      "name": ":method",
      "value": "POST"
    }
  ],
  "body": "hello"
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
`\r`, `\n` and `\u0000` special characters can be used in the body, header field names and header field values
```json
{
  "headers": [
    {
      "name": "x-example",
      "value": "foo\r\nx-another-header: bar"
    }
  ]
}
```
When `addDefaultHeaders` is true, adding a pseudo-header with the same name to `headers` will replace the default value.  
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