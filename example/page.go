package example

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom-example/example/demos"
	"github.com/mlctrez/godom-example/example/navbar"
	"github.com/mlctrez/godom/app"
)

var _ app.Handler = (*page)(nil)

type page struct {
	Root     godom.Element `go:"html"`
	DarkMode godom.Element `go:"darkMode"`
}

func New() app.Handler {
	return &page{}
}

func (p *page) Prepare(ctx *app.ServerContext) {
	ctx.Main = "server/main.go"
	ctx.Output = "build/app.wasm"
	ctx.Watch = []string{"example"}
	ctx.Address = ":8080"
	ctx.ShowWasmSize = true
}

//go:embed page.html
var pageHtml string

func (p *page) Html(ctx *app.Context) {
	if ctx.IsWasm() {
		return
	}
	// TODO: use cookies to store dark mode
	p.Root.SetAttribute("lang", "en")
	p.Root.SetAttribute("data-bs-theme", "dark")
	pHtml := ctx.Doc.H(pageHtml)
	p.Root.Body(pHtml.GetElementsByTagName("head")[0])
	p.Root.AppendChild(p.Body(ctx))
}

//go:embed index.html
var indexHtml string

func (p *page) Body(ctx *app.Context) godom.Element {
	d := ctx.Doc

	titleElement := p.Root.GetElementsByTagName("title")[0]
	// TODO: should be able to use titleElement.ReplaceWith(...)
	titleElement.RemoveChild(titleElement.ChildNodes()[0].This())
	titleElement.AppendChild(d.T(fmt.Sprintf("page %s", ctx.URL.Path)))

	fmt.Println("ctx.URL.String()", ctx.URL.String())

	body := d.El("body").Body(navbar.Render(ctx))
	switch ctx.URL.Path {
	case "/":
		body.AppendChild(d.H(indexHtml))
		p.DarkMode.AddEventListener("click", func(event godom.Value) {
			htmlElement := p.Root.This()
			if "dark" == htmlElement.Call("getAttribute", "data-bs-theme").String() {
				htmlElement.Call("removeAttribute", "data-bs-theme")
			} else {
				htmlElement.Call("setAttribute", "data-bs-theme", "dark")
			}
		})
		return body
	case "/exampleOne":
		return body.Body(demos.ExampleOne(ctx))
	case "/diff":
		return body.Body(demos.Diff(ctx))
	case "/editor":
		return body.Body(demos.Editor(ctx))
	default:
		return body.Body(d.H(strings.Replace(was404, "@@page@@", ctx.URL.String(), 1)))
	}
}

//go:embed bootstrap.min.css
var bootCss string

//go:embed bootstrap.bundle.min.js
var bootJs string

func (p *page) Serve(request app.Request, response app.Response) bool {

	if strings.HasSuffix(request.URL().String(), "bootstrap.min.css") {
		response.SetHeader("Content-Type", "text/css")
		_, _ = response.Write([]byte(bootCss))
		return true
	}
	if strings.HasSuffix(request.URL().String(), "bootstrap.bundle.min.js") {
		response.SetHeader("Content-Type", "application/javascript")
		_, _ = response.Write([]byte(bootJs))
		return true
	}

	switch request.URL().Path {

	case "/diff/api":
		demos.DiffServe(request, response)
		return true
	}
	return false
}

//go:embed 404.html
var was404 string
