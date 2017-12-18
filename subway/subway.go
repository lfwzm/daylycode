package subway

import (
	"container/list"
	//"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

var mutex sync.Mutex
var DataListSingleton = NewDataList()

type DataList struct {
	ListNumber int64
	Datas      map[string]*list.List
	Listeners  map[string][]Listener
}

//存放接收者信息
type Listener struct {
	Nodeid int64 //接收者节点号
	Conn   *net.Conn
}

func NewDataList() *DataList {
	ret := &DataList{}
	ret.Datas = make(map[string]*list.List, 256)
	ret.Listeners = make(map[string][]Listener, 32)
	return ret
}

type HandleFunc func(*net.Conn) error

type StationHandler interface {
	HandleConn(Conn *net.Conn) error
}

type NormalHandler struct {
}

//服务节点
type SubWayStation struct {
	Name     string         //站点名称
	Ip       string         //站点所在的ip
	Port     string         //站点对应的端口号
	listener net.Listener   //服务监听器
	Worker   StationHandler //接收到请求后的处理器 修改为接口
}

func NewSubWayStation() *SubWayStation {
	ret := SubWayStation{}
	ret.Worker = nil
	return &ret
}

func (n *SubWayStation) GetListener() (err error) {
	n.listener, err = net.Listen("tcp", n.Ip+":"+n.Port)

	return
}

func (n *SubWayStation) AcceptAndHandle() error {

	if n == nil || n.listener == nil {
		return errors.New(n.Name + ".listener is nil")
	}

	go func() {
		//i := 0
		for {
			mutex.Lock()
			//从处理后数据map中获取数据

			if len(DataListSingleton.Datas) == 0 {
				time.Sleep(1 * time.Millisecond)
				mutex.Unlock()
				continue
			}
			/*
				if i%10 == 0 {
					log.Println("DataListSingleton.Datas len is :", len(DataListSingleton.Datas))
					log.Println("DataListSingleton.Listeners len is :", len(DataListSingleton.Listeners))
				}
				i++
			*/
			fmt.Println("DataListSingleton.Listeners len is :", len(DataListSingleton.Listeners))
			fmt.Println("DataListSingleton.Datas len is :", len(DataListSingleton.Datas))
			for key, datas := range DataListSingleton.Datas {
				value := datas
				for {
					//列表中无数据退出循环
					if value.Len() == 0 {
						break
					}
					//获取第一个数据
					data := value.Front()
					retvalue := data.Value.([]byte)
					var head *PassengerHead = *(**PassengerHead)(unsafe.Pointer(&retvalue))

					//listener := string(head.StrListenList[:])
					fmt.Println("key is", key)

					id := head.Id
					fmt.Println("head id is:", id)

					_, ok := DataListSingleton.Listeners[key]
					if !ok {
						//返回错误
						fmt.Println("continue")
						continue
					}

					//从监听列表中获取监听客户端
					for _, v := range DataListSingleton.Listeners[key] {

						listener := v
						if listener.Nodeid == id {
							fmt.Print("nodeid is: ", listener.Nodeid)
							go func() {
								fmt.Println("in write Listener is : ", listener)
								_, err := (*listener.Conn).Write(retvalue)
								if err != nil {
									//fmt.Println("subway wLen is:", wLen)
									fmt.Println(err)
								}

							}()
							value.Remove(data)
						}
					}
				}
				delete(DataListSingleton.Datas, key)
			}
			mutex.Unlock()
		}
	}()

	for {
		conn, err := (n.listener).Accept()
		if err != nil {
			//异常处理

			return err
		}
		//fmt.Println(conn.LocalAddr())
		go n.Worker.HandleConn(&conn) //每个连接开一个协程处理
		//处理每个连接请求

	}
	return nil
}

//入参信息结构定义
type ParaInfo struct {
	StrEnName [32]byte // 参数英文名（参数代码）
	StrChName [32]byte // 参数中文名
	BisNeed   bool     // 是否必送
	Ioffset   int32    // 在参数内容中偏移量
	Ilen      int32    // 在参数内容中长度
}

//乘客头信息
type PassengerHead struct {
	IsGiz         int64    //是否需要压缩
	StrRuntime    [32]byte //执行时间 （格式化，或者测试下timestamp是否为定长）
	StrListenList [32]byte //从某个队列获取返回结果
	StrPushList   [32]byte //入参发送到对应的队列
	Flag          int64    //0 入参， 1出参
	Id            int64    //发送或者接收者的Id号
	IputNumber    int8     //参数个数
	BodyLen       int64    //消息体长度
	IsLast        int64    //是否为同一个调用的最后一次
	ListenTime    int64    //监听时长
	IsBroken      int      //主动断开
}

type Mydata struct {
	Name  string `json:Name`
	Phone string `json:Phone`
}

var doTime int64 = 0

//需要解决何时关闭socket的问题。
func (n *NormalHandler) HandleConn(Conn *net.Conn) error {
	//fmt.Println("connect ok")
	//var readData []byte = make([]byte, 4096) //每次都new是否影响性能
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Println(string(buf))
		}
	}()

	pwd, _ := os.Getwd()
	file_path := pwd + "/src.log"
	newfile, err := os.Create(file_path)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(newfile)
	err = (*Conn).SetDeadline(time.Now().Add(time.Duration(3600) * time.Second))
	if err != nil {
		//
		fmt.Println("SetDeadline eror")
		return err
	}
	var ReadData = make([]byte, 4096)
	var ReadfixData = make([]byte, 4096)
	index := 0
	for {
		var readDataTmp []byte = make([]byte, 4096) //每次都new是否影响性能。

		Len, err := (*Conn).Read(readDataTmp)
		if err != nil {
			//
			//fmt.Println("in HandleConn")
			fmt.Println(err)
			return err
		}

		//fmt.Println("read Len is:", Len)
		for i := 0; i < Len; i++ {

			//fmt.Println("index is: ", index, " i is : ", i)
			ReadData[index] = readDataTmp[i]
			index++
			if index == 4096 {
				fixLen := copy(ReadfixData, ReadData)
				if fixLen != 4096 {
					//报错
					fmt.Println("fixLen is:", fixLen)
				}
				index = 0
				Input := ReadfixData
				log.Println("Input is:", Input)
				go EquelHandle(Input, Conn, fixLen)
			}

		}

		//当接收到的数据小于长度小于4096的时候
		//copylen := copy(ReadData[index:], readDataTmp)

		//在连接数量比较小的时候确实可以达到提高性能的作用。但是如果链接数量达到一定数量级的时候。需要考虑是否每一次read都开一个协程处理是否合理
		/*
			Input := readDataTmp
			log.Println("Inputnumber is :", Len)
			go EquelHandle(Input, Conn, Len)
		*/
		/*
			go func() {

				var head *PassengerHead = *(**PassengerHead)(unsafe.Pointer(&readData))
				headlen := unsafe.Sizeof(*head)
				if head.IsBroken == 1 {
					(*Conn).Close()
					return
				}

				//fmt.Println("Inputnumber is :", head.IputNumber)
				if Len < (int)(headlen) {
					//增加错误日志
					return
				}

				tolist := string(head.StrPushList[:])
				mutex.Lock()

				//把数据插入队列请求
				_, ok := DataListSingleton.Datas[tolist]
				if !ok {
					DataListSingleton.Datas[tolist] = list.New()
				}

				pushret := DataListSingleton.Datas[tolist].PushFront(readData)
				if pushret == nil {
					//fmt.Println(pushret)
					//增加出错处理
					return
				}

				listenlist := string(head.StrListenList[:])
				_, ok = DataListSingleton.Listeners[listenlist]
				//接收节点不存在则新增
				if !ok {

					listener := Listener{Nodeid: head.Id, Conn: Conn}
					DataListSingleton.Listeners[listenlist] = make([]Listener, 1)
					DataListSingleton.Listeners[listenlist] = append(DataListSingleton.Listeners[listenlist], listener)
					//fmt.Println("appand listener ok , listener is: ", listener, " len of DataListSingleton.Listeners[mytest] is ", len(DataListSingleton.Listeners[listenlist]))
				} else {

					fundNodeid := false
					for _, v := range DataListSingleton.Listeners[listenlist] {
						Listener := v
						if Listener.Nodeid == head.Id {
							Listener.Conn = Conn
							fundNodeid = true
						}
					}

					//节点并不存在，增加节点
					if !fundNodeid {
						listener := Listener{Nodeid: head.Id, Conn: Conn}
						DataListSingleton.Listeners[listenlist] = append(DataListSingleton.Listeners[listenlist], listener)
					}

				}
				mutex.Unlock()
			}()
		*/
	}
	return nil
}

func EquelHandle(readData []byte, Conn *net.Conn, Len int) {
	//fmt.Println("in EquelHandle")

	var head *PassengerHead = *(**PassengerHead)(unsafe.Pointer(&readData))
	headlen := unsafe.Sizeof(*head)
	/*
		if head.IsBroken == 1 {
			(*Conn).Close()
			return
		}
	*/
	//log.Println("Inputnumber is :", readData)
	if Len < (int)(headlen) {
		//增加错误日志
		return
	}
	mutex.Lock()
	if head.StrPushList[0] != 0 {
		tolist := string(head.StrPushList[:])
		//fmt.Println("tolist is:", tolist)
		//把数据插入队列请求
		_, ok := DataListSingleton.Datas[tolist]
		if !ok {
			DataListSingleton.Datas[tolist] = list.New()
			fmt.Println("create a new list")
		}

		pushret := DataListSingleton.Datas[tolist].PushFront(readData)
		if pushret == nil {
			//fmt.Println(pushret)
			//增加出错处理
			log.Println("DataListSingleton.Datas[", tolist, "].PushFront(", readData, ") failed")
			return
		}
	}

	listenlist := string(head.StrListenList[:])
	//fmt.Println("listenlist is:", listenlist)
	_, ok := DataListSingleton.Listeners[listenlist]
	//接收节点不存在则新增
	if !ok {

		listener := Listener{Nodeid: head.Id, Conn: Conn}
		fmt.Println("Listener is:", listener)
		DataListSingleton.Listeners[listenlist] = make([]Listener, 1)
		DataListSingleton.Listeners[listenlist] = append(DataListSingleton.Listeners[listenlist], listener)
		//fmt.Println("appand listener ok , listener is: ", listener, " len of DataListSingleton.Listeners[mytest] is ", len(DataListSingleton.Listeners[listenlist]))
	} else {

		fundNodeid := false
		for i, Listener := range DataListSingleton.Listeners[listenlist] {
			//Listener := v
			fmt.Println("Listener is:", Listener)
			if Listener.Nodeid == head.Id {
				//Listener.Conn = Conn
				DataListSingleton.Listeners[listenlist][i].Conn = Conn
				fundNodeid = true
			}
		}

		//节点并不存在，增加节点
		if !fundNodeid {
			listener := Listener{Nodeid: head.Id, Conn: Conn}
			DataListSingleton.Listeners[listenlist] = append(DataListSingleton.Listeners[listenlist], listener)
			fmt.Println("add new conn", listener)
		}

	}
	mutex.Unlock()
}
