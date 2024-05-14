package logic

import (
	"fmt"
	"github.com/xhigher/hzgo/utils"
	"testing"
)

func TestUsersIds(t *testing.T) {

	for i := 0; i < 10; i++ {
		id := utils.IntToBase36(utils.NowTimeMillis() - 888999000000 + int64(i))
		name := utils.RandLetterString(20)

		fmt.Println("id=", id, "name="+name)
		//time.Sleep(1 * time.Second)
	}
}
