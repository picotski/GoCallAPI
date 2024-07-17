package call

import (
	"database/sql"
	"time"
)

type Call struct {
	ID        int       `json:"id"`
	Caller    string    `json:"caller"`
	Recipient string    `json:"recipient"`
	Status    string    `json:"status"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// Get calls with pagination
func GetCalls(db *sql.DB, page, count int) ([]Call, error) {
	rows, err := db.Query(
		`SELECT id, caller, recipient, status, start_time, end_time 
		FROM calls LIMIT $1 OFFSET $2`,
		count, count*page,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	calls := []Call{}

	for rows.Next() {
		var c Call
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
func (c *Call) GetCall(db *sql.DB) error {
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
func (c *Call) UpdateCall(db *sql.DB) error {
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
func (c *Call) DeleteCall(db *sql.DB) error {
	_, err := db.Exec(
		"DELETE FROM calls WHERE id=$1",
		c.ID,
	)

	return err
}

// Create one call
func (c *Call) CreateCall(db *sql.DB) error {
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

// Start call
func (c *Call) StartCall() {
	c.Status = "Ongoing"
	c.StartTime = time.Now()
	c.EndTime = time.Time{}
}

// Stop call
func (c *Call) StopCall(db *sql.DB) error {
	c.Status = "Ended"
	c.EndTime = time.Now()

	_, err := db.Exec(
		"UPDATE calls SET status=$1, end_time=$2 WHERE id=$3",
		c.Status,
		c.EndTime,
		c.ID,
	)

	return err
}

// Init calls table
func CreateCallTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE calls(
			id SERIAL PRIMARY KEY,
			caller TEXT,
			recipient TEXT,
			status TEXT,
			start_time TIMESTAMP,
			end_time TIMESTAMP
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
