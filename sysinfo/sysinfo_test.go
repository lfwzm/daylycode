package sysinfo

import (
	"fmt"
	"testing"
)

func TestGetSysInfo(t *testing.T) {
	ret, err := GetSysInfo()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("ret is: %+v \n", ret)
}

func TestGetProcessInfoNorMal(t *testing.T) {
	ret, err := GetProcessInfoNorMal()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("ret is : \n %+v \n", ret)
}
