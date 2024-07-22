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
func GetCall(db *sql.DB, id int) (Call, error) {
	c := Call{ID: id}

	err := db.QueryRow(
		"SELECT id, caller, recipient, status, start_time, end_time FROM calls WHERE id=$1",
		id,
	).Scan(
		&c.ID,
		&c.Caller,
		&c.Recipient,
		&c.Status,
		&c.StartTime,
		&c.EndTime,
	)

	return c, err
}

// Delete one call
func DeleteCall(db *sql.DB, id int) (Call, error) {
	_, err := db.Exec(
		"DELETE FROM calls WHERE id=$1",
		id,
	)

	c := Call{ID: id}

	return c, err
}

// Create one call
func CreateCall(db *sql.DB, c Call) (Call, error) {
	call := Call{}

	err := db.QueryRow(
		"INSERT INTO calls(caller, recipient, status, start_time, end_time) VALUES($1, $2, $3, $4, $5) RETURNING *",
		c.Caller,
		c.Recipient,
		c.Status,
		c.StartTime,
		c.EndTime,
	).Scan(
		&call.ID,
		&call.Caller,
		&call.Recipient,
		&call.Status,
		&call.StartTime,
		&call.EndTime,
	)

	return call, err
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
func StopCall(db *sql.DB, id int) (Call, error) {
	_, err := db.Exec(
		"UPDATE calls SET status=$1, end_time=$2 WHERE id=$3",
		"Ended",
		time.Now(),
		id,
	)

	if err != nil {
		return Call{}, err
	}

	call, err := GetCall(db, id)

	return call, err
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
