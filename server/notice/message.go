package notice

import (
	"time"
)

type Message struct {
	Type      string
	From      string
	To        string
	Data      []byte
	Timestamp time.Time
}
