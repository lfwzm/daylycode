package store

import (
	"errors"
	"os"
	"strconv"
	"subway/store"
)

//以文件的形式存储消息
type FileStore struct {
	Topic      string //消息主题
	baseDir    string //初始目录
	maxIndex   int64  //目前的最大消息号
	maxFileNum int64  //最大的消息文件号
	minFileNum int64  //最小的消息文件号
	msgnumber  int64  //消息最大号
	curentfd   *os.File
}

func NewFileStore(topic string, base string) *FileStore {
	if len(topic) == 0 {
		return nil
	}

	if f, err := os.Stat(base); err != nil || f.IsDir() != true {
		//增加日志
		return nil
	}

	return &FileStore{Topic: topic, baseDir: base, maxFileNum: 1}
}

var NullPtr = errors.New("Ptr is nil")

func (pF *FileStore) Init() error {
	if pF == nil {
		return NullPtr
	}
	pF.maxIndex = 0
	pF.maxFileNum = 1
	for pF.maxIndex >= 1024*1024*1024 || pF.curentfd == nil {
		if pF.curentfd != nil {

			pF.curentfd.Close()
			pF.curentfd = nil
			pF.maxFileNum = pF.maxFileNum + 1
		}

		curentfd, err := os.OpenFile(pF.baseDir+"/"+pF.Topic+strconv.Itoa(int(pF.maxFileNum)), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		fileLen, err := curentfd.Seek(0, os.SEEK_END)
		if err != nil {
			return err
		}

		pF.maxIndex = fileLen
		pF.curentfd = curentfd
	}
	//需要增加对msgnumber的初始化
	return nil
}

func (pF *FileStore) StoreMsg(data []byte, topic string) (index *store.MsgIndex, err error) {
	if pF == nil {
		return nil, NullPtr
	}
	if len(data) == 0 {
		return nil, errors.New("(pF *FileStore)StoreMsg : data len is 0")
	}
	if topic != pF.Topic {
		return nil, errors.New("pF *FileStore)StoreMsg : topic error")
	}

	for pF.maxIndex >= 1024*1024*1024 || pF.curentfd == nil { //有问题
		if pF.curentfd != nil {

			pF.curentfd.Close()
			pF.curentfd = nil
		}
		pF.maxFileNum = pF.maxFileNum + 1
		curentfd, err := os.OpenFile(pF.baseDir+"/"+pF.Topic+strconv.Itoa(int(pF.maxFileNum)), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		fileLen, err := curentfd.Seek(0, os.SEEK_END)
		if err != nil {
			return nil, err
		}
		pF.maxIndex = fileLen
		pF.curentfd = curentfd
	}
	writeLen, err := pF.curentfd.Write(data)
	if err != nil {
		return nil, err
	}
	index = &store.MsgIndex{}
	index.MsgFileNumber = pF.maxFileNum
	index.Index = pF.maxIndex
	index.Len = int64(writeLen)
	pF.maxIndex = pF.maxIndex + int64(writeLen)
	pF.msgnumber = pF.msgnumber + 1
	index.MsgNumber = pF.msgnumber
	return
}

func (pF *FileStore) GetMaxMsgnumber() (index *store.MsgIndex) {
	if pF == nil {
		return nil
	}
	index = &store.MsgIndex{}
	index.MsgNumber = pF.msgnumber
	index.Len = 0
	index.MsgFileNumber = pF.maxFileNum
	index.Index = pF.maxIndex
	return index

}

func (pF *FileStore) GetMinMsgnumber() (index *store.MsgIndex) {
	if pF == nil {
		return nil
	}
	/*
		index = &store.MsgIndex{}
		index.MsgNumber = pF.msgnumber
		index.Len = 0
		index.MsgFileNumber = pF.maxFileNum
		index.Index = pF.maxIndex
	*/
	return index

}

func (pF *FileStore) SetMinMsgFileNum(num int64) error {
	if pF == nil {
		return NullPtr
	}
	if num < pF.minFileNum {
		return errors.New("SetMinMsgFileNum : number set invalue")
	}
	if num > pF.maxFileNum {
		return errors.New("SetMinMsgFileNum : number too big")
	}
	if num == pF.maxFileNum {
		pF = nil
		return nil
	}
	pF.minFileNum = num
	return nil
}

func (pF *FileStore) Close() error {
	if pF == nil {
		return NullPtr
	}
	if pF.curentfd != nil {
		pF.curentfd.Close()
		pF.curentfd = nil
	}
	return nil
}
