package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"net/http"
	"text/template"
)

func main() {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		20,
	)))

	e.Static("/", "public")

	NewTemplateRenderer(e, "public/*.html")

	e.GET("/hello", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.Logger.Fatal(e.Start(":4040"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplateRenderer(e *echo.Echo, pattern string) {
	t := &Template{
		templates: template.Must(template.ParseGlob(pattern)),
	}
	e.Renderer = t
}
