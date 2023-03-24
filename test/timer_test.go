package test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var (
	count = 0
	timer *time.Timer
)

func TestTimer(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	executeTask()


	wg.Wait()

}

func executeTask(){
	fmt.Println("executeTask count=", count)
	timer = time.AfterFunc(5*time.Second, func() {
		count ++
		fmt.Println("AfterFunc count=", count)
		timer = time.AfterFunc(3*time.Second, executeTask)
	})
}

