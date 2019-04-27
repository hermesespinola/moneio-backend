package main

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
)

func writeCSV(w io.Writer, serialCode string) {
	var bills Bills
	if serialCode == "" {
		bills = AllBillEntries(1000, 0)
	} else {
		bills = BillEntries(serialCode, 1000, 0)
	}
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	header := []string{"serialCode", "latitude", "longitude", "date", "notes", "image"}
	csvWriter.Write(header)
	for _, bill := range bills {
		record := []string{
			bill.SerialCode,
			strconv.FormatFloat(bill.Latitude, 'f', -1, 64),
			strconv.FormatFloat(bill.Longitude, 'f', -1, 64),
			bill.CurrentDate,
			bill.Notes,
			bill.Image,
		}
		csvWriter.Write(record)
	}
	err := csvWriter.Error()
	if err != nil {
		log.Fatalln(err)
	}
}
