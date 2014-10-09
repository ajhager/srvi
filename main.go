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
	"strings"
)

func programHandler(w http.ResponseWriter, r *http.Request) {
	// Determine the program to build.
	path := r.URL.Path[1:]
	if strings.HasSuffix(r.URL.Path, "/") {
		path += "main"
	}

	// Open the program's source.
	f, err := os.Open(fmt.Sprintf("%s.go", path))
	if err != nil {
		errorTemplate.Execute(w, &Error{err})
		return
	}

	// Compile the program's source.
	var out bytes.Buffer
	err = gopherjslib.Build(f, &out, nil)
	if err != nil {
		errorTemplate.Execute(w, &Error{err})
		return
	}
	script := out.String()

	// Plug in the compiled javascript.
	successTemplate.Execute(w, &Success{template.JS(script)})
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
		fmt.Fprintln(os.Stderr, "Configure me with these flags!")
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
