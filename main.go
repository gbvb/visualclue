package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

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
	DELETE /:id/:imageid - deletes a specific image from the collection
	DELETE /:id  - deletes the collection
	POST /images/:imageid - add an image
	GET /images/:imageid - retrieve the image


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
	mux.HandleFunc("/images/{imageid:[0-9]+}", AddImage).Methods("PUT").Headers("Content-Type", "multipart/form-data")

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
	Name    string //Name for the image
	Time    int    //Time in milliseconds on this image. -1 to not proceed next
	ImageId int    //image id
}

//Add an image as an upload
func AddImage(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	image_id := vars["imageid"]
	fmt.Fprintf(w, "imageid: %v", image_id)

	file, _, err := req.FormFile("file")

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	out, err := os.Create("/tmp/img")
	if err != nil {
		fmt.Fprintf(w, "unable to create file for writing")
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	return
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
