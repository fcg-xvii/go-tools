package ami

import "time"

type EventListener struct {
	uuid       int64
	eventChan  chan Event
	timeActual time.Time
}

func (s *EventListener) incomingEvent(e Event) {
	s.timeActual = time.Now().Add(time.Minute)
	s.eventChan <- e
}

func (s *EventListener) close() {
	close(s.eventChan)
}
