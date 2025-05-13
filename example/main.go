package main

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/prongbang/fiberresp"
	"golang.org/x/text/language"
	"net/http"
)

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
			"field": "email",
		},
	}
}

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

	app.Get("/test-not-found", func(c *fiber.Ctx) error {
		return fiberresp.With(c).Response(NewNotFound())
	})

	app.Get("/test-unauthorized", func(c *fiber.Ctx) error {
		return fiberresp.With(c).Response(NewUnauthorized())
	})

	app.Get("/test-field-required", func(c *fiber.Ctx) error {
		return fiberresp.With(c).Response(NewFieldRequired())
	})

	app.Get("/test-username-length", func(c *fiber.Ctx) error {
		err := fiberresp.New("VAL002", "validation.username.length").
			WithParam("min", 3).
			WithParam("max", 20)
		return fiberresp.With(c).Response(err)
	})

	app.Get("/test-field", func(c *fiber.Ctx) error {
		err := fiberresp.New("VAL001", "validation.field.required").
			WithParam("field", "อีเมล")
		return fiberresp.With(c).Response(err)
	})

	log.Fatal(app.Listen(":3000"))
}
