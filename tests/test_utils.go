package tests

import (
	"io"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

func NewEchoContext(method, path string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, body)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}
