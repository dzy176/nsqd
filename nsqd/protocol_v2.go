package nsqd

import (
	"net"
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

// 根据TCP客户端不同的请求做出相应的处理
func (p *protocolV2) IOLoop(conn net.Conn) error {

}
