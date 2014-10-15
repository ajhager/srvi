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

const failure = `<html><head><title>SRVi</title></head><body style="color:#555555;background:#eeeeee;font-family:Arial;font-size:36px;text-align:center;margin-top:80px;">%s</body></html>`

var success = `<html><head><title>SRVi</title></head><body><script src="./main.go.js" type="text/javascript"></script></body></html>`

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

func jsHandler(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers["Content-Type"] = []string{"application/javascript"}
	fmt.Fprint(w, code)
}

func main() {
	index = flag.String("index", "", "The html file to use as an index")
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

	http.HandleFunc("/", buildHandler)
	http.HandleFunc("/main.go.js", jsHandler)

	fmt.Println(banner)
	fmt.Printf("Open your browser to http://%s:%d!\n", *host, *port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
