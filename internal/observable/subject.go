package observable

import "sync"

type Subject struct {
	observers []Observer
	mu        sync.RWMutex
}

func (s *Subject) Attach(o Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, o)
}

func (s *Subject) Detach(o Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, observer := range s.observers {
		if observer == o {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *Subject) Notify(message string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, observer := range s.observers {
		observer.Update(message)
	}
}