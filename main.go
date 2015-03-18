// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/go-on/gopherjslib"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

const failure = `<html><head><title>SRVi</title></head><body style="color:#555555;background:#ffeedd;font-family:Arial;font-size:20px;margin:75px;">%s</body></html>`

var success = `<html><head><title>SRVi</title></head><body><script src="/app.go.js" type="text/javascript"></script></body></html>`

var index *string

var code = ""

func buildHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		http.ServeFile(w, r, r.URL.Path[1:])
		return
	}

	var out bytes.Buffer
	builder := gopherjslib.NewBuilder(&out, nil)

	names := []string{"main.go"}
	if flag.NArg() > 0 {
		names = flag.Args()
	}

	for _, name := range names {
		if path.Ext(name) != ".go" {
			continue
		}

		file, err := os.Open(name)
		if err != nil {
			fmt.Fprintf(w, failure, err)
			return
		}
		defer file.Close()

		builder.Add(name, file)
	}

	if err := builder.Build(); err != nil {
		fmt.Fprintf(w, failure, err)
		return
	}

	code = out.String()

	if *index != "" {
		data, err := ioutil.ReadFile(*index)
		if err != nil {
			fmt.Fprintf(w, failure, err)
			return
		}
		fmt.Fprint(w, string(data))
	} else {
		fmt.Fprint(w, success)
	}
}

func main() {
	index = flag.String("index", "", "The html file to use as an index")
	endpoint := flag.String("endpoint", "app.go.js", "The name of the compiled javascript file")
	host := flag.String("http", "localhost:8080", "The host at which to serve")

	banner := `   _______ _   ___ 
  / __/ _ \ | / (_)
 _\ \/ , _/ |/ / / 
/___/_/|_||___/_/  says... 
`

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, banner)
		fmt.Fprintln(os.Stderr, "List your go files as arguments!")
		flag.PrintDefaults()
	}

	flag.Parse()

	http.HandleFunc("/", buildHandler)
	http.HandleFunc("/"+*endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header()["Content-Type"] = []string{"text/javascript"}
		fmt.Fprint(w, code)
	})

	fmt.Println(banner)
	fmt.Printf("Open your browser to http://%s!\n", *host)

	err := http.ListenAndServe(fmt.Sprintf("%s", *host), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
