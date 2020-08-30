package nsqd

import (
	"io"
	"net"
	"nsqd/internal/protocol"
	"sync"
)

type tcpServer struct {
	ctx   *context
	conns sync.Map
}

func (p *tcpServer) Handle(clientConn net.Conn) {
	p.ctx.nsqd.logf(INFO, "TCP: new clinet(%s)", clientConn.RemoteAddr())
	// 客户端建立连接后先发送4字节的协议号
	buf := make([]byte, 4)
	_, err := io.ReadFull(clientConn, buf)
	if err != nil {
		p.ctx.nsqd.logf(ERROR, "fail to read protocol version - %s", err)
		clientConn.Close()
		return
	}
	protocolMagic := string(buf)

	p.ctx.nsqd.logf(INFO, "CLIENT(%s): desired protocol magic: %s",
		clientConn.RemoteAddr(), protocolMagic)

	var prot protocol.Protocol
	switch protocolMagic {
	case "  V2":
		prot = &pro

	}
}
