package gziptest

import (
	"os"
	"testing"
)

func TestCompressFile(t *testing.T) {
	pwd, _ := os.Getwd()
	newfile, err := os.Create(pwd + "/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	newfile.Write([]byte("hello world!!!!"))
	newfile.Close()

	err = CompressFile(pwd+"/test.gz", pwd+"/test.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeCompressFile(t *testing.T) {
	pwd, _ := os.Getwd()

	err := DeCompressFile(pwd+"/test2.txt", pwd+"/test.gz")
	if err != nil {
		t.Fatal(err)
	}
}
