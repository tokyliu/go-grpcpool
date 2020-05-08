package pool

import (
	"google.golang.org/grpc"
	"sync"
)

// net.Conn's Close() method.
type PoolConn struct {
	Conn     *grpc.ClientConn
	mu       sync.RWMutex
	c        *channelPool
	unusable bool
}

// Close() puts the given connects back to the pool instead of closing it.
func (p *PoolConn) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.Conn != nil {
			return p.Conn.Close()
		}
		return nil
	}
	return p.c.put(p.Conn)
}

// MarkUnusable() marks the connection not usable any more, to let the pool close it instead of returning it to pool.
func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

func (c *channelPool) wrapConn(conn *grpc.ClientConn) *PoolConn {
	p := &PoolConn{c: c}
	p.Conn = conn
	return p
}
