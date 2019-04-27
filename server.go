package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func postBill(w http.ResponseWriter, r *http.Request) {
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

	err = UploadBill(serialCode, latitude, longitude, denomination, notes, imFileHeader)
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

func getAllBillEntries(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	pageSizeArr := queryValues["pageSize"]
	pageSize := 100 // Default page size
	var err error
	if len(pageSizeArr) > 0 {
		pageSize, err = strconv.Atoi(pageSizeArr[0])
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
	}
	pageArr := queryValues["page"]
	page := 0
	if len(pageArr) > 0 {
		page, err = strconv.Atoi(pageArr[0])
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
	}
	bills := AllBillEntries(pageSize, page)
	billsJSON, err := json.Marshal(bills)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}
	fmt.Fprintf(w, "%s", billsJSON)
}

func getBillEntries(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	vars := mux.Vars(r)
	serialCode := vars["serialCode"]
	pageSizeArr := queryValues["pageSize"]
	pageSize := 100 // Default page size
	var err error
	if len(pageSizeArr) > 0 {
		pageSize, err = strconv.Atoi(pageSizeArr[0])
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
	}
	pageArr := queryValues["page"]
	page := 0
	if len(pageArr) > 0 {
		page, err = strconv.Atoi(pageArr[0])
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
	}
	bills := BillEntries(serialCode, pageSize, page)
	billsJSON, err := json.Marshal(bills)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}
	fmt.Fprintf(w, "%s", billsJSON)
}

func getCSV(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if serialCode, ok := vars["serialCode"]; ok {
		writeCSV(w, serialCode)
	} else {
		writeCSV(w, "")
	}
}

func main() {
	defer db.Close()
	r := mux.NewRouter()
	r.HandleFunc("/uploadBill", postBill)
	r.HandleFunc("/billEntries/{serialCode}", getBillEntries)
	r.HandleFunc("/billEntries", getAllBillEntries)
	r.HandleFunc("/csv", getCSV)
	r.HandleFunc("/csv/{serialCode}", getCSV)
	server := cors(r)
	port := os.Getenv("PORT")
	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, server))
}
