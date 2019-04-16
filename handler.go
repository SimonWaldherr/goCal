package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/samedi/caldav-go/data"
	"github.com/samedi/caldav-go/global"
	"github.com/samedi/caldav-go/handlers"
)

var userdata map[string]string

func init() {
	userdata = loadUserDataFromFile("user.csv")
}

func RequestHandler(writer http.ResponseWriter, request *http.Request) {
	var password string
	var realm string = "Please enter your username and password for this site"

	addr := request.RemoteAddr
	if i := strings.LastIndex(addr, ":"); i != -1 {
		addr = addr[:i]
	}
	fmt.Printf("%s - - [%s] %q %d %d %q %q\n",
		addr,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		fmt.Sprintf("%s %s %s", request.Method, request.URL, request.Proto),
		//StatusCode,
		//ContentLength,
		0,
		0,
		request.Referer(),
		request.UserAgent())

	user, pass, ok := request.BasicAuth()

	if ok {
		log.Println("Check Authentication ", user)
		password, ok = userdata[user]
	}

	if !ok || pass != password {
		log.Println("Not authorized ", user)
		writer.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
		writer.WriteHeader(401)
		writer.Write([]byte("Unauthorised.\n"))
		return
	}

	response := HandleRequest(request)
	response.Write(writer)
}

func HandleRequest(request *http.Request) *handlers.Response {
	handler := handlers.NewHandler(request)
	return handler.Handle()
}

func HandleRequestWithStorage(request *http.Request, stg data.Storage) *handlers.Response {
	SetupStorage(stg)
	return HandleRequest(request)
}

func SetupStorage(stg data.Storage) {
	global.Storage = stg
}
