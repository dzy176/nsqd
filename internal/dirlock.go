package internal

import (
	"fmt"
	"os"
	"syscall"
)

type DirLock struct {
	dir string
	f   *os.File
}

func New(dir string) *DirLock {
	return &DirLock{
		dir: dir,
	}
}

func (l *DirLock) Lock() error {
	f, err := os.Open(l.dir)
	if err != nil {
		return err
	}

	l.f = f

	// LOCK_EX: 排他锁，不允许其他人读写
	// LOCK_NB: 无法锁定文件时不阻塞，立马返回给进程
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return fmt.Errorf("floclk failed, err: %v", err)
	}
	return nil
}

func (l *DirLock) Unlock() error {
	defer l.f.Close()
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}
