package subway

import (
	"errors"
	"fmt"
	"net"
	"time"
	"unsafe"
)

var (
	PoitNilError       = errors.New("*Passenger is nil")
	defaultDailTimeOut = 60 * time.Second
	head               PassengerHead
	HeadSize           = unsafe.Sizeof(head)
	SendSize           = 4096 - HeadSize
)

type Passenger struct {
	Stationinfo StationInfo
	Conn        *net.Conn     //网络连接
	IsConnect   bool          //是否已经连接
	Head        PassengerHead //头部信息
}

type StationInfo struct {
	NearStationIp   string //最近站点的IP
	NearStationPort string //最近站点的端口
	IsConnect       bool   //是否已经连接
}

//新建
func NewPassenger() *Passenger {
	return &Passenger{}
}

//设置最近站点信息
func (pP *Passenger) SetNearStation(ip, port string) (err error) {
	if pP == nil {
		return PoitNilError
	}

	pP.Stationinfo.NearStationIp, pP.Stationinfo.NearStationPort = ip, port
	return nil
}

//查看最近站点信息
func (pP *Passenger) GetNearStation() (ret StationInfo, err error) {
	if pP == nil {
		return ret, PoitNilError
	}

	return pP.Stationinfo, nil
}

//链接最近站点
func (pP *Passenger) ConnectNearStation() error {
	if pP == nil {
		return PoitNilError
	}

	TimeOut := defaultDailTimeOut
	conn, err := net.DialTimeout("tcp", pP.Stationinfo.NearStationIp+":"+pP.Stationinfo.NearStationPort, TimeOut)
	if err == nil {
		pP.Conn = &conn
		pP.IsConnect = true
	}
	return err
}

//设置头部信息
func (pP *Passenger) SetPassengerHead(head *PassengerHead) error {
	if pP == nil || head == nil {
		return PoitNilError
	}
	pP.Head = *head
	return nil
}

//适配发送数据，以减少发送期间的数据拷贝。其实也减少不了多少^()^
func AdaptSend(pData *[]byte) (pRet *[]byte, err error) {
	if pData == nil {
		return nil, errors.New("AdaptSend: pData is nil")
	}

	iLen := len(*pData)
	ret := make([]byte, iLen+int(HeadSize))
	icpyLen := copy(ret[HeadSize-1:], *pData)
	if icpyLen != iLen {
		return nil, errors.New("AdaptSend: copy failed")
	}
	return &ret, nil
}

//发送数据, 基本发送功能ok 后续需要做优化，尽量减少拷贝
func (pP *Passenger) SendData(pData *[]byte) error {
	if pP == nil {
		return PoitNilError
	}
	//fmt.Println("ok")
	if pData == nil {
		return errors.New("SendData: pData is nil")
	}

	if pP.IsConnect == false {
		return errors.New("SendData: no Connect to station")
	}

	dataLen := len(*pData)
	sendData := make([]byte, 4096)
	i := 0
	for i < dataLen {
		//增加处理

		sendLen := copy(sendData[int(HeadSize):], (*pData)[int(i):])
		if sendLen < int(SendSize) {
			pP.Head.IsLast = 1
		}
		pP.Head.IsLast = int64(sendLen)

		var head2 *PassengerHead = *(**PassengerHead)(unsafe.Pointer(&sendData))
		*head2 = pP.Head
		//发送数据
		if pP.IsConnect == false {
			return errors.New("SendData: no Connect to station")
		}
		writeLen, err := (*pP.Conn).Write(sendData)
		if err != nil || writeLen != 4096 {
			fmt.Println(err)
			fmt.Println("writeLen is :", writeLen)
			//增加写日志
			pP.IsConnect = false
			return err
		}
		//fmt.Println("-----------in passenger writeLen is :", writeLen)
		//fmt.Println("-----------in passenger sendData is :", string(sendData[HeadSize:]))
		i = i + int(writeLen) - int(HeadSize)
		//大部分情况下，客户端发给服务端的参数都是小于4094字节大小的。
		if i < dataLen {
			for index, _ := range sendData {
				sendData[index] = 0
			}
		}
	}
	return nil
}

//接收数据

//主动断开连接

//发送心跳数据
