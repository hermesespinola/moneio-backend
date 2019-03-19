package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

// Bill contains the serialCode
type Bill struct {
	serialCode string
}

func connect() (*sql.DB, error) {
	fmtStr := "host=%s port=%s user=%s " +
		"password=%s dbname=%s sslmode=disable"
	psqlInfo := fmt.Sprintf(
		fmtStr,
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DBNAME"),
	)
	db, err := sql.Open("postgres", psqlInfo)
	return db, err
}

func saveBillImage(imFileHeader *multipart.FileHeader, serialCode string) error {
	// Read image
	im, err := imFileHeader.Open()
	defer im.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Ensure dir exists and create final file
	os.MkdirAll("./images/bills/"+serialCode, os.ModePerm)
	file, err := os.Create("./images/bills/" + serialCode + "/" + imFileHeader.Filename)
	defer file.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// write image file to dir
	_, err = io.Copy(file, im)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func uploadBill(serialCode string, latitude, longitude float64, denomination int, notes string, imFileHeader *multipart.FileHeader) error {
	db, err := connect()
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	// Check for bill
	sqlStatement := `
		SELECT COUNT(serialCode), denomination
		FROM bills
		WHERE serialCode = $1
		GROUP BY denomination;
	`
	row := db.QueryRow(sqlStatement, serialCode)
	billExists := 0
	billDenomination := 0
	err = row.Scan(&billExists, &billDenomination)
	switch err {
	case nil:
	case sql.ErrNoRows:
		break
	default:
		log.Println(err.Error())
		return err
	}

	// Insert the new bill
	if billExists == 0 {
		sqlStatement = `
			INSERT INTO bills (serialCode, denomination) VALUES ($1, $2);
		`
		_, err = db.Query(sqlStatement, serialCode, denomination)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else if billDenomination != denomination {
		return errors.New("This bill should be $" + strconv.FormatInt(int64(billDenomination), 10))
	}

	sqlStatement = `
		INSERT INTO billEntry (serialCode, latitude, longitude, notes)
    	VALUES ($1, $2, $3, $4);
	`
	_, err = db.Query(sqlStatement, serialCode, latitude, longitude, notes)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if imFileHeader != nil {
		return saveBillImage(imFileHeader, serialCode)
	}
	return nil
}
