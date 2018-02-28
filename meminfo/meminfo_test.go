package meminfo

import (
	"fmt"
	"testing"
)

func Test_GetMemInfo(t *testing.T) {
	ret, err := GetMemInfo()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("ret is: %+v \n", ret)
}
