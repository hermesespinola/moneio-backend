package main

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func saveBillImage(imFileHeader *multipart.FileHeader, serialCode, id string) error {
	// Read image
	im, err := imFileHeader.Open()
	defer im.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Ensure dir exists and create final file
	os.MkdirAll("./images/bills/"+serialCode, os.ModePerm)
	file, err := os.Create("./images/bills/" + serialCode + "/" + id)
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

// UploadBill uploads a new bill entry to the database. If there is no previous entry of a bill, it creates a new bill row to keep track
// of the number of entries, it also creates the new entry for the bill.
func UploadBill(serialCode string, latitude, longitude float64, denomination int, notes string, imFileHeader *multipart.FileHeader) error {
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
	err := row.Scan(&billExists, &billDenomination)
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
		INSERT INTO billEntries (serialCode, latitude, longitude, notes)
    	VALUES ($1, $2, $3, $4);
	`
	_, err = db.Query(sqlStatement, serialCode, latitude, longitude, notes)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	sqlStatement = `
		SELECT id
		FROM billEntries
		WHERE serialCode = $1
		ORDER BY creationDate ASC
		LIMIT 1;
	`
	row = db.QueryRow(sqlStatement, serialCode)
	var id string
	err = row.Scan(&id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if imFileHeader != nil {
		return saveBillImage(imFileHeader, serialCode, id)
	}
	return nil
}
