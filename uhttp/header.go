package uhttp

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"

	"github.com/msw-x/moon/refl"
	"github.com/msw-x/moon/ufmt"
)

func HeaderString(h http.Header) string {
	if len(h) == 0 {
		return ""
	}
	var l []string
	for name, vals := range h {
		l = append(l, fmt.Sprintf("%s: %s", name, ufmt.JoinSliceWith("; ", vals)))
	}
	sort.Slice(l, func(i, j int) bool {
		return l[i] < l[j]
	})
	return ufmt.JoinSliceWith("\n", l)
}

func HeaderTo(h http.Header, s any) (err error) {
	refl.WalkOnTags(s, HeaderTag, func(v reflect.Value, name string, flags []string) {
		if v.CanSet() && err == nil {
			s := h.Get(name)
			if s != "" {
				if v.Kind() == reflect.Ptr {
					v.Set(reflect.New(v.Type().Elem()))
					v = v.Elem()
				}
				err = refl.SetFromString(v, s)
			}
		}
	})
	return
}
