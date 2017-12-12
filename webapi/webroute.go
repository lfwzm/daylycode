package webapi

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type NodeHandler func(http.ResponseWriter, *http.Request)

type DirComponent interface {
	Add(DirComponent) bool
	Remove(string) bool
	GetChild() ([]DirComponent, bool)
	GetName() string
	Handler(w http.ResponseWriter, req *http.Request)
	IsDir() bool
}

type Node struct {
	Name        string
	Url         string
	NodeHandler NodeHandler
}

func NewNode(name string, handler NodeHandler) *Node {
	return &Node{Name: name, NodeHandler: handler}
}

func (pN *Node) Add(iDirComponent DirComponent) bool {
	return false
}

func (pN *Node) Remove(name string) bool {
	return false
}

func (pN *Node) GetChild() (sziDirComponent []DirComponent, b bool) {
	return
}

func (pN *Node) Handler(w http.ResponseWriter, req *http.Request) {
	pN.NodeHandler(w, req)
}

func (pN *Node) GetName() string {
	return pN.Name
}

func (pN *Node) IsDir() bool {
	return false
}

type Dir struct {
	Name  string
	Sons  []DirComponent
	mutex sync.Mutex
}

func NewDir(name string) *Dir {
	return &Dir{Name: name}
}

func (pD *Dir) Add(iDirComponent DirComponent) bool {
	pD.mutex.Lock()
	for _, dirComponent := range pD.Sons {
		if iDirComponent.GetName() == dirComponent.GetName() {
			return false
		}
	}
	pD.Sons = append(pD.Sons, iDirComponent)
	pD.mutex.Unlock()
	return true
}

func (pD *Dir) Remove(name string) bool {
	pD.mutex.Lock()
	for i, dirComponent := range pD.Sons {
		if name == dirComponent.GetName() {
			pD.Sons = append(pD.Sons[:i], pD.Sons[i+1:]...)
			pD.mutex.Unlock()
			return true
		}
	}
	pD.mutex.Unlock()
	return false
}

//dir in web must have child
func (pD *Dir) GetChild() (dirs []DirComponent, b bool) {
	if len(pD.Sons) == 0 {
		return dirs, false
	}
	return pD.Sons, true
}

func (pD *Dir) GetName() string {
	return pD.Name
}

func (pD *Dir) IsDir() bool {
	return true
}

//return error
func (pD *Dir) Handler(w http.ResponseWriter, req *http.Request) {
	return
}

func AddDir(url string, node DirComponent, root DirComponent) bool {
	if node == nil || root == nil || node.IsDir() {
		return false
	}

	dirs := strings.Split(url, "/")
	dirDeth := len(dirs)
	rootTmp := root

	for i, d := range dirs {
		if i == 0 {
			continue
		}

		if i == dirDeth-1 {
			sons, ok := rootTmp.GetChild()
			if ok {
				for _, dir := range sons {
					if d == dir.GetName() {
						return false
					}
				}
			}
			fmt.Println("bbb")
			rootTmp.Add(node)
			return true
		}
		//获取子节点
		sons, ok := rootTmp.GetChild()
		if !ok {
			newdir := NewDir(d)
			rootTmp.Add(newdir)
			sons, _ = rootTmp.GetChild()
		}

		for _, dir := range sons {
			if d == dir.GetName() {
				rootTmp = dir
				break
			}
		}
	}
	return true
}

func DelDir(url string, root DirComponent) bool {
	if root == nil {
		return false
	}

	dirs := strings.Split(url, "/")
	dirDeth := len(dirs)
	rootTmp := root
	for i, d := range dirs {
		if i == 0 {
			continue
		}

		sons, ok := rootTmp.GetChild()
		if !ok {
			return false
		}

		for _, dir := range sons {
			if d == dir.GetName() {
				if i == dirDeth-1 {
					fmt.Println("last")
					return rootTmp.Remove(d)
				}
				rootTmp = dir
				break
			}
		}

	}

	return false
}

func UpdateDir(url string, node DirComponent, root DirComponent) bool {
	if DelDir(url, root) {
		return AddDir(url, node, root)
	}
	return false
}

func FindNode(url string, root DirComponent) (DirComponent, error) {
	if root == nil {
		return nil, errors.New("root is nil")
	}

	dirs := strings.Split(url, "/")
	dirDeth := len(dirs)
	rootTmp := root
	for i, d := range dirs {
		if i == 0 {
			continue
		}
		sons, ok := rootTmp.GetChild()
		if !ok {
			return nil, errors.New("dir no exist!")
		}

		for _, dir := range sons {
			if d == dir.GetName() {
				if i == dirDeth-1 {

					return dir, nil
				}
				rootTmp = dir
				break
			}
		}
	}
	return nil, errors.New("no find!")
}
