package semaphore

type Semaphore struct {
	limit chan struct{}
}

func NewSemaphore(size uint64) *Semaphore {
	return &Semaphore{limit: make(chan struct{}, size)}
}

func (s *Semaphore) Acquire() {
	s.limit <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.limit
}
