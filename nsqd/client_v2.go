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

	pubCounts map[string]uint64 // 如果是生产者客户端，则统计已生产消息数量

	writeLock sync.RWMutex
	metaLock  sync.RWMutex

	ID        int64 // 新增一个客户端，原子+1
	ctx       *context
	UserAgent string

	net.Conn

	flateWriter *flate.Writer

	Reader *bufio.Reader // tcp连接的reader
	Writer *bufio.Writer // tcp连接的writer

	OutputBufferSize    int           // 写缓冲区大小
	OutputBufferTimeout time.Duration // 写数据缓冲超时时间

	HeartbeatInterval time.Duration

	MsgTimeout time.Duration // 消息超时时间

	State          int32
	ConnectTime    time.Time // 首次连接时间
	Channel        *Channel
	ReadyStateChan chan int // client是否ready
	ExitChan       chan int // client断开信号

	ClientID string
	Hostname string

	SampleTate int32 // 采样频率  TODO: 含义未知

	IdentifyEventChan chan identifyEvent
	SubEventChan      chan *Channel

	Snappy  int32
	Deflate int32

	lenBuf   [4]byte
	lenSlice []byte
}

func newClientV2(id int64, conn net.Conn, ctx *context) *clinetV2 {
	var identifier string
	if conn != nil {
		identifier, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	}

	c := &clinetV2{
		ID:     id,
		ctx:    ctx,
		Conn:   conn,
		Reader: bufio.NewReaderSize(conn, defaultBufferSize),
		Writer: bufio.NewWriterSize(conn, defaultBufferSize),

		OutputBufferSize:    defaultBufferSize,
		OutputBufferTimeout: ctx.nsqd.getOpts().OutputBufferTimeout,

		MsgTimeout: ctx.nsqd.getOpts().MsgTimeout,

		ReadyStateChan: make(chan int, 1),
		ExitChan:       make(chan int),
		ConnectTime:    time.Now(),
		State:          stateInit,

		ClientID: identifier,
		Hostname: identifier,

		SubEventChan:      make(chan *Channel, 1),
		IdentifyEventChan: make(chan identifyEvent, 1),

		HeartbeatInterval: ctx.nsqd.getOpts().ClientTimeout / 2,

		pubCounts: make(map[string]uint64),
	}
}
