package nsqd

import "sync"

type guid int64

type guidFactory struct {
	sync.Mutex

	nodeID        int64
	sequence      int64
	lastTimestamp int64
	lastID        guid
}
