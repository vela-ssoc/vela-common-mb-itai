package sigar

import (
	"bufio"
	"context"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func CPU() (*CPUStat, error) {
	f := cpuStat("/proc/stat")
	return f.Read()
}

func CPUPercent(ctx context.Context, interval time.Duration) (float64, error) {
	if interval <= 0 {
		return 0, nil
	}

	c1, err := CPU()
	if err != nil {
		return 0, err
	}

	if err = Park(ctx, interval); err != nil {
		return 0, err
	}

	c2, err := CPU()
	if err != nil {
		return 0, err
	}
	per := c1.Percent(c2)

	return per, nil
}

type cpuStat string

func (cs cpuStat) Read() (*CPUStat, error) {
	file, err := os.Open(string(cs))
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	ret := new(CPUStat)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue
		}
		if fields[0] != "cpu" {
			continue
		}
		for i := 1; i < 11; i++ {
			val, _ := strconv.ParseUint(fields[i], 10, 64)
			ret.setValue(i, val)
		}
		break
	}

	return ret, nil
}

type CPUStat struct {
	User      uint64 // 从系统启动开始累积到当前时刻，处于用户态的运行时间，不包含 nice 值为负的进程。
	Nice      uint64 // 从系统启动开始累积到当前时刻，nice 值为负的进程所占用的 CPU 时间。
	System    uint64 // 从系统启动开始累积到当前时刻，处于核心态的运行时间。
	Idle      uint64 // 从系统启动开始累积到当前时刻，除 IO 等待时间以外的其他等待时间。
	IOWait    uint64 // 从系统启动开始累积到当前时刻，IO 等待时间。(since 2.5.41)
	IRQ       uint64 // 从系统启动开始累积到当前时刻，硬中断时间。(since 2.6.0-test4)
	SoftIRQ   uint64 // 从系统启动开始累积到当前时刻，软中断时间。(since 2.6.0-test4)
	Steal     uint64
	Guest     uint64
	GuestNice uint64
}

func (s CPUStat) Total() uint64 {
	return s.User + s.Nice + s.System + s.Idle + s.IOWait +
		s.IRQ + s.SoftIRQ + s.Steal + s.Guest + s.GuestNice
}

func (s CPUStat) Used() uint64 {
	total := s.Total()
	return total - s.Idle - s.IOWait
}

func (s CPUStat) Percent(after *CPUStat) float64 {
	t1, u1 := s.Total(), s.Used()
	t2, u2 := after.Total(), after.Used()
	if u2 <= u1 {
		return 0
	}
	if t2 <= t1 {
		return 100
	}

	return math.Min(100, math.Max(0, float64(u2-u1)/float64(t2-t1)*100))
}

func (s *CPUStat) setValue(idx int, val uint64) {
	switch idx {
	case 1:
		s.User = val
	case 2:
		s.Nice = val
	case 3:
		s.System = val
	case 4:
		s.Idle = val
	case 5:
		s.IOWait = val
	case 6:
		s.IRQ = val
	case 7:
		s.SoftIRQ = val
	case 8:
		s.Steal = val
	case 9:
		s.Guest = val
	case 10:
		s.GuestNice = val
	}
}
