package transport

import "time"

const (
	maxTimeWithoutReceiving = 10 * time.Second
	heartbeatRate           = 5 * time.Second

	heartbeatCommandName = "heartbeat"
)

var (
	timer *time.Timer
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
	timer = time.NewTimer(maxTimeWithoutReceiving)
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

func (h *Health) Kill() {
	h.dead = true
	timer.Stop()
}

func (h *Health) HeartbeatSent() {
	h.lastTimeSent = time.Now()
}

func (h *Health) HeartbeatReceived() {
	h.lastTimeReceived = time.Now()
	h.timer.Reset(maxTimeWithoutReceiving)
}
