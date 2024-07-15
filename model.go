package main

import (
	"database/sql"
	"errors"
)

type call struct {
	ID     int    `json:"id"`
	Caller string `json:"caller"`
	Status string `json:"status"`
}

func (c *call) getCall(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, caller, status FROM calls WHERE id=$1",
		c.ID,
	).Scan(&c.ID, &c.Caller, &c.Status)
}

func (c *call) updateCall(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (c *call) deleteCall(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (c *call) createCall(db *sql.DB) error {
	return errors.New("Not implemented")
}
