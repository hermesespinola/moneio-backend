package main

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("DBNAME"),
	)
	db, err := sql.Open("postgres", psqlInfo)
	return db, err
}

func saveBillImage(imFileHeader *multipart.FileHeader, id int) (err error) {
	// ensure dir exists and create final file
	os.MkdirAll("./packs", os.ModePerm)
	file, err := os.Create("./packs/" + imFileHeader.Filename)
	defer file.Close()
	if err != nil {
		println(err.Error())
		return err
	}

	// Read image
	im, err := imFileHeader.Open()
	defer im.Close()
	if err != nil {
		println(err.Error())
		return err
	}

	// Ensure dir exists and create final file
	os.MkdirAll("./bills", os.ModePerm)
	file, err = os.Create("./images/" + (string)(id) + "/" + imFileHeader.Filename)
	defer file.Close()
	if err != nil {
		println(err.Error())
		return err
	}

	// write image file to dir
	_, err = io.Copy(file, im)
	if err != nil {
		println(err.Error())
		return err
	}

	return nil
}

func uploadBill(serialCode string, denomination int, coords string, notes string, imFileHeader *multipart.FileHeader) {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `
		INSERT INTO bills (serialCode, denomination, latitude, longitude, notes)
    	VALUES ($1, $2, $3, $4)
		RETURNING serialCode
	`
	id := 0
	err = db.QueryRow(sqlStatement, serialCode, denomination, coords, notes).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)

	saveBillImage(imFileHeader, id)
}
