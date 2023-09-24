package example

import (
	_ "embed"
	"strings"

	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom-example/example/demos"
	"github.com/mlctrez/godom-example/example/navbar"
	"github.com/mlctrez/godom/app"
)

var _ app.Handler = (*router)(nil)

type router struct{}

func New() app.Handler {
	return &router{}
}

func (e *router) Prepare(ctx *app.ServerContext) {
	ctx.Main = "example/bin/main.go"
	ctx.Output = "build/app.wasm"
	ctx.Watch = []string{"example"}
	ctx.Address = ":8080"
	ctx.ShowWasmSize = true
}

//go:embed head.html
var headHtml string

func (e *router) Headers(ctx *app.Context, header godom.Element) {
	if header.Parent() != nil {
		header.Parent().SetAttribute("data-bs-theme", "dark")
	}
	doc := ctx.Doc
	for _, node := range header.ChildNodes() {
		if node.NodeName() == "title" {
			// TODO: simplify this
			node.ChildNodes()[0] = doc.T("godom-example")
		}
	}
	if len(header.GetElementsByTagName("link")) == 0 {
		header.Body(doc.H(headHtml).ChildNodes()...)
	}
}

func (e *router) Body(ctx *app.Context) godom.Element {
	d := ctx.Doc

	body := d.El("body").Body(navbar.Render(ctx))
	switch ctx.URL.Path {
	case "/", "/godom-example/":
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

func (e *router) Serve(request app.Request, response app.Response) bool {

	switch request.URL().Path {
	case "/diff/api":
		demos.DiffServe(request, response)
		return true
	}
	return false
}

//go:embed 404.html
var was404 string
