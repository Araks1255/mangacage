package rotate_logger

import (
	"time"
)

func (rl *RotateLogger) Start() {
	for {
		nextMidnight := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)
		timer := time.NewTimer(time.Until(nextMidnight))
		<-timer.C

		rl.rotate()
	}
}
