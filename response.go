package fiberresp

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
)

type ResponseBody struct {
	Code         string         `json:"code"`
	Data         any            `json:"data"`
	Message      string         `json:"message"`
	Cause        *string        `json:"cause,omitempty"`
	StatusCode   int            `json:"-"`
	LocaleParams map[string]any `json:"-"`
}

func (e *ResponseBody) Error() string {
	return e.Message
}

// New the message support text and localize key
func New(code string, message string) *ResponseBody {
	return &ResponseBody{
		Code:         code,
		Data:         nil,
		Message:      message,
		StatusCode:   http.StatusBadRequest,
		LocaleParams: make(map[string]any),
	}
}

func (e *ResponseBody) WithData(data any) *ResponseBody {
	e.Data = data
	return e
}

func (e *ResponseBody) WithStatusCode(statusCode int) *ResponseBody {
	e.StatusCode = statusCode
	return e
}

func (e *ResponseBody) WithCause(cause string) *ResponseBody {
	e.Cause = &cause
	return e
}

func (e *ResponseBody) WithMessage(message string) *ResponseBody {
	e.Message = message
	return e
}

func (e *ResponseBody) WithParams(params map[string]any) *ResponseBody {
	for k, v := range params {
		e.LocaleParams[k] = v
	}
	return e
}

func (e *ResponseBody) WithParam(key string, value any) *ResponseBody {
	e.LocaleParams[key] = value
	return e
}

type Config struct {
	ctx *fiber.Ctx
}

func With(ctx *fiber.Ctx) *Config {
	return &Config{ctx: ctx}
}

func (c *Config) Response(err *ResponseBody) error {
	return ResponseWith(c.ctx, err)
}

func ResponseWith(ctx *fiber.Ctx, resp *ResponseBody) error {
	if ctx != nil {
		var localizedMessage string
		var localizeErr error

		if len(resp.LocaleParams) > 0 {
			localizedMessage, localizeErr = fiberi18n.Localize(ctx, &i18n.LocalizeConfig{
				MessageID:    resp.Message,
				TemplateData: resp.LocaleParams,
			})
		} else {
			localizedMessage, localizeErr = fiberi18n.Localize(ctx, resp.Message)
		}

		if localizeErr == nil {
			resp.Message = localizedMessage
		}

		_ = ctx.SendStatus(resp.StatusCode)
		return ctx.JSON(resp)
	}
	return resp
}
