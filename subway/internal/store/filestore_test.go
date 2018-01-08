package store

import (
	"fmt"
	"os"
	"testing"
)

func Test_NewFileStore(t *testing.T) {
	pwd, _ := os.Getwd()
	if newFileStore := NewFileStore("", pwd); newFileStore != nil {
		t.Fatal("empty topic failed")
	}
	newFileStore := NewFileStore("test", pwd)
	if newFileStore == nil {
		t.Fatal("NewFileStore failed")
	}

	if newFileStore.baseDir != pwd {
		t.Fatal("NewFileStore add pwd failed")
	}

	if newFileStore.Topic != "test" {
		t.Fatal("NewFileStore add topic failed")
	}

	if newFileStore.curentfd != nil {
		t.Fatal("NewFileStore new curentfd error")
	}

	if newFileStore.maxFileNum != 1 {
		t.Fatal("NewFileStore new maxFileNum error")
	}

	if newFileStore.minFileNum != 0 {
		t.Fatal("NewFileStore new minFileNum error")
	}

	if newFileStore.msgnumber != 0 {
		t.Fatal("NewFileStore new msgnumber error")
	}
}

func Test_Init(t *testing.T) {
	pwd, _ := os.Getwd()
	newFileStore := NewFileStore("test", pwd)
	if newFileStore == nil {
		t.Fatal("NewFileStore failed")
	}

	err := newFileStore.Init()
	if err != nil {
		t.Fatal(err)
	}
	if newFileStore.curentfd == nil {
		t.Fatal("After Init curentfd is nil")
	}

	if newFileStore.curentfd.Name() != pwd+"/"+"test1" {
		fmt.Println(newFileStore.curentfd.Name())
		t.Fatal("After Init createfileName error")
	}

	if newFileStore.maxFileNum != 1 {
		t.Fatal("After Init maxFileNum error")
	}

	if newFileStore.maxIndex != 0 {
		t.Fatal("After Init maxIndex error")
	}
}

func Test_StoreMsg(t *testing.T) {
	pwd, _ := os.Getwd()
	newFileStore := NewFileStore("test", pwd)
	if newFileStore == nil {
		t.Fatal("NewFileStore failed")
	}

	err := newFileStore.Init()
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 4096)
	for i := 0; i < 4096; i++ {
		buf[i] = 'a'
	}

	for j := 0; j < 1024*256*2; j++ {
		_, err := newFileStore.StoreMsg(buf, "test")
		if err != nil {
			panic(err)
		}
	}

	if newFileStore.maxFileNum != 2 {
		t.Fatal("After StoreMsg maxFileNum error")
	}

	if newFileStore.maxIndex != 1024*1024*1024 {
		t.Fatal("After StoreMsg maxIndex error")
	}

	if newFileStore.msgnumber != 1024*256*2 {
		t.Fatal("After StoreMsg msgnumber error")
	}

	index := newFileStore.GetMaxMsgnumber()
	if index.Index != 1024*1024*1024 {
		t.Fatal("After GetMaxMsgnumber get Index error")
	}

	if index.MsgNumber != 1024*256*2 {
		t.Fatal("After GetMaxMsgnumber get MsgNumber error")
	}

	if index.MsgFileNumber != 2 {
		t.Fatal("After GetMaxMsgnumber get MsgFileNumber error")
	}

	if index.Len != 0 {
		t.Fatal("After GetMaxMsgnumber get Len error")
	}

	err = newFileStore.SetMinMsgFileNum(100)
	if err == nil {
		t.Fatal("After SetMinMsgFileNum set too long ok")
	}

	err = newFileStore.SetMinMsgFileNum(-1)
	if err == nil {
		t.Fatal("After SetMinMsgFileNum set too short ok")
	}

	err = newFileStore.SetMinMsgFileNum(1)
	if err != nil || newFileStore.minFileNum != 1 {

		t.Fatal("After SetMinMsgFileNum error")
	}

	newFileStore.Close()
	if newFileStore.curentfd != nil {
		t.Fatal("After Close error")
	}
}

/*
func Test_StoreMsg(t *testing.T) {
	pwd, _ := os.Getwd()
	newFileStore := NewFileStore("test", pwd)
	if newFileStore == nil {
		t.Fatal("NewFileStore failed")
	}

	err := newFileStore.Init()
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 4096)
	for i := 0; i < 4096; i++ {
		buf[i] = 'a'
	}

	for j := 0; j < 1024*256*2; j++ {
		_, err := newFileStore.StoreMsg(buf, "test")
		if err != nil {
			panic(err)
		}
	}

	if newFileStore.maxFileNum != 6 {
		t.Fatal("After StoreMsg maxFileNum error")
	}

	if newFileStore.maxIndex != 1024*1024*1024 {
		t.Fatal("After StoreMsg maxIndex error")
	}

	if newFileStore.msgnumber != 1024*256*2 {
		t.Fatal("After StoreMsg msgnumber error")
	}
}
*/
