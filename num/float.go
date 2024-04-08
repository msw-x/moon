package num

type Float64d1 float64

type Float64d2 float64

type Float64d3 float64

type Float64d4 float64

type Float64d5 float64

type Float64d6 float64

func (o Float64d1) Value() float64 {
	return float64(o)
}

func (o Float64d2) Value() float64 {
	return float64(o)
}

func (o Float64d3) Value() float64 {
	return float64(o)
}

func (o Float64d4) Value() float64 {
	return float64(o)
}

func (o Float64d5) Value() float64 {
	return float64(o)
}

func (o Float64d6) Value() float64 {
	return float64(o)
}

func (o Float64d1) String() string {
	return DelicateFloat64(o.Value(), 1)
}

func (o Float64d2) String() string {
	return DelicateFloat64(o.Value(), 1)
}

func (o Float64d3) String() string {
	return DelicateFloat64(o.Value(), 1)
}

func (o Float64d4) String() string {
	return DelicateFloat64(o.Value(), 1)
}

func (o Float64d5) String() string {
	return DelicateFloat64(o.Value(), 1)
}

func (o Float64d6) String() string {
	return DelicateFloat64(o.Value(), 1)
}
