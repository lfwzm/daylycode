package subway

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
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

	if n.listener == nil {
		return errors.New(n.Name + ".listener is nil")
	}

	go func() {
		for {
			mutex.Lock()
			for key, datas := range DataListSingleton.Datas {
				//fmt.Println("++++++++++++++key is: ", key)
				value := datas
				for {
					if value.Len() == 0 {
						break
					}

					data := value.Front()
					retvalue := data.Value.([]byte)
					var head *XpHead = *(**XpHead)(unsafe.Pointer(&retvalue))

					//listener := string(head.StrListenList[:])
					id := head.Id
					_, ok := DataListSingleton.Listeners[key]
					if !ok {
						//返回错误
						continue
					}
					var i int64
					for _, v := range DataListSingleton.Listeners[key] {
						i = i + 1
						listener := v
						if listener.Nodeid == id {
							go func() {
								(*listener.Conn).Write(retvalue)
								datas.Remove(data)
							}()

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

		go n.worker.HandleConn(&conn) //高并发情况下，会有些问题。修改问使用goroutine
		//处理每个连接请求

	}
	return nil
}

/*
// LBM调用入参信息结构
struct ST_LBM_PARA_INFO
{
  char   szENName[32];          // 参数英文名（参数代码）
  char   szCNName[32];          // 参数中文名
  bool   bIsNeed;               // 是否必送
  short  siOffset;              // 在参数内容中偏移量
  short  siLen;                 // 在参数内容中长度
};
*/

type ParaInfo struct {
	StrEnName [32]byte
	StrChName [32]byte
	BisNeed   bool
	Ioffset   int32
	Ilen      int32
}

type XpHead struct {
	IsGiz         int64    //是否需要压缩
	StrRuntime    [32]byte //执行时间
	StrListenList [32]byte //从某个队列获取返回结果
	StrPushList   [32]byte //入参发送到对应的队列
	Flag          int64    //0 入参， 1出参
	IsFirst       int      // 是否第一次调用，1 是 0 不是 是第一次调用要注册对应的监听队列。
	Id            int64    //发送或者接收者的Id号
	IputNumber    int64    //参数个数
	BodyLen       int64    //消息体长度
}

type Mydata struct {
	Name  string `json:Name`
	Phone string `json:Phone`
}

func (n *NormalHandler) HandleConn(Conn *net.Conn) error {

	var readData []byte = make([]byte, 4096) //每次都new是否影响性能。
	Len, err := (*Conn).Read(readData)
	if err != nil {
		//
		fmt.Println("read eror")
		return err
	}

	var head *XpHead = *(**XpHead)(unsafe.Pointer(&readData))
	headlen := unsafe.Sizeof(*head)
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
					if err := (*Listener.Conn).Close(); err != nil {
						fmt.Println(err) //修改为记录日志
					}
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
	return nil
}
