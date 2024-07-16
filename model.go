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

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE calls(
			id SERIAL PRIMARY KEY,
			caller TEXT,
			recipient TEXT,
			status TEXT,
			startTime TEXT,
			endTime TEXT
		)
	`)

	if err != nil {
		return err
	}
	return nil
}
