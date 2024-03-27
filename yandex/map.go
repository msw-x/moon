package yandex

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/msw-x/moon/ufmt"
)

type MapType string

const (
	MapShema     MapType = "map"
	MapSatellite MapType = "sat"
	MapHybrid    MapType = "skl"
)

const MapTraffic MapType = "trf"
const MapPanoramas MapType = "stv"

type Location struct {
	Latitude  float32
	Longitude float32
}

func (o Location) Empty() bool {
	return o.Latitude == 0 || o.Longitude == 0
}

func (o Location) String() string {
	return ufmt.JoinWith(",", o.Longitude, o.Latitude)
}

type Points []Location

func (o Points) Empty() bool {
	return len(o) == 0
}

func (o Points) String() string {
	return ufmt.JoinSliceWith("~", o)
}

type Map struct {
	Type      MapType
	Traffic   bool
	Panoramas bool
	Scale     int // 1-19
	Center    Location
	Point     Location
	Points    Points
	Text      string
}

func (o Map) Url() string {
	var u url.URL
	u.Scheme = "https"
	u.Host = "yandex.ru"
	u.Path = "maps"
	q := u.Query()
	if !o.Center.Empty() {
		q.Set("ll", o.Center.String())
	}
	if o.Text == "" {
		var points Points
		if !o.Point.Empty() {
			points = append(points, o.Point)
		}
		points = append(points, o.Points...)
		if !points.Empty() {
			q.Set("pt", points.String())
		}
	} else {
		q.Set("text", o.Text)
	}
	if o.Scale != 0 {
		q.Set("z", strconv.Itoa(o.Scale))
	}
	var l []MapType
	if o.Type != "" {
		l = append(l, o.Type)
	}
	if o.Traffic {
		l = append(l, MapTraffic)
	}
	if o.Panoramas {
		l = append(l, MapPanoramas)
	}
	if len(l) != 0 {
		q.Set("l", ufmt.JoinSliceWith(",", l))
	}
	u.RawQuery = q.Encode()
	return strings.ReplaceAll(u.String(), "%2C", ",")
}
