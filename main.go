package main

import (
	"nsqd/nsqd"
	"sync"
)

type program struct {
	once sync.Once
	nsqd *nsqd.NSQD
}

func main() {

}
