package sigar

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

// Memory https://www.kernel.org/doc/Documentation/filesystems/proc.txt
func Memory() (*Meminfo, error) {
	f := linuxMeminfo("/proc/meminfo")
	return f.Read()
}

type linuxMeminfo string

func (lm linuxMeminfo) Read() (*Meminfo, error) {
	file, err := os.Open(string(lm))
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	ret := new(Meminfo)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		plain := sc.Text()
		key, val := parseLine(plain)
		ret.setValue(key, val)
	}

	return ret, nil
}

func (lm linuxMeminfo) parseLine(line string) (string, uint64) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", 0
	}

	val, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return "", 0
	}
	key := strings.TrimRight(fields[0], ":")
	var unit string
	if len(fields) > 2 {
		unit = fields[2]
	}
	switch unit {
	case "kB":
		val *= 1024
	}

	return key, val
}

func parseLine(plain string) (string, uint64) {
	fields := strings.Fields(plain)
	if len(fields) < 2 {
		return "", 0
	}

	val, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return "", 0
	}
	key := strings.TrimRight(fields[0], ":")
	var unit string
	if len(fields) > 2 {
		unit = fields[2]
	}
	switch unit {
	case "kB":
		val *= 1024
	}

	return key, val
}

type Meminfo struct {
	MemTotal          uint64 `json:"mem_total"`
	MemFree           uint64 `json:"mem_free"`
	MemAvailable      uint64 `json:"mem_available"`
	Buffers           uint64 `json:"buffers"`
	Cached            uint64 `json:"cached"`
	SwapCached        uint64 `json:"swap_cached"`
	Active            uint64 `json:"active"`
	Inactive          uint64 `json:"inactive"`
	ActiveAnon        uint64 `json:"active_anon"`
	InactiveAnon      uint64 `json:"inactive_anon"`
	ActiveFile        uint64 `json:"active_file"`
	InactiveFile      uint64 `json:"inactive_file"`
	Unevictable       uint64 `json:"unevictable"`
	Mlocked           uint64 `json:"mlocked"`
	SwapTotal         uint64 `json:"swap_total"`
	SwapFree          uint64 `json:"swap_free"`
	Zswap             uint64 `json:"zswap"`
	Zswapped          uint64 `json:"zswapped"`
	Dirty             uint64 `json:"dirty"`
	Writeback         uint64 `json:"writeback"`
	AnonPages         uint64 `json:"anon_pages"`
	Mapped            uint64 `json:"mapped"`
	Shmem             uint64 `json:"shmem"`
	KReclaimable      uint64 `json:"k_reclaimable"`
	Slab              uint64 `json:"slab"`
	SReclaimable      uint64 `json:"s_reclaimable"`
	SUnreclaim        uint64 `json:"s_unreclaim"`
	KernelStack       uint64 `json:"kernel_stack"`
	PageTables        uint64 `json:"pagetables"`
	SecPageTables     uint64 `json:"sec_page_tables"`
	NFSUnstable       uint64 `json:"nfs_unstable"`
	Bounce            uint64 `json:"bounce"`
	WritebackTmp      uint64 `json:"writeback_tmp"`
	CommitLimit       uint64 `json:"commit_limit"`
	CommittedAS       uint64 `json:"committed_as"`
	VmallocTotal      uint64 `json:"vmalloc_total"`
	VmallocUsed       uint64 `json:"vmalloc_used"`
	VmallocChunk      uint64 `json:"vmalloc_chunk"`
	Percpu            uint64 `json:"percpu"`
	HardwareCorrupted uint64 `json:"hardware_corrupted"`
	AnonHugePages     uint64 `json:"anon_huge_pages"`
	ShmemHugePages    uint64 `json:"shmem_huge_pages"`
	ShmemPmdMapped    uint64 `json:"shmem_pmd_mapped"`
	FileHugePages     uint64 `json:"file_huge_pages"`
	FilePmdMapped     uint64 `json:"file_pmd_mapped"`
	CmaTotal          uint64 `json:"cma_total"`
	CmaFree           uint64 `json:"cma_free"`
	HugePagesTotal    uint64 `json:"huge_pages_total"`
	HugePagesFree     uint64 `json:"huge_pages_free"`
	HugePagesRsvd     uint64 `json:"huge_pages_rsvd"`
	HugePagesSurp     uint64 `json:"huge_pages_surp"`
	Hugepagesize      uint64 `json:"hugepagesize"`
	Hugetlb           uint64 `json:"hugetlb"`
	DirectMap4K       uint64 `json:"direct_map_4k"`
	DirectMap2M       uint64 `json:"direct_map_2m"`
	DirectMap1G       uint64 `json:"direct_map_1g"`
}

func (m Meminfo) Used() uint64 {
	return m.MemTotal - m.MemFree - m.Buffers - m.Cached
}

func (m Meminfo) Percent() float64 {
	used := m.Used()
	return float64(used) / float64(m.MemTotal) * 100.0
}

func (m Meminfo) String() string {
	s, _ := json.MarshalIndent(m, "", "  ")
	return string(s)
}

func (m *Meminfo) setValue(key string, val uint64) {
	switch key {
	case "MemTotal":
		m.MemTotal = val
	case "MemFree":
		m.MemFree = val
	case "MemAvailable":
		m.MemAvailable = val
	case "Buffers":
		m.Buffers = val
	case "Cached":
		m.Cached = val
	case "SwapCached":
		m.SwapCached = val
	case "Active":
		m.Active = val
	case "Inactive":
		m.Inactive = val
	case "Active(anon)":
		m.ActiveAnon = val
	case "Inactive(anon)":
		m.InactiveAnon = val
	case "Active(file)":
		m.ActiveFile = val
	case "Inactive(file)":
		m.InactiveFile = val
	case "Unevictable":
		m.Unevictable = val
	case "Mlocked":
		m.Mlocked = val
	case "SwapTotal":
		m.Mlocked = val
	case "SwapFree":
		m.SwapTotal = val
	case "Zswap":
		m.SwapFree = val
	case "Zswapped":
		m.Zswap = val
	case "Dirty":
		m.Zswapped = val
	case "Writeback":
		m.Dirty = val
	case "AnonPages":
		m.Writeback = val
	case "Mapped":
		m.AnonPages = val
	case "Shmem":
		m.Mapped = val
	case "KReclaimable":
		m.KReclaimable = val
	case "Slab":
		m.Slab = val
	case "SReclaimable":
		m.SReclaimable = val
	case "SUnreclaim":
		m.SUnreclaim = val
	case "KernelStack":
		m.KernelStack = val
	case "PageTables":
		m.PageTables = val
	case "SecPageTables":
		m.SecPageTables = val
	case "NFS_Unstable":
		m.NFSUnstable = val
	case "Bounce":
		m.Bounce = val
	case "WritebackTmp":
		m.WritebackTmp = val
	case "CommitLimit":
		m.CommitLimit = val
	case "Committed_AS":
		m.CommittedAS = val
	case "VmallocTotal":
		m.VmallocTotal = val
	case "VmallocUsed":
		m.VmallocUsed = val
	case "VmallocChunk":
		m.VmallocChunk = val
	case "Percpu":
		m.Percpu = val
	case "HardwareCorrupted":
		m.HardwareCorrupted = val
	case "AnonHugePages":
		m.AnonHugePages = val
	case "ShmemHugePages":
		m.ShmemHugePages = val
	case "ShmemPmdMapped":
		m.ShmemPmdMapped = val
	case "FileHugePages":
		m.FileHugePages = val
	case "FilePmdMapped":
		m.FilePmdMapped = val
	case "CmaTotal":
		m.CmaTotal = val
	case "CmaFree":
		m.CmaFree = val
	case "HugePages_Total":
		m.HugePagesTotal = val
	case "HugePages_Free":
		m.HugePagesFree = val
	case "HugePages_Rsvd":
		m.HugePagesRsvd = val
	case "HugePagesSurp":
		m.HugePagesSurp = val
	case "Hugepagesize":
		m.Hugepagesize = val
	case "Hugetlb":
		m.Hugetlb = val
	case "DirectMap4k":
		m.DirectMap4K = val
	case "DirectMap2M":
		m.DirectMap2M = val
	case "DirectMap1G":
		m.DirectMap1G = val
	}
}
