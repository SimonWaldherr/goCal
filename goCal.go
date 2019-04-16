package main

import (
	"net/http"

	"simonwaldherr.de/go/golibs/file"
	"simonwaldherr.de/go/gwv"
)

const (
	SSLCert   = "ssl.cert"
	SSLKey    = "ssl.key"
	HTTPPort  = 8080
	HTTPSPort = 443
)

func main() {
	if !file.Exists("ssl.cert") || !file.Exists("ssl.key") {
		options := map[string]string{}
		options["certPath"] = SSLCert
		options["keyPath"] = SSLKey
		options["host"] = "*"
		options["countryName"] = "DE"
		options["provinceName"] = "Bavaria"
		options["organizationName"] = "Lorem Ipsum Ltd"
		options["commonName"] = "*"

		gwv.GenerateSSL(options)
	}

	startServer()
}

func handler(rw http.ResponseWriter, req *http.Request) (string, int) {
	RequestHandler(rw, req)
	return "", 0
}

func startServer() {
	HTTPD := gwv.NewWebServer(HTTPPort, 60)
	HTTPD.ConfigSSL(HTTPSPort, SSLKey, SSLCert, true)

	HTTPD.URLhandler(
		gwv.URL("^/", handler, gwv.MANUAL),
	)

	HTTPD.Start()
	HTTPD.WG.Wait()
}
