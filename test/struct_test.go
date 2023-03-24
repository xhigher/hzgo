package test

import (
	"fmt"
	"github.com/xhigher/hzgo/utils"
	"testing"
)

type StructA struct {
	Id string
	Name string
}

type StructB struct {
	*StructA
	Role string
}

func TestStruct(t *testing.T) {
	sb1 := &StructB{
		&StructA{
			Id: "1",
			Name:"aaaa",
		},
		"1",
	}
	sb2 := &StructB{
		&StructA{
			Id: "1",
			Name:"aaaa",
		},
		"1",
	}
	var alist []*StructA
	var blist []*StructB

	alist = append(alist, sb1.StructA)
	alist = append(alist, sb2.StructA)

	blist = append(blist, sb1)
	blist = append(blist, sb2)

	for _,st := range alist {
		fmt.Println("alist", utils.JSONString(st))
	}

	for _,st := range blist {
		fmt.Println("blist", utils.JSONString(st))
	}
}
