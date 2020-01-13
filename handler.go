package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"simonwaldherr.de/go/caldav-go/handlers"
	"simonwaldherr.de/go/golibs/file"
)

func checkMethod(access accessType, method string) bool {
	switch access {
	case readonly:
		if method == "GET" || method == "PROPFIND" || method == "OPTIONS" {
			return true
		}
	case writeonly:
		if method == "PUT" || method == "PROPFIND" || method == "OPTIONS" {
			return true
		}
	case readwrite:
		return true
	}
	return false
}

func CalDAVHandler(writer http.ResponseWriter, request *http.Request) {
	var udata userData
	var realm string = "Please enter your username and password for this site"

	addr := request.RemoteAddr
	if i := strings.LastIndex(addr, ":"); i != -1 {
		addr = addr[:i]
	}

	fmt.Printf("%s - - [%s] %q %d %d %q %q\n",
		addr,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		fmt.Sprintf("%s %s %s", request.Method, request.URL, request.Proto),
		0,
		0,
		request.Referer(),
		request.UserAgent())

	user, pass, ok := request.BasicAuth()

	if ok {
		log.Printf("Check Authentication %v\n", user)
		udata, ok = userdata[user]
	}

	if (!ok || pass != udata.password) || !checkMethod(udata.access, request.Method) {
		log.Printf("Not authorized user: %v, access type: %v\n", user, udata.access)
		writer.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
		writer.WriteHeader(401)
		writer.Write([]byte("Unauthorised.\n"))
		return
	}

	response := HandleCalDAVRequest(request)
	response.Write(writer)
}

func FeedHandler(writer http.ResponseWriter, request *http.Request) {
	var realm string = "Please enter your username and password for this site"

	addr := request.RemoteAddr
	if i := strings.LastIndex(addr, ":"); i != -1 {
		addr = addr[:i]
	}
	fmt.Printf("%s - - [%s] %q %d %d %q %q\n",
		addr,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		fmt.Sprintf("%s %s %s", request.Method, request.URL, request.Proto),
		0,
		0,
		request.Referer(),
		request.UserAgent())

	if udata, ok := userdata[""]; !ok {
		user, pass, ok := request.BasicAuth()

		if ok {
			log.Printf("Check Authentication %v\n", user)
			udata, ok = userdata[user]
		}

		if (!ok || pass != udata.password) || !checkMethod(udata.access, request.Method) {
			log.Printf("Not authorized user: %v, access type: %v\n", user, udata.access)
			writer.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			writer.WriteHeader(401)
			writer.Write([]byte("Unauthorised.\n"))
			return
		}
	}

	writer.Write([]byte("BEGIN:VCALENDAR\n"))
	file.Each(StoragePath, true, func(filename, extension, filepath string, dir bool, fileinfo os.FileInfo) {
		if extension == "ics" && !dir {
			str, _ := file.Read(filepath)
			str = strings.Replace(str, "BEGIN:VCALENDAR", "", 1)
			str = strings.Replace(str, "END:VCALENDAR", "", 1)
			writer.Write([]byte(str))
		}
	})
	writer.Write([]byte("END:VCALENDAR\n"))
}

func HandleCalDAVRequest(request *http.Request) *handlers.Response {
	handler := handlers.NewHandler(request)
	return handler.Handle()
}
