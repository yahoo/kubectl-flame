package api

type ProfilingEvent string

const (
	Cpu         ProfilingEvent = "cpu"
	Alloc       ProfilingEvent = "alloc"
	Lock        ProfilingEvent = "lock"
	CacheMisses ProfilingEvent = "cache-misses"
	Wall        ProfilingEvent = "wall"
	Itimer      ProfilingEvent = "itimer"
)

var (
	supportedEvents = []ProfilingEvent{Cpu, Alloc, Lock, CacheMisses, Wall, Itimer}
)

func AvailableEvents() []ProfilingEvent {
	return supportedEvents
}

func IsSupportedEvent(event string) bool {
	return containsEvent(ProfilingEvent(event), AvailableEvents())
}

func containsEvent(e ProfilingEvent, events []ProfilingEvent) bool {
	for _, current := range events {
		if e == current {
			return true
		}
	}
	return false
}
