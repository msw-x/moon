package hw

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"github.com/msw-x/moon/fs"
)

func ThermalZoneFile(zone int) string {
	return fmt.Sprintf("/sys/devices/virtual/thermal/thermal_zone%d/temp", zone)
}

func ThermalZone(zone int) (temp int, err error) {
	raw := []byte{}
	raw, err = ioutil.ReadFile(ThermalZoneFile(zone))
	if err != nil {
		return
	}
	s := string(raw)
	s = strings.TrimSuffix(s, "\n")
	var i int
	i, err = strconv.Atoi(s)
	if err != nil {
		return
	}
	if math.Abs(float64(i)) > 1000 {
		i /= 1000
	}
	temp = i
	return
}

func CpuTemp() int {
	for i := 8; i != 0; i-- {
		if fs.Exist(ThermalZoneFile(i)) {
			temp, err := ThermalZone(i)
			if err == nil && temp > -50 && temp < 200 {
				return temp
			}
		}
	}
	return -273
}
