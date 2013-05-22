/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// go11 is a web server that announces whether or not Go 1.1 has been tagged.
package main

import (
	"expvar"
	"flag"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

const changeURL = "https://code.google.com/p/go/source/detail?r=go1.1"

var (
	httpAddr   = flag.String("http", "localhost:8080", "Listen address")
	pollPeriod = flag.Duration("poll", 5*time.Second, "Poll period")
)

// Exported variables OMIT
var (
	hitCount       = expvar.NewInt("hitCount")
	pollCount      = expvar.NewInt("pollCount")
	pollError      = expvar.NewString("pollError")
	pollErrorCount = expvar.NewInt("pollErrorCount")
)

func main() {
	flag.Parse()
	go poll(*pollPeriod)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

var state struct {
	sync.RWMutex
	yes bool // whether Go 1.1 has been tagged.
}

func poll(period time.Duration) {
	for !isTagged() {
		time.Sleep(period)
	}
	state.Lock()
	state.yes = true
	state.Unlock()
}

func isTagged() bool {
	pollCount.Add(1) // HL
	r, err := http.Head(changeURL)
	if err != nil {
		log.Print(err)
		pollError.Set(err.Error()) // HL
		pollErrorCount.Add(1)      // HL
		return false
	}
	return r.StatusCode == http.StatusOK
}

func handler(w http.ResponseWriter, r *http.Request) {
	hitCount.Add(1) // HL
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
