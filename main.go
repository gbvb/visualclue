package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

/*
	The application to be built for this needs to do the following
	* Manage a named collection of images (CRUD)
		* Attach timer to each of the images to be shown as a property to the image
		* order of the image is also part of the collection

	GET / - Lists the collections
	GET /:id/ - list the named collection
	GET /:id/:imageid - list the named image
	POST /:id - Creates an empty collection
	POST /:id/:imageid - Adds :imageid to the list :id (creating ;id when one does not exist)
	PUT /:id/:imageid - replaces the current :imageid
	DELETE /:id/:imageid - deletes a specific image
	DELETE /:id  - deletes the collection


	* provide a url addressable endpoint to see those images in sequence
	* Get those images in sequence

*/

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/", GetListOfCollections).Methods("GET")
	mux.HandleFunc("/{colname}", GetSpecificCollection).Methods("GET")
	mux.HandleFunc("/{colname}/{id:[0-9]+}", GetSpecificImageAndAttributes).Methods("GET")
	mux.HandleFunc("/{colname}", CreateNamedCollection).Methods("POST")
	mux.HandleFunc("/{colname}/{id:[0-9]+}", AddToCollection).Methods("POST")

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}

func GetListOfCollections(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "GetListOfCollections")
}

func GetSpecificCollection(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	colname := vars["colname"]
	fmt.Fprintf(w, "colname : %v", colname)
}

func GetSpecificImageAndAttributes(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	colname := vars["colname"]
	id := vars["id"]
	fmt.Fprintf(w, "colname: %v, id : %v", colname, id)
}

func CreateNamedCollection(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	colname := vars["colname"]
	fmt.Fprintf(w, "colname : %v", colname)
}

type image_attributes struct {
	Name  string
	Timer int
}

func AddToCollection(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	colname := vars["colname"]
	id := vars["id"]
	fmt.Fprintf(w, "colname: %v, id : %v", colname, id)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic("Oy, something happened! ")
	}

	var i image_attributes
	err = json.Unmarshal(body, &i)
	if err != nil {
		panic("Cant parse the json")
	}
}
