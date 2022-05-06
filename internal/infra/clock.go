package infra

import (
	"github.com/yuridevx/app-example/internal/config"
	"time"
)

type Clock struct {
	loc *time.Location
}

func (c *Clock) Now() time.Time {
	return time.Now().In(c.loc)
}

func NewClock(conf *config.Config) *Clock {
	loc, err := time.LoadLocation(conf.Timezone)
	if err != nil {
		panic(err)
	}
	return &Clock{loc: loc}
}
