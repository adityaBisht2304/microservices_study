package handlers

import (
	"log"
	"net/http"

	"github.com/microservices_study/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		p.getProducts(rw, req)
	}

	// For all the other methods
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	// Using Encoder rather than Marshal is a better approach
	// Encoder directly writes to ioWriter (in our case ResponseWriter implements ioWriter)
	// This saves us the need to create a temporary buffer for storing the encoded value like
	// in marshaling
	// Encoding also make the code a litte bit faster

	// d, err := json.Marshal(lp)
	// if err != nil {
	// 	http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	// }
	// rw.Write(d)

	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
