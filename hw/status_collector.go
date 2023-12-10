package hw

type StatusCollector struct {
	Temperature    *Series[int]
	CpuUsedPercent *Series[int]
	RamTotal       *Series[uint64]
	RamUsed        *Series[uint64]
	DiskTotal      *Series[uint64]
	DiskUsed       *Series[uint64]
}

func NewStatusCollector(capacity int) *StatusCollector {
	o := new(StatusCollector)
	o.Temperature = NewSeries[int](capacity)
	o.CpuUsedPercent = NewSeries[int](capacity)
	o.RamTotal = NewSeries[uint64](capacity)
	o.RamUsed = NewSeries[uint64](capacity)
	o.DiskTotal = NewSeries[uint64](capacity)
	o.DiskUsed = NewSeries[uint64](capacity)
	return o
}

func (o *StatusCollector) Poll() Status {
	s := GetStatus()
	o.Temperature.Add(s.Temperature)
	o.CpuUsedPercent.Add(s.CpuUsedPercent)
	o.RamTotal.Add(s.RamTotal)
	o.RamUsed.Add(s.RamUsed)
	o.DiskTotal.Add(s.DiskTotal)
	o.DiskUsed.Add(s.DiskUsed)
	return o.Get()
}

func (o *StatusCollector) Get() (s Status) {
	s.Temperature = o.Temperature.Average()
	s.CpuUsedPercent = o.CpuUsedPercent.Average()
	s.RamTotal = o.RamTotal.Average()
	s.RamUsed = o.RamUsed.Average()
	s.DiskTotal = o.DiskTotal.Average()
	s.DiskUsed = o.DiskUsed.Average()
	return s.Enrich()
}
