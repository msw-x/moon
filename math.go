package moon

import "math/rand"

func Percent[Minor AnyInt, Major AnyInt](minor Minor, major Major) int {
	mn := int64(minor)
	mj := int64(major)
	if mn == 0 {
		return 0
	}
	return int(float64(mn) / float64(mj) * 100)
}

func Rand(min, max int) int {
	return rand.Intn(max-min) + min
}
