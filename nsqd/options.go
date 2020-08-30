package nsqd

import (
	"crypto/md5"
	"hash/crc32"
	"io"
	"log"
	"os"
	"time"
)

type Options struct {
	ID        int64    `flag:"node-id" cfg:"id"`
	LogLevel  LogLevel `flag:"log-level"`
	LogPrefix string   `flag:"log-prefix"`
	Logger    Logger

	TCPAddress               string        `flag:"tcp-address"`
	HTTPAddress              string        `flag:"http-address"`
	BroadcastAddress         string        `flag:"broadcast-address"`
	NSQLookupdTCPAddress     []string      `flag:"lookupd-tcp-address" cfg:"nsqlookup_tcp_address"`
	AuthHTTPAddress          []string      `flag:"auth-http-address" cfg:"auth_http_address"`
	HTTPClientConnectTimeout time.Duration `flag:"http-client-connect-timeout" cfg:"http_client_connect_timeout"`
	HTTPClientRequestTimeout time.Duration `flag:"http-client-request-timeout" cfg:"http_client_request_timeout"`

	// disk queue config
	DataPath        string        `flag:"data-path"`
	MemQueueSize    int64         `flag:"mem-queue-size"`
	MaxBytesPerFile int64         `flag:"max-bytes-per-file"`
	SyncEvery       int64         `flag:"sync-every"`
	SyncTimeout     time.Duration `flag:"sync-timeout"`

	// 扫描延时队列或者未ack队列的时间间隔
	QueueScanInterval time.Duration
	// 调整工作线程和刷新topic.channel的时间间隔
	QueueScanRefreshInterval time.Duration
	// 扫描channel的数量
	QueueScanSelectionCount int `flag:"queue-scan-selection-count"`
	// 扫描各个队列的最大工作线程
	QueueScanWorkPoolMax int `flag:"queue-scan-work-pool-max"`
	// 允许扫描失败的比例
	QueueScanDirtyPercent float64

	// msg and command options
	// msg超过这个时长未被ack则重新入队，重新消费
	MsgTimeout time.Duration `flag:"msg-timeout"`
	//
	MaxMsgTimeOut time.Duration `flag:"max-msg-timeout"`
	MaxMsgSize    int64         `flag:"max-msg-size"`
	MaxBodySize   int64         `flag:"max-body-size"`
	MaxReqTimeout time.Duration `flag:"max-req-timeout"`
	ClientTimeout time.Duration

	// client overridable configuration options
	MaxHeartbeatInterval   time.Duration `flag:"max-heartbeat-interval"`
	MaxRdyCount            int64         `flag:"max-rdy-count"`
	MaxOutputBufferSize    int64         `flag:"max-output-buffer-size"`
	MaxOutputBufferTimeout time.Duration `flag:"max-output-buffer-timeout"`
	MinOutputBufferTimeout time.Duration `flag:"min-output-buffer-timeout"`
	OutputBufferTimeout    time.Duration `flag:"output-buffer-timeout"`
	MaxChannelConsumers    int           `flag:"max-channel-consumers"`
}

func NewOptions() *Options {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	h := md5.New()
	io.WriteString(h, hostname)
	defaltID := int64(crc32.ChecksumIEEE(h.Sum(nil)) % 1024)

	return &Options{
		ID:        defaltID,
		LogPrefix: "[nsqd] ",
		LogLevel:  INFO,

		TCPAddress:       "[::]:4150",
		HTTPAddress:      "[::]:4151",
		BroadcastAddress: hostname,

		NSQLookupdTCPAddress: make([]string, 0),
		AuthHTTPAddress:      make([]string, 0),

		HTTPClientConnectTimeout: 2 * time.Second,
		HTTPClientRequestTimeout: 5 * time.Second,

		MemQueueSize:    10000,
		MaxBytesPerFile: 100 * 1024 * 1024,
		SyncEvery:       2500,

		QueueScanInterval:        100 * time.Millisecond,
		QueueScanRefreshInterval: 5 * time.Second,
		QueueScanSelectionCount:  20,
		QueueScanWorkPoolMax:     4,
		QueueScanDirtyPercent:    0.25,

		MsgTimeout:    time.Minute,
		MaxMsgTimeOut: 15 * time.Minute,
		MaxMsgSize:    1024 * 1024,
		MaxBodySize:   5 * 1024 * 1024,
		MaxReqTimeout: 1 * time.Hour,
		ClientTimeout: 60 * time.Second,

		MaxHeartbeatInterval:   60 * time.Second,
		MaxRdyCount:            2500,
		MaxOutputBufferSize:    64 * 1024,
		MaxOutputBufferTimeout: 30 * time.Second,
		MinOutputBufferTimeout: 25 * time.Millisecond,
		OutputBufferTimeout:    250 * time.Millisecond,
		MaxChannelConsumers:    0,
	}
}
