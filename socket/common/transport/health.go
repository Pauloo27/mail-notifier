package transport

import "time"

const (
	maxTimeWithoutReceiving = 40 * time.Second
	heartbeatRate           = 30 * time.Second
)

type Health struct {
	lastTimeSent     time.Time
	lastTimeReceived time.Time
	timer            *time.Timer
	unhealthCallback func()
	dead             bool
}

func newHealth(unhealthCallback func()) *Health {
	now := time.Now()
	timer := time.NewTimer(maxTimeWithoutReceiving)
	go func() {
		<-timer.C
		unhealthCallback()
	}()
	return &Health{
		lastTimeSent:     now,
		lastTimeReceived: now,
		timer:            timer,
		unhealthCallback: unhealthCallback,
		dead:             false,
	}
}

func (h *Health) HeartbeatSent() {
	h.lastTimeSent = time.Now()
}

func (h *Health) HeartbeatReceived() {
	h.lastTimeReceived = time.Now()
	h.timer.Reset(maxTimeWithoutReceiving)
}
