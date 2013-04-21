// go11 is a web server that announces whether or not Go 1.1 has been tagged.
package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

const changeURL = "https://code.google.com/p/go/source/detail?r=go1.1"

var (
	httpAddr     = flag.String("http", "localhost:8080", "Listen address")
	pollInterval = flag.Duration("poll", time.Second*5, "Poll interval")
)

func main() {
	flag.Parse()
	go poll(*pollInterval)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

var state struct {
	sync.RWMutex
	yes bool // True if Go 1.1 has been tagged.
}

func poll(interval time.Duration) {
	for !isTagged() {
		time.Sleep(interval)
	}
	state.Lock()
	state.yes = true
	state.Unlock()
}

func isTagged() bool {
	r, err := http.Head(changeURL)
	if err != nil {
		log.Print(err)
		return false
	}
	return r.StatusCode == http.StatusOK
}

func handler(w http.ResponseWriter, r *http.Request) {
	state.RLock()
	data := struct {
		Yes bool
		URL string
	}{
		Yes: state.yes,
		URL: changeURL,
	}
	state.RUnlock()
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Print(err)
	}
}

var tmpl = template.Must(template.New("root").Parse(`
<!DOCTYPE html><html><body><center>
	<h2>Is Go 1.1 out yet?</h2>
	<h1>
	{{if .Yes}}
		<a href="{{.URL}}">YES!</a>
	{{else}}
		No.
	{{end}}
	</h1>
</center></body></html>
`))
