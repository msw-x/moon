package hw

import (
	"fmt"
	"runtime"

	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/umath"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type Status struct {
	Uptime          uint64
	Temperature     int
	CpuUsedPercent  int
	RamTotal        uint64
	RamUsed         uint64
	RamUsedPercent  int
	DiskTotal       uint64
	DiskUsed        uint64
	DiskUsedPercent int
}

func (o Status) Cpu() string {
	s := fmt.Sprintf("%d %%", o.CpuUsedPercent)
	if o.Temperature > 0 {
		s = fmt.Sprintf("%s / %d ℃", s, o.Temperature)
	}
	return fmt.Sprintf("%s - %d core(s)", s, runtime.NumCPU())
}

func (o Status) Ram() string {
	return fmt.Sprintf("%s / %s [%d %%]", ufmt.ByteSize(o.RamUsed), ufmt.ByteSize(o.RamTotal), o.RamUsedPercent)
}

func (o Status) Disk() string {
	return fmt.Sprintf("%s / %s [%d %%]", ufmt.ByteSize(o.DiskUsed), ufmt.ByteSize(o.DiskTotal), o.DiskUsedPercent)
}

func GetStatus() Status {
	cpuPercent, _ := cpu.Percent(0, false)
	vmStat, _ := mem.VirtualMemory()
	diskStat, _ := disk.Usage("/")
	uptime, _ := host.Uptime()
	return Status{
		Uptime:          uptime,
		Temperature:     CpuTemp(),
		CpuUsedPercent:  int(cpuPercent[0]),
		RamTotal:        vmStat.Total,
		RamUsed:         vmStat.Used,
		RamUsedPercent:  umath.Percent(vmStat.Used, vmStat.Total),
		DiskTotal:       diskStat.Total,
		DiskUsed:        diskStat.Used,
		DiskUsedPercent: umath.Percent(diskStat.Used, diskStat.Total),
	}
}
