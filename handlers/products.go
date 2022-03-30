package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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
		return
	}

	if req.Method == http.MethodPost {
		p.addProduct(rw, req)
		return
	}

	if req.Method == http.MethodPut {
		p.l.Println("PUT")
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(req.URL.Path, -1)
		if len(g) != 1 {
			p.l.Println("Invalid URL : More than one id")
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URL : Unable to convert to number")
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}
		p.l.Println("id : ", id)

		p.updateProduct(id, rw, req)
	}

	// For all the other methods
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// To execute GET in WINDOWS
// curl -v localhost:9090 | jq
func (p *Products) getProducts(rw http.ResponseWriter, req *http.Request) {
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
func (p *Products) addProduct(rw http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle POST Products")
	prod := &data.Product{}

	err := prod.FromJSON(req.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	p.l.Printf("Prod: %#v", prod)

	data.AddProduct(prod)
}

// To execute PUT in WINDOWS
// curl -v localhost:9090/2 -XPUT -d"{\"id\":2,\"name\":\"Tea\",\"description\":\"hot cup of tea\"}" | jq
func (p *Products) updateProduct(id int, rw http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle PUT Products")
	prod := &data.Product{}

	err := prod.FromJSON(req.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	p.l.Printf("Prod: %#v", prod)

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
