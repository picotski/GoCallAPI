package call

import (
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

// Start call
func (c *Call) StartCall() {
	c.Status = "Ongoing"
	c.StartTime = time.Now()
	c.EndTime = time.Time{}
}
