package sysinfo

import (
	//_ "myfirstbeego/routers"
	//"github.com/astaxie/beego"
	//"fmt"
	//"log"
	"os/exec"
	"strconv"
	"strings"
)

type LoadAverage struct {
	OneMin  float64
	FiMin   float64
	Quarter float64
}

type CPUPer struct {
	UserPer     float64
	SysPer      float64
	SkipLevePer float64
	IdlePer     float64
	WaPer       float64
	HiPer       float64
	SiPer       float64
	StPer       float64
}

type Mem struct {
	Total int
	Free  int
	Used  int
	Cache int
}

type SwapMem struct {
	Total    int
	Free     int
	Used     int
	AvailMem int
}

type SysStat struct {
	GetTime            string
	RuningTime         string
	Users              int
	SysLoad            LoadAverage
	TaskNumber         int
	RunningTaskNumber  int
	SleepingTaskNumber int
	StopedTaskNumber   int
	ZombieTaskNumber   int
	CPU                CPUPer
	MemInfo            Mem
	SwapInfo           SwapMem
}

func GetSysInfo() (sysStat SysStat, err error) {

	cmd := exec.Command("top", "-b", "-n 1")

	topres, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	strTop := string(topres)
	strTops := strings.SplitAfter(strTop, "avail Mem")
	SysInfo := strTops[0]
	SysInfos := strings.Fields(SysInfo)
	sysStat.GetTime = SysInfos[2]
	sysStat.RuningTime = strings.TrimRight(SysInfos[4], ",")

	sysStat.Users, err = strconv.Atoi(SysInfos[5])
	if err != nil {
		return
	}
	sysStat.SysLoad.OneMin, err = strconv.ParseFloat(strings.TrimRight(SysInfos[9], ","), 3)
	if err != nil {
		return
	}

	sysStat.SysLoad.FiMin, err = strconv.ParseFloat(strings.TrimRight(SysInfos[10], ","), 3)
	if err != nil {
		return
	}

	sysStat.SysLoad.Quarter, err = strconv.ParseFloat(SysInfos[11], 3)
	if err != nil {
		return
	}

	sysStat.TaskNumber, err = strconv.Atoi(SysInfos[13])
	if err != nil {
		return
	}

	sysStat.RunningTaskNumber, err = strconv.Atoi(SysInfos[15])
	if err != nil {
		return
	}

	sysStat.SleepingTaskNumber, err = strconv.Atoi(SysInfos[17])
	if err != nil {
		return
	}

	sysStat.StopedTaskNumber, err = strconv.Atoi(SysInfos[19])
	if err != nil {
		return
	}

	sysStat.ZombieTaskNumber, err = strconv.Atoi(SysInfos[21])
	if err != nil {
		return
	}

	sysStat.CPU.UserPer, err = strconv.ParseFloat(SysInfos[24], 3)
	if err != nil {
		return
	}

	sysStat.CPU.SysPer, err = strconv.ParseFloat(SysInfos[26], 3)
	if err != nil {
		return
	}

	sysStat.CPU.SkipLevePer, err = strconv.ParseFloat(SysInfos[28], 3)
	if err != nil {
		return
	}

	sysStat.CPU.IdlePer, err = strconv.ParseFloat(SysInfos[30], 3)
	if err != nil {
		return
	}

	sysStat.CPU.WaPer, err = strconv.ParseFloat(SysInfos[32], 3)
	if err != nil {
		return
	}

	sysStat.CPU.SiPer, err = strconv.ParseFloat(SysInfos[34], 3)
	if err != nil {
		return
	}

	sysStat.CPU.StPer, err = strconv.ParseFloat(SysInfos[36], 3)
	if err != nil {
		return
	}

	sysStat.MemInfo.Total, err = strconv.Atoi(SysInfos[43])
	if err != nil {
		return
	}
	sysStat.MemInfo.Free, err = strconv.Atoi(SysInfos[45])
	if err != nil {
		return
	}
	sysStat.MemInfo.Used, err = strconv.Atoi(SysInfos[47])
	if err != nil {
		return
	}
	sysStat.MemInfo.Cache, err = strconv.Atoi(SysInfos[49])
	if err != nil {
		return
	}

	sysStat.SwapInfo.Total, err = strconv.Atoi(SysInfos[53])
	if err != nil {
		return
	}

	sysStat.SwapInfo.Free, err = strconv.Atoi(SysInfos[55])
	if err != nil {
		return
	}

	sysStat.SwapInfo.Used, err = strconv.Atoi(SysInfos[57])
	if err != nil {
		return
	}

	sysStat.SwapInfo.AvailMem, err = strconv.Atoi(SysInfos[59])

	return
	//fmt.Printf("sysStat is %+v \n", sysStat)
}

type ProcessInfo struct {
	Command     string
	Pid         int
	User        string
	Priority    string
	Nice        int
	Virt        int //SWAP +RES
	Res         int //CODE + data
	Shr         int
	Stat        string
	CPUPer      float64 //cpuper from last top to now
	MemPer      float64
	RunningTime string //Process RunningTime
}

type ProcessInfos struct {
	Process []ProcessInfo
}

func GetProcessInfoNorMal() (ret ProcessInfos, err error) {
	//删除个性化配置
	cmd := exec.Command("rm", "~/.toprc", "-f")

	topres, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	cmd = exec.Command("top", "-b", "-n 1")

	topres, err = cmd.CombinedOutput()
	if err != nil {
		return
	}

	strTop := string(topres)
	strTops := strings.SplitAfter(strTop, "COMMAND")

	i := 0
	ProcessInfos := strings.Fields(strTops[1])

	for i < len(ProcessInfos) {
		var info ProcessInfo
		info.Pid, err = strconv.Atoi(ProcessInfos[i])
		if err != nil {
			return
		}

		info.User = ProcessInfos[i+1]
		info.Priority = ProcessInfos[i+2]
		info.Nice, err = strconv.Atoi(ProcessInfos[i+3])
		if err != nil {
			return
		}

		info.Virt, err = strconv.Atoi(ProcessInfos[i+4])
		if err != nil {
			return
		}

		info.Res, err = strconv.Atoi(ProcessInfos[i+5])
		if err != nil {
			return
		}

		info.Shr, err = strconv.Atoi(ProcessInfos[i+6])
		if err != nil {
			return
		}

		info.Stat = ProcessInfos[i+7]

		info.CPUPer, err = strconv.ParseFloat(ProcessInfos[i+8], 3)
		if err != nil {
			return
		}

		info.MemPer, err = strconv.ParseFloat(ProcessInfos[i+9], 3)
		if err != nil {
			return
		}

		info.RunningTime = ProcessInfos[i+10]
		info.Command = ProcessInfos[i+11]

		ret.Process = append(ret.Process, info)
		i = i + 12

	}

	return
}
