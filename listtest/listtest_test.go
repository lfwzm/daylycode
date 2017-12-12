package listtest

import (
	"container/list"
	"fmt"
	"testing"
)

func TestEnv(t *testing.T) {
	l := list.New()
	e4 := l.PushBack(4)
	e1 := l.PushFront(1)
	l.InsertBefore("abc", e4)
	l.InsertAfter(2, e1)

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	var inf interface{}
	inf = "abc"
	fmt.Println(inf)
}
list.List

func BenchEnv(b *testing.B) {

}
