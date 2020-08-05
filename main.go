package main

import (
	"fmt"

	"github.com/anz-bank/sysl-playground/pkg/syslUtil"
	"github.com/anz-bank/sysl-playground/pkg/urls"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

// PageView is our main page component.
type PageView struct {
	vecty.Core
	Input     string
	Command   string
	Link      string
	InputLink string
	Plantuml  string
}

func main() {
	input, cmd := setup()

	vecty.SetTitle("Sysl Playground")
	vecty.RenderBody(&PageView{
		Input:    input,
		Command:  cmd,
		Plantuml: "http://plantuml.com/plantuml",
	})
}

func setup() (string, string) {
	playgroundUrl, _ := urls.LoadQueryParams()
	input, cmd := urls.DecodeQueryParams(playgroundUrl)

	if input == "" {
		input = `MobileApp:
        Login:
                        Server <- Login
        !type LoginData:
                        username <: string
                        password <: string
        !type LoginResponse:
                        message <: string
Server:
        Login(data <: MobileApp.LoginData):
                        return MobileApp.LoginResponse`
	}
	if cmd == "" {
		cmd = "sysl sd -o \"project.svg\" -s \"MobileApp <- Login\" tmp.sysl"
	}
	return input, cmd
}

// Render implements the vecty.Component interface.
func (p *PageView) Render() vecty.ComponentOrHTML {
	return elem.Body(
		vecty.Markup(
			vecty.Class("body"),
		),
		elem.Header(
			vecty.Text("Sysl Playground"),
			vecty.Markup(
				vecty.Style("font-family", "monospace"),
				vecty.Style("font-size", "30px"),
			),
		),
		elem.Article(
			vecty.Text("Welcome to the Sysl Playground"),
			vecty.Markup(
				vecty.Style("font-family", "monospace"),
				vecty.Style("font-size", "20px"),
			),
		),
		elem.Article(
			vecty.Text("Service to use:"),
			vecty.Markup(
				vecty.Style("font-family", "monospace"),
				vecty.Style("font-size", "15px"),
			),
		),

		// Display a textarea on the right-hand side of the page.
		elem.Table(
			elem.TableRow(
				elem.TableData(
					elem.Table(
						elem.TableRow(
							elem.TableData(
								elem.TextArea(
									vecty.Markup(
										vecty.Style("font-family", "monospace"),
										vecty.Style("font-size", "17px"),
										vecty.Property("rows", 2),
										vecty.Property("cols", 70),

										// When input is typed into the textarea, update the local
										// component state and rerender.
										event.Input(func(e *vecty.Event) {
											p.Plantuml = e.Target.Get("value").String()
											vecty.Rerender(p)
										}),
									),
									vecty.Text(p.Plantuml),
								),
							),
						),
						elem.TableRow(
							elem.TableData(
								elem.TextArea(
									vecty.Markup(
										vecty.Style("font-family", "monospace"),
										vecty.Style("font-size", "17px"),
										vecty.Property("rows", 14),
										vecty.Property("cols", 70),

										// When input is typed into the textarea, update the local
										// component state and rerender.
										event.Input(func(e *vecty.Event) {
											p.Input = e.Target.Get("value").String()
											vecty.Rerender(p)
										}),
									),
									vecty.Text(p.Input), // initial textarea text.
								),
							),
						),
						elem.TableRow(
							elem.TableData(
								elem.TextArea(
									vecty.Markup(
										vecty.Style("font-family", "monospace"),
										vecty.Style("font-size", "17px"),
										vecty.Property("rows", 1),
										vecty.Property("cols", 70),

										// When input is typed into the textarea, update the local
										// component state and rerender.
										event.Input(func(e *vecty.Event) {
											p.Command = e.Target.Get("value").String()
											vecty.Rerender(p)
										}),
									),
									vecty.Text(p.Command), // initial textarea text.
								),
							),
						),
						elem.TableRow(
							elem.TableData(

								elem.Button(
									vecty.Markup(
										vecty.UnsafeHTML("Share"),
										vecty.Style("width", "75px"),
										vecty.Style("height", "30px"),
										event.Click(func(e *vecty.Event) {
											p.Link = urls.EncodeUrl(p.Input, p.Command)
											vecty.Rerender(p)
										}),
									),
								),
							),
						),
						elem.TableRow(
							elem.TableData(
								elem.TextArea(
									vecty.Markup(
										vecty.Style("font-family", "monospace"),
										vecty.Style("font-size", "17px"),
										vecty.Property("rows", 7),
										vecty.Property("cols", 70),
										vecty.Property("wrap", "hard"),
										event.Input(func(e *vecty.Event) {
											p.InputLink = e.Target.Get("value").String()
										}),
									),
									vecty.Text(p.Link),
								),
							),
						),
					),
				),
				elem.TableData(
					&Markdown{Input: p.Input, Command: p.Command, Plantuml: p.Plantuml},
				),
			),
		),
	)

}

// Markdown is a simple component which renders the Input markdown as sanitized
// HTML into a div.
type Markdown struct {
	vecty.Core
	Input    string `vecty:"prop"`
	Command  string `vecty:"prop"`
	Plantuml string `vecty:"prop"`
}

// Render implements the vecty.Component interface.
func (m *Markdown) Render() (res vecty.ComponentOrHTML) {
	defer func() {
		// If panic, then print the error
		if r := recover(); r != nil {
			res = elem.Div(
				vecty.Markup(
					vecty.UnsafeHTML(fmt.Sprintf("%s", r)),
				),
			)
		}
	}()

	output, err := syslUtil.Parse(m.Input, m.Command)
	check(err)
	image := fmt.Sprintf(`<img src="%s/svg/~1%s" width="150% height="150%">`, m.Plantuml, string(output))

	return elem.TableData(vecty.Markup(vecty.UnsafeHTML(image)))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
