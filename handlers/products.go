package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/microservices_study/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// To execute GET in WINDOWS
// curl -v localhost:9090 | jq
func (p *Products) GetProducts(rw http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle GET Products")

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

// To execute POST in WINDOWS
// curl -v localhost:9090 -d"{\"id\":4,\"name\":\"Tea\",\"description\":\"hot cup of tea\"}" | jq
func (p *Products) AddProduct(rw http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle POST Products")
	prod := req.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}

// To execute PUT in WINDOWS
// curl -v localhost:9090/2 -XPUT -d"{\"id\":2,\"name\":\"Tea\",\"description\":\"hot cup of tea\"}" | jq
func (p *Products) UpdateProduct(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusInternalServerError)
		return
	}

	p.l.Println("Handle PUT Products", id)

	prod := req.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(req.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		}

		p.l.Printf("Prod: %#v", prod)

		ctx := context.WithValue(req.Context(), KeyProduct{}, prod)
		req = req.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
