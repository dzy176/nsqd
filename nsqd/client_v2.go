package nsqd

import (
	"bufio"
	"compress/flate"
	"net"
	"sync"
	"time"
)

const (
	defaultBufferSize = 16 * 1024

	stateInit = iota
	stateDisconnected
	stateConnected
	stateClosing
)

type identifyDataV2 struct {
	ClientID           string `json:"client_id"`
	Hostname           string `json:"hostname"`
	HeartbeatInterval  int    `json:"heartbeat_interval"`
	OutputBufferSize   int    `json:"output_buffer_size"`
	FeatureNegotiation bool   `json:"feature_negotiation"`
	Deflate            bool   `json:"deflate"`       // 是否deflate压缩
	DeflateLever       int    `json:"deflate_lever"` // 压缩级别，级别越高，CPU负载越高
	Snappy             bool   `json:"snappy"`        // 是否snappy压缩
	SampleRate         int32  `json:"sample_rate"`
	UserAgent          string `json:"user_agent"`
	MsgTimeout         int    `json:"msg_timeout"`
}

type identifyEvent struct {
	OutputBufferTimeout time.Duration
	HeartbeatInterval   time.Duration
	SampleRate          int32
	MsgTimeout          time.Duration
}

type clinetV2 struct {
	ReadyCount    int64
	InFlightCount int64
	MessageCount  uint64
	FinishCount   uint64
	RequeueCount  uint64

	pubCounts map[string]uint64

	writeLock sync.RWMutex
	metaLock  sync.RWMutex

	ID        int64
	ctx       *context
	UserAgent string

	net.Conn

	flateWriter *flate.Writer

	Reader *bufio.Reader
	Writer *bufio.Writer

	OutputBufferSize    int
	OutputBufferTimeout time.Duration

	HeartbeatInterval time.Duration

	MsgTimeout time.Duration

	State          int32
	ConnectTime    time.Time
	Channel        *Channel
	ReadyStateChan chan int
	ExitChan       chan int

	ClientID string
	Hostname string

	// 采样率
	SampleTate int32

	IdentifyEventChan chan identifyEvent
	SubEventChan      chan *Channel

	Snappy  int32
	Deflate int32

	lenBuf   [4]byte
	lenSlice []byte
}
