# fiberresp

A standardized response handling package for Fiber web framework with i18n support.

[![Go Reference](https://pkg.go.dev/badge/github.com/prongbang/fiberresp.svg)](https://pkg.go.dev/github.com/prongbang/fiberresp)
[![Go Report Card](https://goreportcard.com/badge/github.com/prongbang/fiberresp)](https://goreportcard.com/report/github.com/prongbang/fiberresp)

## Features

- Standardized JSON response format
- Internationalization (i18n) support
- Customizable error codes and messages
- Fluent interface for building responses
- Support for status codes, causes, and localization parameters

## Installation

```bash
go get github.com/prongbang/fiberresp
```

## Usage

### Basic Usage

```go
import (
    "github.com/gofiber/fiber/v2"
    "github.com/prongbang/fiberresp"
    "net/http"
)

func handler(c *fiber.Ctx) error {
    err := fiberresp.New("ERR001", "errors.not_found")
    return fiberresp.With(c).Response(err)
}
```

### Setting Up i18n

Configure i18n in your Fiber application:

```go
import (
    "github.com/gofiber/contrib/fiberi18n/v2"
    "github.com/gofiber/fiber/v2"
    "golang.org/x/text/language"
)

func main() {
    app := fiber.New()
    app.Use(fiberi18n.New(&fiberi18n.Config{
        RootPath:         "./localize",
        FormatBundleFile: "json",
        AcceptLanguages:  []language.Tag{language.Thai, language.English},
        DefaultLanguage:  language.English,
    }))

    // Routes
    // ...
}
```

### Predefined Error Responses

Create reusable error responses:

```go
func NewBadRequest() *fiberresp.ResponseBody {
    return &fiberresp.ResponseBody{
        Code:         "CLT001",
        Message:      "errors.bad_request",
        StatusCode:   http.StatusBadRequest,
        LocaleParams: make(map[string]any),
    }
}

func NewNotFound() *fiberresp.ResponseBody {
    return &fiberresp.ResponseBody{
        Code:         "CLT002",
        Message:      "errors.not_found",
        StatusCode:   http.StatusNotFound,
        LocaleParams: make(map[string]any),
    }
}

func NewUnauthorized() *fiberresp.ResponseBody {
    return &fiberresp.ResponseBody{
        Code:         "AUT001",
        Message:      "errors.unauthorized",
        StatusCode:   http.StatusUnauthorized,
        LocaleParams: make(map[string]any),
    }
}

func NewFieldRequired() *fiberresp.ResponseBody {
    return &fiberresp.ResponseBody{
        Code:       "VAL001",
        Message:    "validation.field.required",
        StatusCode: http.StatusBadRequest,
        LocaleParams: map[string]any{
            "field": "username",
        },
    }
}
```

### Using Localization Parameters

Add dynamic parameters to error messages:

```go
// Field validation example
app.Get("/validate-email", func(c *fiber.Ctx) error {
    err := fiberresp.New("VAL001", "validation.field.required").
        WithParam("field", "email")
    return fiberresp.With(c).Response(err)
})

// Length validation example
app.Get("/validate-username", func(c *fiber.Ctx) error {
    err := fiberresp.New("VAL002", "validation.username.length").
        WithParam("min", 3).
        WithParam("max", 20)
    return fiberresp.With(c).Response(err)
})
```

### Complete Example

```go
package main

import (
    "github.com/gofiber/contrib/fiberi18n/v2"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/log"
    "github.com/prongbang/fiberresp"
    "golang.org/x/text/language"
    "net/http"
)

func main() {
    app := fiber.New()
    app.Use(fiberi18n.New(&fiberi18n.Config{
        RootPath:         "./localize",
        FormatBundleFile: "json",
        AcceptLanguages:  []language.Tag{language.Thai, language.English},
        DefaultLanguage:  language.English,
    }))

    app.Get("/test-bad-request", func(c *fiber.Ctx) error {
        return fiberresp.With(c).Response(NewBadRequest())
    })

    app.Get("/test-username-required", func(c *fiber.Ctx) error {
        return fiberresp.With(c).Response(NewFieldRequired())
    })

    app.Get("/test-email-required", func(c *fiber.Ctx) error {
        err := fiberresp.New("VAL001", "validation.field.required").
            WithParam("field", "email")
        return fiberresp.With(c).Response(err)
    })

    log.Fatal(app.Listen(":3000"))
}

func NewBadRequest() *fiberresp.ResponseBody {
    return &fiberresp.ResponseBody{
        Code:         "CLT001",
        Message:      "errors.bad_request",
        StatusCode:   http.StatusBadRequest,
        LocaleParams: make(map[string]any),
    }
}

func NewFieldRequired() *fiberresp.ResponseBody {
    return &fiberresp.ResponseBody{
        Code:       "VAL001",
        Message:    "validation.field.required",
        StatusCode: http.StatusBadRequest,
        LocaleParams: map[string]any{
            "field": "username",
        },
    }
}
```

## Localization Example

Example of localization JSON files:

`./localize/en.json`:
```json
{
  "errors.bad_request": "Bad Request",
  "errors.not_found": "Not Found",
  "errors.unauthorized": "Unauthorized",
  "validation.field.required": "The {{.field}} field is required",
  "validation.username.invalid": "Invalid username",
  "validation.username.length": "Username must be between {{.min}} and {{.max}} characters"
}
```

`./localize/th.json`:
```json
{
  "errors.bad_request": "คำขอไม่ถูกต้อง",
  "errors.not_found": "ไม่พบข้อมูล",
  "errors.unauthorized": "ไม่มีสิทธิ์เข้าถึง",
  "validation.field.required": "กรุณากรอกข้อมูล {{.field}}",
  "validation.username.invalid": "ชื่อผู้ใช้ไม่ถูกต้อง",
  "validation.username.length": "ชื่อผู้ใช้ต้องมีความยาวระหว่าง {{.min}} ถึง {{.max}} ตัวอักษร"
}
```

## Response Structure

The JSON response format:

```json
{
  "code": "ERR001",
  "message": "Not Found",
  "data": null,
  "cause": "Optional error cause"
}
```

## API Reference

### ResponseBody

```go
type ResponseBody struct {
    Code         string         `json:"code"`
    Data         any            `json:"data"`
    Message      string         `json:"message"`
    Cause        *string        `json:"cause,omitempty"`
    StatusCode   int            `json:"-"`
    LocaleParams map[string]any `json:"-"`
}
```

### Methods

- `New(code string, message string) *ResponseBody`: Create a new error response
- `WithData(data any) *ResponseBody`: Set response data
- `WithStatusCode(statusCode int) *ResponseBody`: Set HTTP status code
- `WithCause(cause string) *ResponseBody`: Set error cause
- `WithMessage(message string) *ResponseBody`: Set error message
- `WithParams(params map[string]any) *ResponseBody`: Set multiple localization parameters
- `WithParam(key string, value any) *ResponseBody`: Set a single localization parameter
- `With(ctx *fiber.Ctx) *Config`: Initialize response with Fiber context
- `Response(resp *ResponseBody) error`: Send response to client
- `Response(ctx *fiber.Ctx, resp *ResponseBody) error`: Send response to client

## License

MIT
