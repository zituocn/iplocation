package iplocation

import (
	"fmt"
	"testing"
)

func TestIPSearch(t *testing.T) {
	ip, err := NewIPLocation("/Users/samsong/mygo/src/github.com/zituocn/xuedao/src/template/ip/qqzeng-ip-china-utf8.dat")
	if err != nil {
		t.Errorf(err.Error())
	}
	result := ip.Get("218.88.127.69")
	fmt.Printf("struct: %v \n", result)
	fmt.Printf("ToJSON: %v \n", result.ToJSON())
	fmt.Printf("ToString: %v \n", result.ToString())
}
