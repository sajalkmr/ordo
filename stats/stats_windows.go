//go:build windows
// +build windows

package stats

import (
	"syscall"
	"unsafe"
)

type Stats struct {
	MemStats  *MemInfo
	DiskStats *Disk
	CpuStats  *CPUStat
	LoadStats *LoadAvg
	TaskCount int
}

type MemInfo struct {
	MemTotal     uint64
	MemAvailable uint64
}

type Disk struct {
	All  uint64
	Free uint64
	Used uint64
}

type CPUStat struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	IOWait  uint64
	IRQ     uint64
	SoftIRQ uint64
	Steal   uint64
}

type LoadAvg struct {
	Last1Min  float64
	Last5Min  float64
	Last15Min float64
}

func (s *Stats) MemUsedKb() uint64 {
	return s.MemStats.MemTotal - s.MemStats.MemAvailable
}

func (s *Stats) MemUsedPercent() uint64 {
	if s.MemStats.MemTotal == 0 {
		return 0
	}
	return (s.MemStats.MemTotal - s.MemStats.MemAvailable) * 100 / s.MemStats.MemTotal
}

func (s *Stats) MemAvailableKb() uint64 {
	return s.MemStats.MemAvailable
}

func (s *Stats) MemTotalKb() uint64 {
	return s.MemStats.MemTotal
}

func (s *Stats) DiskTotal() uint64 {
	return s.DiskStats.All
}

func (s *Stats) DiskFree() uint64 {
	return s.DiskStats.Free
}

func (s *Stats) DiskUsed() uint64 {
	return s.DiskStats.Used
}

func (s *Stats) CpuUsage() float64 {
	idle := s.CpuStats.Idle + s.CpuStats.IOWait
	nonIdle := s.CpuStats.User + s.CpuStats.Nice + s.CpuStats.System + s.CpuStats.IRQ + s.CpuStats.SoftIRQ + s.CpuStats.Steal
	total := idle + nonIdle

	if total == 0 && idle == 0 {
		return 0.00
	}

	return (float64(total) - float64(idle)) / float64(total)
}

func GetStats() *Stats {
	return &Stats{
		MemStats:  GetMemoryInfo(),
		DiskStats: GetDiskInfo(),
		CpuStats:  GetCpuStats(),
		LoadStats: GetLoadAvg(),
	}
}

func GetMemoryInfo() *MemInfo {
	var memStatus syscall.MemoryStatusEx
	memStatus.Length = uint32(unsafe.Sizeof(memStatus))
	syscall.GlobalMemoryStatusEx(&memStatus)

	return &MemInfo{
		MemTotal:     uint64(memStatus.TotalPhys) / 1024, // Convert to KB
		MemAvailable: uint64(memStatus.AvailPhys) / 1024, // Convert to KB
	}
}

func GetDiskInfo() *Disk {
	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64
	path := "C:\\"
	syscall.GetDiskFreeSpaceEx(
		syscall.StringToUTF16Ptr(path),
		&freeBytesAvailable,
		&totalNumberOfBytes,
		&totalNumberOfFreeBytes,
	)

	return &Disk{
		All:  totalNumberOfBytes / 1024, // Convert to KB
		Free: totalNumberOfFreeBytes / 1024,
		Used: (totalNumberOfBytes - totalNumberOfFreeBytes) / 1024,
	}
}

func GetCpuStats() *CPUStat {
	// Windows doesn't provide detailed CPU stats like Linux
	// Return a default CPUStat with zeros
	return &CPUStat{}
}

func GetLoadAvg() *LoadAvg {
	// Windows doesn't provide load average like Linux
	// Return a default LoadAvg with zeros
	return &LoadAvg{}
}
