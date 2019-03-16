package main

import (
	"fmt"
	"log"
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
	coords := r.FormValue("coords")
	imFileHeader := r.MultipartForm.File["image"][0]

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	uploadBill(serialCode, denomination, coords, notes, imFileHeader)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func main() {
	http.HandleFunc("/upload-bill", billHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
