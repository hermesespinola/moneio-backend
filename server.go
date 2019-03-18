package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func billHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data
	err := r.ParseMultipartForm(0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		fmt.Fprint(w, "Error parsing multiform data")
		return
	}

	serialCode := r.FormValue("serialCode")
	log.Println(serialCode)
	denomination, err := strconv.Atoi(r.FormValue("denomination"))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	notes := r.FormValue("notes")
	latitude, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	longitude, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	imArray := r.MultipartForm.File["image"]
	var imFileHeader *multipart.FileHeader
	if len(imArray) > 0 {
		imFileHeader = r.MultipartForm.File["image"][0]
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	err = uploadBill(serialCode, latitude, longitude, denomination, notes, imFileHeader)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "Ok")
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		h.ServeHTTP(w, r)
	})
}

func main() {
	// Set custom env variable
	os.Setenv("PG_HOST", "127.0.0.1")
	os.Setenv("PG_PORT", "5432")
	os.Setenv("PG_USER", "postgres")
	os.Setenv("PG_PASSWORD", "postgres")
	os.Setenv("PG_DBNAME", "moneio")

	mux := http.NewServeMux()
	mux.HandleFunc("/upload-bill", billHandler)
	server := cors(mux)
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", server))

}
