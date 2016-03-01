package seqdb

import (
	"time"
	"fmt"
)

const TIMEOUT=10

type BucketLocks struct {
	locks map[string]bool
}

func NewBucketLock() *BucketLocks {
	var l BucketLocks
	l.locks = make(map[string]bool)
	return &l
}

/**
	Lock a bucket
 */
func (b *BucketLocks) AddLock(s string) {
	b.locks[s] = true
}

/**
	Unlock a bucket
 */
func (b *BucketLocks) RemoveLock(s string) {
	b.locks[s] = false
}

/**
	Check if the bucket is locked
 */
func (b *BucketLocks) IsLocked(s string) bool {
	if lockState, ok := b.locks[s]; ok {
		return lockState
	}

	return false
}

/**
	Wait for the bucket to be unlocked
 */
func (b *BucketLocks) WaitForLock(s string) error {

	start := time.Now().Second()

	for b.IsLocked(s) {
		// while the bucket is locked
		if cur := time.Now().Second() - start; cur > TIMEOUT {
			// we have waited for too long
			// break
			return fmt.Errorf("Lock timed out for %s bucket", s)
		}
		// sleep tenth of a second and then check again
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

/**
	Waits while a block is locked and immediately sets the lock again
 */
func (b *BucketLocks) WaitAndSet(s string) {
	b.WaitForLock(s)
	b.AddLock(s)
}