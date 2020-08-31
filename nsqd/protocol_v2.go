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

	// messagePump 将需要发送的client的message从缓存池子里面pump出来
	messagePumpStartChan := make(chan bool)
	go p.messagePump(client, messagePumpStartChan)
	<-messagePumpStartChan

	for {
		// nsqd与客户端建立的tcp连接conn被包装在clientV2中，nsqd和客户端都会向conn中读写数据
		// nsqd会向客户端周期性发送heartbeat，并期待客户端响应，这个heartbeat是由客户端配置的
		// 如果nsqd两个周期内没有收到回应，认为客户端不在线，nsqd将关闭服务端连接
		if client.HeartbeatInterval > 0 {
			client.SetReadDeadline(time.Now().Add(client.HeartbeatInterval * 2))
		}
	}
}
