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
func (srv *Server) idleTimeout() time.Duration {
	if srv.IdleTimeout != 0 {
		return srv.IdleTimeout
	}
	return srv.ReadTimeout
}
func (srv *Server) readTimeout() time.Duration {
	if srv.ReadTimeout != 0 {
		return srv.ReadTimeout
	}
	return srv.ReadTimeout
}

func (srv *Server) shuttingDown() bool {
	return srv.inShutdown.Load()
}
