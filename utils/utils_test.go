package utils

import (
	"fmt"
	"testing"
)

func TestIntToBase36(t *testing.T) {

	micro := NowTimeMicro()-1040000000000000

	code := IntToBase36(micro)
	fmt.Println("TestIntToBase36: ", micro, code)
}
