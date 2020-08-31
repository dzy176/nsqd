package nsqd

import "sync/atomic"

type ClientStats struct {
	ClientID      string     `json:"client_id"`
	Hostname      string     `json:"hostname"`
	Version       string     `json:"version"`
	RemoteAddress string     `json:"remote_address"`
	State         int32      `json:"state"`
	ReadyCount    int64      `json:"ready_count"`
	InFlightCount int64      `json:"in_flight_count"`
	MessageCount  uint64     `json:"message_count"`
	FinishedCount uint64     `json:"finished_count"`
	RequeueCount  uint64     `json:"requeue_count"`
	ConnectTime   int64      `json:"connect_time"`
	SampleRate    int32      `json:"sample_rate"`
	Deflate       bool       `json:"deflate"`
	Snappy        bool       `json:"snappy"`
	UserAgent     string     `json:"user_agent"`
	PubCounts     []PubCount `json:"pub_counts,omitempty"`
}

type PubCount struct {
	Topic string `json:"topic"`
	Count uint64 `json:"count"`
}

type TopicStats struct {
	TopicName            string           `json:"topic_name"`
	Channels             []ChannelStats   `json:"channels"`
	Depth                int64            `json:"depth"`
	BackendDepth         int64            `json:"backend_depth"`
	MessageCount         uint64           `json:"message_count"`
	MessageBytes         uint64           `json:"message_bytes"`
	Paused               bool             `json:"paused"`
	E2eProcessingLatency *quantile.Result `json:"e2e_processing_latency"`
}

func NewTopicStats(t *Topic, channels []ChannelStats) TopicStats {
	return TopicStats{
		TopicName:            t.name,
		Channels:             channels,
		Depth:                t.Depth(),
		BackendDepth:         t.backend.Depth(),
		MessageCount:         atomic.LoadUint64(&t.messageCount),
		MessageBytes:         atomic.LoadUint64(&t.messageBytes),
		Paused:               t.IsPaused(),
		E2eProcessingLatency: t.AggreateChannelE2eProcessingLatency().Result(),
	}
}
