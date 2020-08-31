package nsqd

import (
	"net"
	"nsqd/internal"
	"nsqd/internal/dirlock"
	"sync"
	"sync/atomic"
	"time"
)

type Client interface {
	Stats() ClientStats
	IsProducer() bool
}

type NSQD struct {
	clinetIDSquence int64

	sync.RWMutex

	opts atomic.Value

	dl        *dirlock.DirLock
	isLoading int32

	//　atomic.value: 将任意类型的数据的读写操作封装成原子操作（中间态对外不可见）
	// 例如普通写入一个int64，底层需要进行两次写操作（低32位，高32位）
	// 如果一个线程刚写完低32，另一个线程读取这个变量，则读到了一个中间量
	errValue  atomic.Value
	startTime time.Time

	topicMap map[string]*Topic

	clientLock sync.RWMutex
	clients    map[int64]Client

	lookupPeers atomic.Value

	tcpServer   *tcpServer
	tcpListener net.Listener
	httpListern net.Listener

	poolSize int

	notifyChan           chan interface{}
	optsNotificationChan chan struct{}
	exitChan             int
	waitGroup            internal.WaitGroupWrapper
}

func (n *NSQD) logf(level LogLevel, f string, args ...interface{}) {
	opts := n.getOpts()
	Logf(opts.Logger, opts.LogLevel, level, f, args...)
}

func (n *NSQD) getOpts() *Options {
	return n.opts.Load().(*Options)
}

func (n *NSQD) AddClient(clientId int64, client Client) {
	n.clientLock.Lock()
	n.clients[clientId] = client
	n.clientLock.Unlock()
}
