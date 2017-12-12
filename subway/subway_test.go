package subway

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"
	"unsafe"
)

type SliceHeader struct {
	addr uintptr
	len  int
	cap  int
}

func TestHandleConn(t *testing.T) {
	server := NewSubWayStation()
	server.Ip = "127.0.0.1"
	server.Port = "9090"
	server.worker = &NormalHandler{}
	err := server.GetListener()
	if err != nil {
		t.Fail()
	}
	go server.AcceptAndHandle()
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", server.Ip+":"+server.Port)
	if err != nil {
		t.Fatal(err)
	}

	head := &XpHead{}
	head.IputNumber = 1
	head.IsGiz = 2
	Len := unsafe.Sizeof(*head)

	newhead := &SliceHeader{
		addr: uintptr(unsafe.Pointer(head)),
		cap:  int(Len),
		len:  int(Len),
	}
	fmt.Println(newhead)
	data := *(*[]byte)(unsafe.Pointer(newhead))
	fmt.Println("datalen is: ", len(data))

	var body = make(map[string]interface{}, 1024)
	body["Name"] = "wuzhiming"
	body["Phone"] = "18826418902"
	bdata, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//copy(data2[len(data):], bdata)

	var data2 = make([]byte, 4096)
	var head2 *XpHead = *(**XpHead)(unsafe.Pointer(&data2))
	head2.Flag = 1
	head2.Id = 1024
	head2.IputNumber = 10000
	head2.IsFirst = 1
	tolist := []byte("mytest")
	listenlist := []byte("mytest")
	copy(head2.StrPushList[:], tolist)
	copy(head2.StrListenList[:], listenlist)
	head2.BodyLen = int64(copy(data2[len(data):], bdata))
	fmt.Println("head is: ", head2)

	Len2, _ := conn.Write(data2)
	fmt.Println("Len()TEST  is: ", Len2)
	time.Sleep(1 * time.Second)

	data3 := make([]byte, 4096)
	fmt.Println("before read data")
	_, err = conn.Read(data3)
	if err != nil {
		fmt.Println(err)
		return
	}

	bodyret := data3[len(data):]
	fmt.Println("+++++++++bodyret is: ", string(bodyret))
}
