package meminfo

import (
	"os/exec"
	"strconv"
	"strings"
)

type Meminfo struct {
	TotalMem     int
	UsedMem      int
	FreeMem      int
	ShareMem     int
	CacheMem     int
	AvailableMem int
}

func GetMemInfo() (ret *Meminfo, err error) {
	ret = &Meminfo{}
	cmd := exec.Command("free")

	freeret, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	strFree := string(freeret)
	strFrees := strings.SplitAfter(strFree, ":")
	meminfos := strings.Fields(strFrees[1])
	ret.TotalMem, err = strconv.Atoi(meminfos[0])
	if err != nil {
		return nil, err
	}

	ret.UsedMem, err = strconv.Atoi(meminfos[1])
	if err != nil {
		return nil, err
	}

	ret.FreeMem, err = strconv.Atoi(meminfos[2])
	if err != nil {
		return nil, err
	}

	ret.ShareMem, err = strconv.Atoi(meminfos[3])
	if err != nil {
		return nil, err
	}

	ret.CacheMem, err = strconv.Atoi(meminfos[4])
	if err != nil {
		return nil, err
	}

	ret.AvailableMem, err = strconv.Atoi(meminfos[5])
	if err != nil {
		return nil, err
	}

	return ret, nil
}
