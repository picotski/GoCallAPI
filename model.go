package main

import (
	"database/sql"
)

type call struct {
	ID        int    `json:"id"`
	Caller    string `json:"caller"`
	Recipient string `json:"recipient"`
	Status    string `json:"status"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// Get calls with pagination
func getCalls(db *sql.DB, page, count int) ([]call, error) {
	rows, err := db.Query(
		`SELECT id, caller, recipient, status, start_time, end_time 
		FROM calls LIMIT $1 OFFSET $2`,
		count, count*page,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	calls := []call{}

	for rows.Next() {
		var c call
		if err := rows.Scan(
			&c.ID,
			&c.Caller,
			&c.Recipient,
			&c.Status,
			&c.StartTime,
			&c.EndTime,
		); err != nil {
			return nil, err
		}
		calls = append(calls, c)
	}
	return calls, nil
}

// Get one call
func (c *call) getCall(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, caller, recipient, status, start_time, end_time FROM calls WHERE id=$1",
		c.ID,
	).Scan(
		&c.ID,
		&c.Caller,
		&c.Recipient,
		&c.Status,
		&c.StartTime,
		&c.EndTime,
	)
}

// Update one call
func (c *call) updateCall(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE calls SET caller=$1, recipient=$2, status=$3, start_time=$4, end_time=$5 WHERE id=$6",
		c.Caller,
		c.Recipient,
		c.Status,
		c.StartTime,
		c.EndTime,
		c.ID,
	)

	return err
}

// Delete one call
func (c *call) deleteCall(db *sql.DB) error {
	_, err := db.Exec(
		"DELETE FROM calls WHERE id=$1",
		c.ID,
	)

	return err
}

// Create one call
func (c *call) createCall(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO calls(caller, recipient, status, start_time, end_time) VALUES($1, $2, $3, $4, $5) RETURNING id",
		c.Caller,
		c.Recipient,
		c.Status,
		c.StartTime,
		c.EndTime,
	).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

// Get amount of calls
func CountCalls(db *sql.DB) (int, error) {
	rows, err := db.Query("SELECT COUNT(*) FROM calls")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}

// Init calls table
func CreateCallTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE calls(
			id SERIAL PRIMARY KEY,
			caller TEXT,
			recipient TEXT,
			status TEXT,
			start_time TEXT,
			end_time TEXT
		)
	`)

	if err != nil {
		return err
	}
	return nil
}

// Drop calls table
func DeleteCallTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE calls")

	if err != nil {
		return err
	}
	return nil
}
