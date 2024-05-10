package notice

import (
	"time"
)

type Message struct {
	Type      string
	From      string
	To        string
	Message   string
	Timestamp time.Time
}
