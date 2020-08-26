package nsqd

import (
	"net"
	"sync"
)

type tcpServer struct {
	ctx   *context
	conns sync.Map
}

func (p *tcpServer) Handle(clientConn net.Conn) {
	p.ctx.nsqd.logf(INFO, "TCP: new clinet(%s)", clientConn.RemoteAddr())

}
