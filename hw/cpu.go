package hw

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
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

func CpuID() (id string) {
	find := func(s string, exp string) {
		if id == "" {
			re := regexp.MustCompile(`ID: (.*)`)
			sm := re.FindStringSubmatch(s)
			if len(sm) == 2 {
				id = sm[1]
			}
		}
	}
	find(fs.ReadString("/proc/cpuinfo"), `Serial\t*: (.*)`)
	find(fs.ReadStdout("dmidecode", "--type", "processor"), `ID: (.*)`)
	find(fs.ReadStdout("lshw"), `serial: (.*)`)
	id = strings.ReplaceAll(id, " ", "")
	id = strings.ToLower(id)
	return
}
