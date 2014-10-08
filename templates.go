// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "html/template"

type Success struct {
	Script template.JS
}

type Error struct {
	Error error
}

var errorTemplate *template.Template
var successTemplate *template.Template

const errorPage = `
<!DOCTYPE html>
<html>
	<head>
		<title>SRVi</title>
		<link rel="icon" type="image/png" href="/favicon.png">
		<style>
		html, body {
			padding: 0;
			margin: 0;
			font-family: Arial;
			background: #f4f4f4;
			width: 100%;
			height: 100%;
			overflow: hidden;
			font-size: 18px;
		}
		div#error {
			color: #4b5464;
			font-size: 1.75em;
			text-align: center;
			padding: 80px;
		}
		</style>
	</head>
  <body>
		<div id="error">{{.Error}}</div>
  </body>
</html>
`

const successPage = `
<!DOCTYPE html>
<html>
	<head>
		<title>SRVi</title>
		<link rel="icon" type="image/png" href="/favicon.png">
	</head>
  <body>
		<script>
			{{.Script}}
		</script>
  </body>
</html>
`

func init() {
	errorTemplate, _ = template.New("path").Parse(errorPage)
	successTemplate, _ = template.New("path").Parse(successPage)
}
