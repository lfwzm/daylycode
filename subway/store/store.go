package store

//消息索引
type MsgIndex struct {
	MsgNumber     int64
	MsgFileNumber int64
	Index         int64
	Len           int64
}

//消息存储接口
type MsgStore interface {
	Init() error                                                     //初始化
	StoreMsg(data []byte, topic string) (index *MsgIndex, err error) //存储消息
	GetMaxMsgIndex() (index *MsgIndex)                               //获取最大的消息号
	SetMinMsgFileNum(num int64) error                                //设置最小ID号
	Close() error
}

//后续加上备份的接口
