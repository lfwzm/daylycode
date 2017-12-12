package webapi

import (
	"fmt"
	"net/http"
	"testing"
)

func THandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("ok")
}

func TestNewNode(t *testing.T) {
	Node := NewNode("TestNode", THandler)
	if Node.GetName() != "TestNode" {
		t.Fail()
	}
}

func TestNewDir(t *testing.T) {
	Dir := NewDir("TestDir")
	if Dir.GetName() != "TestDir" {
		t.Fail()
	}
}

func TestDirAdd(t *testing.T) {
	Dir := NewDir("TestDir")
	Dir2 := NewDir("SencondDir")
	ok := Dir.Add(Dir2)
	if !ok {
		t.Logf("Dir Add fail")
	}

	ok = Dir.Add(Dir2)
	if ok {
		t.Logf("Dir Add fail")
	}

	sondir, ok := Dir.GetChild()
	if !ok {
		t.Logf("Dir GetChild fail")
	}
	if sondir[0].GetName() != "SencondDir" {
		t.Logf("Dir GetChild fail2")
	}
}

func TestDirRemove(t *testing.T) {
	Dir := NewDir("TestDir")
	Dir2 := NewDir("SencondDir")
	ok := Dir.Add(Dir2)
	if !ok {
		t.Logf("Dir Add fail")
	}

	Node := NewNode("TestNode", THandler)
	ok = Dir.Add(Node)
	if !ok {
		t.Logf("Dir Add fail2")
	}

	ok = Dir.Remove("Null")
	if ok {
		t.Logf("Dir Remove fail")
	}

	ok = Dir.Remove("TestNode")
	if !ok {
		t.Logf("Dir Remove fail2")
	}

	sondir, ok := Dir.GetChild()
	if !ok {
		t.Logf("Dir GetChild fail")
	}
	if sondir[0].GetName() != "SencondDir" {
		t.Logf("Dir GetChild fail2")
	}
}

func TestAddDir1(t *testing.T) {
	dirs := AddDir("/abc/cde/efg", nil, nil)
	if dirs {
		t.Logf("AddDir fail")
	}
}

func TestAddDir2(t *testing.T) {
	root := NewDir("kingdom")
	son1 := NewDir("otc")
	root.Add(son1)
	node := NewNode("wuzhiming", THandler)
	son1.Add(node)

	ok := AddDir("/kingdom/otc/wuzhiming", node, root)
	if ok {
		t.Logf("AddDir fail2")
	}
}

func TestAddDir3(t *testing.T) {
	root := NewDir("kingdom")
	son1 := NewDir("otc")
	root.Add(son1)
	node := NewNode("wuzhiming", THandler)
	//son1.Add(node)

	ok := AddDir("/kingdom/otc/wuzhiming", node, root)
	if !ok {
		t.Logf("AddDir fail3")
	}

	son2, ok := son1.GetChild()
	if !ok {
		t.Logf("AddDir fail3 son1.GetChild failed")
	}

	if son2[0].GetName() != "wuzhiming" {
		t.Logf("AddDir fail3 son2[0].GetName() fail")
	}

}

func TestDelDir(t *testing.T) {
	dir := NewDir("kingdom")
	ok := DelDir("/kingdom/otc/wuzhiming", dir)
	if ok {
		t.Logf("TestDelDir fail")
	}
}

func TestDelDir2(t *testing.T) {
	root := NewDir("kingdom")
	son1 := NewDir("otc")
	root.Add(son1)
	node := NewNode("wuzhiming", THandler)
	son1.Add(node)

	ok := DelDir("/kingdom/otc/wuzhiming", root)
	if !ok {
		t.Logf("TestDelDir2 fail")
	}

	_, ok = son1.GetChild()
	if ok {
		t.Logf("TestDelDir2 fail")
	}

}

func TestFindHandler(t *testing.T) {
	root := NewDir("kingdom")
	son1 := NewDir("otc")
	root.Add(son1)
	node := NewNode("wuzhiming", THandler)
	son1.Add(node)

	_, err := FindNode("/kingdom/otc/wuzhiming", root)
	if err != nil {
		t.Logf("TestFindHandler fail")
	}

}
