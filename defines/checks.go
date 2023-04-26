package defines

import (
	"github.com/xhigher/hzgo/consts"
)

var (
	pageLimits = []int32{10, 20, 50, 100}
)

func CheckPageLimit(limit int32) bool {
	for _, l := range pageLimits {
		if l == limit {
			return true
		}
	}
	return false
}

func CheckChangeStatus(status int32) bool {
	if status == consts.StatusOnline || status == consts.StatusOffline {
		return true
	}
	return false
}
