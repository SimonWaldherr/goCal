package main

import (
	"flag"
	"net/http"

	"simonwaldherr.de/go/caldav-go/files"
	"simonwaldherr.de/go/golibs/file"
	"simonwaldherr.de/go/gwv"
)

var (
	SSLCert     string
	SSLKey      string
	UserFile    string
	StoragePath string
)

const (
	HTTPPort  = 8080
	HTTPSPort = 443
)

var userdata map[string]userData

func init() {
	flag.StringVar(&SSLCert, "sslcrt", "ssl.crt", "path to the ssl/tls crt file")
	flag.StringVar(&SSLKey, "sslkey", "ssl.key", "path to the ssl/tls key file")
	flag.StringVar(&UserFile, "user", "user.csv", "path to the user.csv file")
	flag.StringVar(&StoragePath, "storage", "icsdata", "path to the folder with ics data")
	flag.Parse()

	userdata = loadUserDataFromFile(UserFile)
	files.StoragePath = StoragePath
}

func main() {
	if !file.Exists(SSLCert) || !file.Exists(SSLKey) {
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

func caldavhandler(rw http.ResponseWriter, req *http.Request) (string, int) {
	CalDAVHandler(rw, req)
	return "", 0
}

func feedhandler(rw http.ResponseWriter, req *http.Request) (string, int) {
	FeedHandler(rw, req)
	return "", 0
}

func startServer() {
	HTTPD := gwv.NewWebServer(HTTPPort, 60)
	HTTPD.ConfigSSL(HTTPSPort, SSLKey, SSLCert, true)

	HTTPD.URLhandler(
		gwv.URL("^/icsfeed/", feedhandler, gwv.MANUAL),
		gwv.URL("^/", caldavhandler, gwv.MANUAL),
	)

	HTTPD.Start()
	HTTPD.WG.Wait()
}
