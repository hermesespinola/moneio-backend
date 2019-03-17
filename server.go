package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

func billHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data
	err := r.ParseMultipartForm(0)
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error parsing multiform data")
		return
	}

	serialCode := r.FormValue("serialCode")
	denomination, err := strconv.Atoi(r.FormValue("denomination"))
	notes := r.FormValue("notes")
	latitude, err1 := strconv.ParseFloat(r.FormValue("latitude"), 64)
	longitude, err2 := strconv.ParseFloat(r.FormValue("longitude"), 64)
	imArray := r.MultipartForm.File["image"]
	var imFileHeader *multipart.FileHeader
	if len(imArray) > 0 {
		imFileHeader = r.MultipartForm.File["image"][0]
		if err != nil || err1 != nil || err2 != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	err = uploadBill(serialCode, latitude, longitude, denomination, notes, imFileHeader)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error")
	}
	fmt.Fprintf(w, "Ok")
}

func main() {
	http.HandleFunc("/upload-bill", billHandler)
	println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
