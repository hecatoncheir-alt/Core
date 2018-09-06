package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"os"
)

type Server struct {
	APIVersion string
	HTTPServer *http.Server
	router     *httprouter.Router
	log        *log.Logger
}

func New(apiVersion string) *Server {
	server := Server{
		APIVersion: apiVersion,
		router:     httprouter.New()}

	logPrefix := fmt.Sprintf("HttpServer ")
	server.log = log.New(os.Stdout, logPrefix, 3)

	return &server
}

func (server *Server) SetUp(staticFilesDirectory, host string, port int) error {
	server.router.NotFound = http.FileServer(http.Dir(staticFilesDirectory))
	server.router.GET("/api/version", server.apiVersionCheckHandler)

	server.HTTPServer = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}

	eventMessage := fmt.Sprintf("Http server listen on %v, port:%v \n", host, port)
	server.log.Println(eventMessage)

	server.HTTPServer.Handler = server.router

	err := server.HTTPServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (server *Server) apiVersionCheckHandler(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	data := map[string]string{"apiVersion": server.APIVersion}
	encodedData, err := json.Marshal(data)
	if err != nil {
		server.log.Println(err)
	}

	response.Header().Set("content-type", "application/javascript")
	_, err = response.Write(encodedData)
	if err != nil {
		server.log.Println(err)
	}
}
