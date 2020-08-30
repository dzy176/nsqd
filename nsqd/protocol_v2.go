package nsqd

import (
	"net"
	"sync/atomic"
	"time"
)

const (
	maxTimeout = time.Hour

	frameTypeResponse int32 = 0
	frameTypeError    int32 = 1
	frameTypeMessage  int32 = 2
)

var separatorBytes = []byte(" ")
var heartbeatBytes = []byte("_heartbeat_")
var okBytes = []byte("OK")

type protocolV2 struct {
	ctx *context
}

// 客户端分为生产者与消费者
// IOLoop根据TCP客户端不同的请求做出相应的处理
func (p *protocolV2) IOLoop(conn net.Conn) error {
	var err error
	var line []byte
	var zeroTime time.Time

	clientID := atomic.AddInt64(&p.ctx.nsqd.clinetIDSquence, 1)
	client := newClientV2(clientID, conn, p.ctx)
	p.ctx.nsqd.AddClient(clientID, client)
}
