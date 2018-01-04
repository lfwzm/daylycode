package subway

import (
	//"bytes"
	//"encoding/json"
	"fmt"
	//"net"
	"testing"
	"time"
	//"unsafe"
)

type SliceHeader struct {
	addr uintptr
	len  int
	cap  int
}

func Benchmark_SendData(b *testing.B) {
	//var testData []byte = []byte("wuzhiming111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111wuzhiming111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111wuzhiming111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111wuzhiming111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111wuzhiming111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111wuzhiming111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
	var testData []byte = []byte("abcdefg")
	server := NewSubWayStation()
	server.Ip = "127.0.0.1"
	server.Port = "9000"
	server.Worker = &NormalHandler{}
	/*
		err := server.GetListener()
		if err != nil {
			b.Fail()
		}
		go server.AcceptAndHandle()
	*/
	time.Sleep(1 * time.Second)

	passenger := NewPassenger()
	passenger.SetNearStation(server.Ip, server.Port)

	var head2 PassengerHead
	head2.Flag = 1
	head2.Id = 1024
	head2.IputNumber = 100
	head2.ListenTime = 3600 * 2
	head2.IsBroken = 0
	head2.BodyLen = int64(len(testData))
	tolist := []byte("test")
	//listenlist := []byte("yourtest")
	copy(head2.StrPushList[:], tolist)
	//copy(head2.StrListenList[:], listenlist)

	err := passenger.ConnectNearStation()
	if err != nil {
		fmt.Println(err)
	}

	b.ReportAllocs()
	b.N = 1000000
	passenger.SetPassengerHead(&head2)
	for i := 0; i < b.N; i++ { //use b.N for looping
		head2.IputNumber++
		passenger.SetPassengerHead(&head2)
		b.StartTimer()

		err = passenger.SendData(&testData)
		if err != nil {
			b.Fatal(err)
		}
		b.StopTimer()
		if i%100000 == 99999 {
			fmt.Println(i)
		}
	}
}

/*
func TestHandleConn(t *testing.T) {
	server := NewSubWayStation()
	server.Ip = "127.0.0.1"
	server.Port = "9050"
	server.Worker = &NormalHandler{}
	err := server.GetListener()
	if err != nil {
		t.Fail()
	}
	go server.AcceptAndHandle()
	time.Sleep(1 * time.Second)
	return
	conn, err := net.Dial("tcp", server.Ip+":"+server.Port)
	if err != nil {
		t.Fatal(err)
	}

	head := &PassengerHead{}
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
	var head2 *PassengerHead = *(**PassengerHead)(unsafe.Pointer(&data2))
	head2.Flag = 1
	head2.Id = 1024
	head2.IputNumber = 10
	head2.ListenTime = 3600 * 2
	head2.IsBroken = 0
	tolist := []byte("mytest")
	listenlist := []byte("mytest")
	copy(head2.StrPushList[:], tolist)
	copy(head2.StrListenList[:], listenlist)
	head2.BodyLen = int64(copy(data2[len(data):], bdata))
	fmt.Println("head is: ", head2)

	go func() {
		Len2, ok := conn.Write(data2)
		if ok != nil {
			fmt.Println("Len()TEST  is: ", Len2, ok)
		}
		fmt.Println("Len()TEST  is: ", Len2, ok)

	}()
	time.Sleep(1 * time.Second)
	go func() {
		data3 := make([]byte, 4096)
		fmt.Println("before read data")

		_, err = conn.Read(data3)
		if err != nil {
			fmt.Println(err)
			fmt.Println("write end")
			conn.Close()
			return
		}

		bodyret := data3[len(data):]
		fmt.Println("+++++++++bodyret is: ", string(bodyret))
	}()

	time.Sleep(1 * time.Second)
	return
}

func TestSendData(t *testing.T) {
	var testData []byte = []byte("wuzhiming")
	server := NewSubWayStation()
	server.Ip = "127.0.0.1"
	server.Port = "10086"
	server.Worker = &NormalHandler{}
	err := server.GetListener()
	if err != nil {
		t.Fail()
	}
	go server.AcceptAndHandle()
	time.Sleep(1 * time.Second)

	passenger := NewPassenger()
	passenger.SetNearStation(server.Ip, server.Port)

	var head2 PassengerHead
	head2.Flag = 1
	head2.Id = 1024
	head2.IputNumber = 10
	head2.ListenTime = 3600 * 2
	head2.IsBroken = 0
	head2.BodyLen = int64(len(testData))
	tolist := []byte("mytest")
	listenlist := []byte("mytest")
	copy(head2.StrPushList[:], tolist)
	copy(head2.StrListenList[:], listenlist)

	err = passenger.ConnectNearStation()
	if err != nil {
		fmt.Println(err)
	}
	passenger.SetPassengerHead(&head2)

	err = passenger.SendData(&testData)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func TestDisConnectStation(t *testing.T) {
	var testData []byte = []byte("wuzhiming")
	server := NewSubWayStation()
	server.Ip = "127.0.0.1"
	server.Port = "8891"
	server.Worker = &NormalHandler{}
	err := server.GetListener()
	if err != nil {
		t.Fail()
	}
	//go server.AcceptAndHandle()
	time.Sleep(1 * time.Second)

	passenger := NewPassenger()
	passenger.SetNearStation(server.Ip, server.Port)

	var head2 PassengerHead
	head2.Flag = 1
	head2.Id = 1234
	head2.IputNumber = 10
	head2.ListenTime = 3600 * 2
	head2.IsBroken = 0
	head2.BodyLen = int64(len(testData))
	tolist := []byte("mytest")
	listenlist := []byte("mytest")
	copy(head2.StrPushList[:], tolist)
	copy(head2.StrListenList[:], listenlist)

	err = passenger.ConnectNearStation()
	if err != nil {
		t.Fatal(err)
	}
	passenger.SetPassengerHead(&head2)

	err = passenger.SendData(&testData)
	if err != nil {
		t.Fatal(err)
	}
	err = passenger.DisConnectStation()
	if err != nil {
		t.Fatal(err)
	}
	//测试passenger端的控制

	err = passenger.SendData(&testData)
	if err == nil {
		t.Fatal("have disconnect can not senddata success")
	}

	//测试subway端的控制
	passenger.SetPassengerHead(&head2)
	passenger.IsConnect = true
	time.Sleep(1 * time.Second)
	err = passenger.SendData(&testData)
	if err == nil {
		t.Fatal("have disconnect can not senddata success")
	}

	return
}

func TestAdaptSend(t *testing.T) {
	var testData []byte = []byte("wuzhiming")
	ret, err := AdaptSend(&testData)
	if err == nil {
		fmt.Println(*ret)
	}
}
*/
