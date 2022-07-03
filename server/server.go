package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	port uint16
}

//create server
func NewServer(port uint16) *Server {
	return &Server{port: port}
}

//return port
func (sv *Server) Port() uint16 {
	return sv.port
}

func helloworld(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello world")
}

func (sv *Server) Run() {
	http.HandleFunc("/", helloworld)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(sv.Port())), nil))
}
