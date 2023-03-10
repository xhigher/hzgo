package logic

import (
	"time"
)


type Engine struct {
	ticker   *time.Ticker
}

func newEngine() *Engine{
	return &Engine{
		ticker: time.NewTicker(tickerDuration),
	}
}