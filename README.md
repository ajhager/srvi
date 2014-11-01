# SRVi

SRVi is a utility that eases the development of [GopherJS](http://github.com/gopherjs/gopherjs) applications. It will serve and automatically rebuild your project every time you refresh. Quickly try out new ideas without even needing to supply your own index.html.

## Install

```bash
go get -u github.com/ajhager/srvi
```

## Usage

If a custom index file is supplied, you can either add `<script src="/app.go.js" type="text/javascript"></script>` to your html file, or supply the `-endpoint custom.js` where 'custom' is whatever name you wish to use.

```
   _______ _   ___
  / __/ _ \ | / (_)
 _\ \/ , _/ |/ / /
/___/_/|_||___/_/  says...

List your go files as arguments!
	-endpoint="app.go.js": The name of the compiled javascript file
	-http="localhost:8080": The host at which to serve
	-index="": The html file to use as an index
```
