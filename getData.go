package main

import "log"

// Bill is a struct that contains all of a bill entry's information
type Bill struct {
	SerialCode  string  `json:"serialCode"`
	CurrentDate string  `json:"currentDate"`
	Image       string  `json:"image"`
	Notes       string  `json:"notes"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

// Bills is a slice of Bill objects
type Bills []Bill

// AllBillEntries gets all entries registered on the database, limited by page, which is the number of elements
// that will be obtained in the end.
func AllBillEntries(pageSize, page int) []Bill {
	query := `
		SELECT serialCode, current_date, notes, latitude, longitude
		FROM billEntries
		ORDER BY current_date
		OFFSET $1 LIMIT $2
	`

	offset := page * pageSize
	rows, err := db.Query(query, offset, pageSize)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	bills := Bills{}
	for rows.Next() {
		bill := Bill{}
		err := rows.Scan(
			&bill.SerialCode,
			&bill.CurrentDate,
			&bill.Notes,
			&bill.Latitude,
			&bill.Longitude,
		)
		if err != nil {
			log.Println(err)
			return nil
		}
		bills = append(bills, bill)
	}
	return bills
}

// BillEntries obtains a bill's entries through its serial code, amount of entries is limited by page.
func BillEntries(serialCode string, pageSize, page int) []Bill {
	query := `
		SELECT serialCode, current_date, notes, latitude, longitude
		FROM billEntries
		WHERE serialCode = $1
		ORDER BY current_date
		OFFSET $2 LIMIT $3
	`

	offset := page * pageSize
	rows, err := db.Query(query, serialCode, offset, pageSize)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	bills := Bills{}
	for rows.Next() {
		bill := Bill{}
		err := rows.Scan(
			&bill.SerialCode,
			&bill.CurrentDate,
			&bill.Notes,
			&bill.Latitude,
			&bill.Longitude,
		)
		if err != nil {
			log.Println(err)
			return nil
		}
		bills = append(bills, bill)
	}
	return bills
}
