package websocket

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
