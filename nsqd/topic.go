package nsqd

import (
	"net/mail"
	"nsqd/internal"
	"sync"
)

type Topic struct {
	messageCount uint64
	messageBytes uint64

	sync.RWMutex

	name              string
	channelMap        map[string]*Channel
	backend           BackendQueue
	memoryMsgChan     chan *Message
	startChan         chan int
	exitChan          chan int
	channelUpdateChan chan int
	waitGroup         internal.WaitGroupWrapper
	exitFlag          int32
	idFactory         *guidFactory

	deleteCallback func(topic *Topic)
	deleter        sync.Once

	paused    int32
	pauseChan chan int

	ctx *context
}
