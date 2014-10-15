// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/go-on/gopherjslib"
	"html/template"
	"net/http"
	"os"
	"path"
)

func programHandler(w http.ResponseWriter, r *http.Request) {
	var out bytes.Buffer
	builder := gopherjslib.NewBuilder(&out, nil)

	for _, name := range flag.Args() {
		if path.Ext(name) != ".go" {
			continue
		}

		file, err := os.Open(name)
		if err != nil {
			errorTemplate.Execute(w, &Error{err})
			return
		}
		defer file.Close()

		builder.Add(name, file)
	}

	if err := builder.Build(); err != nil {
		errorTemplate.Execute(w, &Error{err})
		return
	}

	successTemplate.Execute(w, &Success{template.JS(out.String())})
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
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
	http.Handle("/favicon.png", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "."}))
	http.HandleFunc(fmt.Sprintf("/%s/", path.Clean(*static)), staticHandler)

	fmt.Println(banner)
	fmt.Printf("Open your browser to http://%s:%d!\n", *host, *port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
