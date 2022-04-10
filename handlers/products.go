// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/microservices_study/data"
)

// A list of products returns in the response
// swagger:response productsResponse
type productsResponseWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// swagger:response noContent
type productsNoContent struct {
}

// swagger:parameters deleteProduct
type productIDParameterWrapper struct {
	// The id of the product to be deleted
	// in: path
	// required: true
	ID int `json:"id"`
}

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// To execute GET in WINDOWS
// curl -v localhost:9090 | jq

// swagger:route GET /products products listOfProducts
// Returns a list of products
// responses:
//	200: productsResponse

// GetProducts returns the products from the data store
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
// curl -v localhost:9090 -d"{\"id\":4,\"name\":\"Tea\",\"description\":\"hot cup of tea\"}"
// curl -v localhost:9090 -d"{\"id\":4,\"name\":\"Tea\",\"description\":\"hot cup of tea\",\"price\":4.00,\"sku\":\"abc-def-fgh\"}"
func (p *Products) AddProduct(rw http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle POST Products")
	prod := req.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}

// To execute PUT in WINDOWS
// curl -v localhost:9090/2 -XPUT -d"{\"id\":2,\"name\":\"Tea\",\"description\":\"hot cup of tea\"}"
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

// curl -v localhost:9090/1 -XDELETE

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a product
// responses:
//	201: noContent

// DeleteProducts delete a product from the database
func (p *Products) DeleteProduct(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Println("Handle Delete Product", id)

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(req.Body)
		if err != nil {
			p.l.Println("[ERROR] Reading Product", err)
			http.Error(rw, "Error Reading Product", http.StatusBadRequest)
			return
		}

		p.l.Printf("Prod: %#v", prod)

		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] Validating Product", err)
			http.Error(
				rw,
				fmt.Sprintf("Error Validating Product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(req.Context(), KeyProduct{}, prod)
		req = req.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
