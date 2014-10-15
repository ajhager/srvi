// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/go-on/gopherjslib"
	"html/template"
	"net/http"
	"os"
	"path"
)

var failureT *template.Template

const failure = `<html><head><title>SRVi</title></head><body style="color:#555555;background:#eeeeee;font-family:Arial;font-size:36px;text-align:center;margin-top:80px;">{{.}}</body></html>`

var successT *template.Template

const success = `<html><head><title>SRVi</title></head><body><script>{{.}}</script></body></html>`

func programHandler(w http.ResponseWriter, r *http.Request) {
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
			failureT.Execute(w, err)
			return
		}
		defer file.Close()

		builder.Add(name, file)
	}

	if err := builder.Build(); err != nil {
		failureT.Execute(w, err)
		return
	}

	successT.Execute(w, template.JS(out.String()))
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	failureT, _ = template.New("path").Parse(failure)
	successT, _ = template.New("path").Parse(success)

	static := flag.String("static", "data", "The relative path to your assets")
	host := flag.String("host", "127.0.0.1", "The host at which to serve")
	port := flag.Int("port", 8080, "The port at which to serve")

	banner := `   _______ _   ___ 
  / __/ _ \ | / (_)
 _\ \/ , _/ |/ / / 
/___/_/|_||___/_/  says... 
`

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, banner)
		fmt.Fprintln(os.Stderr, "List all go files as arguments!")
		flag.PrintDefaults()
	}

	flag.Parse()

	http.HandleFunc("/", programHandler)
	http.HandleFunc(fmt.Sprintf("/%s/", path.Clean(*static)), staticHandler)

	fmt.Println(banner)
	fmt.Printf("Open your browser to http://%s:%d!\n", *host, *port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
