package main

import(
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request){
		log.Println("Hello World")
		d, err := ioutil.ReadAll(req.Body)		
		if err != nil {
			http.Error(rw, "OOPS", http.StatusBadRequest)

			// Same thing can be done through the below lines as well
			// rw.WriteHeader(http.StatusBadRequest)
			// rw.Write([]byte("OOPS"))
			return
		}
		log.Printf("Data received is %s\n", d)
		fmt.Fprintf(rw, "Hello %s", d)
	})

	// Listen on port 9090 for any ip address
	http.ListenAndServe(":9090", nil)

	// To mention a specific address - in this case it is local host
	//http.ListenAndServe("127.0.0.1:9090", nil)


}
