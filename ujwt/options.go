package ujwt

import "time"

type Options struct {
	Key                   string
	ExpirationTime        time.Duration
	RefreshExpirationTime time.Duration
}

func (o Options) KeyBytes() []byte {
	return []byte(o.Key)
}
