package main

import (
	"html/template"
	"os"
)

const html = `
<script>var foo = {{.Foo}};</script>
<a href="{{.URL}}">
	{{.Text}}
</a>
`

func main() {
	tmpl := template.Must(template.New("example").Parse(html))
	data := struct {
		Foo       string
		URL, Text string
	}{
		Foo:  `Some "quoted" string`,
		URL:  `" onClick="alert('xss!');`,
		Text: "The <- operator is for channel sends and receives",
	}
	tmpl.Execute(os.Stdout, data)
}
