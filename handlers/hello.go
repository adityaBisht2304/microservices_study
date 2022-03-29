package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.l.Println("Hello World")
	d, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "OOPS", http.StatusBadRequest)

		// Same thing can be done through the below lines as well
		// rw.WriteHeader(http.StatusBadRequest)
		// rw.Write([]byte("OOPS"))
		return
	}
	h.l.Printf("Data received is %s\n", d)
	fmt.Fprintf(rw, "Hello %s", d)
}
