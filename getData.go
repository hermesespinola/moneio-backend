package main

import "log"

// Bill is a bill
type Bill struct {
	SerialCode  string  `json:"serialCode"`
	CurrentDate string  `json:"currentDate"`
	Image       string  `json:"image"`
	Notes       string  `json:"notes"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

// Bills are bills
type Bills []Bill

func billEntries(serialCode string, pageSize, page int) []Bill {
	db, err := Connect()
	if err != nil {
		log.Println(err)
		return nil
	}
	defer db.Close()

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
