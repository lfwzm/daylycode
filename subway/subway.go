package subway

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"net"
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

type Listener struct {
	Nodeid int64
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
	worker   StationHandler //接收到请求后的处理器 修改为接口
}

func NewSubWayStation() *SubWayStation {
	ret := SubWayStation{}
	ret.worker = nil
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
		for {
			mutex.Lock()
			//从处理后数据map中获取数据
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
					id := head.Id

					_, ok := DataListSingleton.Listeners[key]
					if !ok {
						//返回错误
						continue
					}
					var i int64
					//从监听列表中获取监听客户端
					for _, v := range DataListSingleton.Listeners[key] {
						i = i + 1
						listener := v
						if listener.Nodeid == id {
							go func() {
								fmt.Println(i)
								(*listener.Conn).Write(retvalue)

							}()
							value.Remove(data)
						}
					}
				}
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
		fmt.Println(conn.LocalAddr())
		go n.worker.HandleConn(&conn) //高并发情况下，会有些问题。修改问使用goroutine
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
	IputNumber    int64    //参数个数
	BodyLen       int64    //消息体长度
	ListenTime    int64    //监听时长
	IsBroken      int      //主动断开
	IsFirst       int      // 是否第一次调用，1 是 0 不是 是第一次调用要注册对应的监听队列。
}

type Mydata struct {
	Name  string `json:Name`
	Phone string `json:Phone`
}

func (n *NormalHandler) HandleConn(Conn *net.Conn) error {
	fmt.Println("connect ok")
	for {
		var readData []byte = make([]byte, 4096) //每次都new是否影响性能。

		Len, err := (*Conn).Read(readData)
		if err != nil {
			//
			fmt.Println(err)
			return err
		}
		fmt.Println("read ok")
		var head *PassengerHead = *(**PassengerHead)(unsafe.Pointer(&readData))
		headlen := unsafe.Sizeof(*head)
		listenTime := head.ListenTime
		err = (*Conn).SetDeadline(time.Now().Add(time.Duration(listenTime) * time.Second))
		if err != nil {
			//
			fmt.Println("SetDeadline eror")
			return err
		}
		if head.IsBroken == 1 {
			fmt.Println("connect is Broken")
			(*Conn).Close()
			break
		}
		fmt.Println("**********data len is: ", Len)
		if Len < (int)(headlen) {
			fmt.Println("to short, Len is ", Len, " headlen is :", (int)(headlen))
			return nil
		}

		fmt.Println("InputNumber is ", head.IputNumber, " id is: ", head.Id)
		tolist := string(head.StrPushList[:])

		mutex.Lock()
		_, ok := DataListSingleton.Datas[tolist]
		if !ok {
			DataListSingleton.Datas[tolist] = list.New()
			//listdata, _ := DataListSingleton.Datas[tolist]
		}

		pushret := DataListSingleton.Datas[tolist].PushFront(readData)
		if pushret != nil {
			fmt.Println(pushret)
		}

		fmt.Println("to list is: ", tolist)
		bodyLen := head.BodyLen

		body := readData[int64(headlen) : int64(headlen)+bodyLen]
		fmt.Println("body is :", string(body))
		var mybody Mydata
		fmt.Println("body len is: ", head.BodyLen)
		/*mybodydata, */ err = json.Unmarshal(body, &mybody)

		if err != nil {
			fmt.Println(err)
			return nil
		}

		fmt.Println("body is :", mybody)
		fmt.Println("body is :", string(body))

		if head.IsFirst == 1 {
			listenlist := string(head.StrListenList[:])
			fmt.Println("listen list is :", listenlist)

			_, ok := DataListSingleton.Listeners[listenlist]
			//将节点加入到监听队列
			if !ok {
				listener := Listener{Nodeid: head.Id, Conn: Conn}
				DataListSingleton.Listeners[listenlist] = make([]Listener, 1)
				DataListSingleton.Listeners[listenlist] = append(DataListSingleton.Listeners[listenlist], listener)
				fmt.Println("appand listener ok , listener is: ", listener, " len of DataListSingleton.Listeners[mytest] is ", len(DataListSingleton.Listeners[listenlist]))

			} else {
				for _, v := range DataListSingleton.Listeners[listenlist] {
					Listener := v
					if Listener.Nodeid == head.Id {
						/*
							//把原有的连接关闭。
							if err := (*Listener.Conn).Close(); err != nil {
								fmt.Println(err) //修改为记录日志
							}
						*/
						Listener.Conn = Conn
					}
				}
			}
		}
		/*
			if head.Flag == 1 {
				sendto := string(head.StrListenList)
				_, ok := DataListSingleton.Listeners[listenlist]
				if !ok {
					fmt.Println("listener error")
					return
				}
			}
		*/
		mutex.Unlock()
	}
	return nil
}
