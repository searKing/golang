package tcp

import "time"

func (srv *Server) initialReadLimitSize() int64 {
	return int64(srv.maxBytes()) + 4096 // bufio slop
}

const DefaultMaxBytes = 1 << 20 // 1 MB
func (srv *Server) maxBytes() int {
	if srv.MaxBytes > 0 {
		return srv.MaxBytes
	}
	return DefaultMaxBytes
}
func (s *Server) idleTimeout() time.Duration {
	if s.IdleTimeout != 0 {
		return s.IdleTimeout
	}
	return s.ReadTimeout
}
func (s *Server) readTimeout() time.Duration {
	if s.ReadTimeout != 0 {
		return s.ReadTimeout
	}
	return s.ReadTimeout
}

func (s *Server) shuttingDown() bool {
	return s.inShutdown.Load()
}
func (s *Server) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getDoneChanLocked()
}

func (s *Server) getDoneChanLocked() chan struct{} {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

func (s *Server) closeDoneChanLocked() {
	ch := s.getDoneChanLocked()
	select {
	case <-ch:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded
		// by s.mu.
		close(ch)
	}
}
