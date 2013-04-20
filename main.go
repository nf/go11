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
	yes bool
}

func poll(interval time.Duration) {
	tick := time.NewTicker(interval)
	defer tick.Stop()
	for {
		<-tick.C
		r, err := http.Head(changeURL)
		if err != nil {
			log.Print(err)
			continue
		}
		if r.StatusCode != http.StatusOK {
			continue
		}
		state.Lock()
		state.yes = true
		state.Unlock()
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	state.RLock()
	var data = struct {
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
<!doctype html><html><body><center>
	<h2>Is Go 1.1 tagged yet?</h2>
	<h1>
	{{if .Yes}}
		<a href="{{.URL}}">YES!</a>
	{{else}}
		No.
	{{end}}
	</h1>
</center></body></html>
`))
