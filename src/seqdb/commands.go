package seqdb
import (
	"github.com/boltdb/bolt"
	"fmt"
	"bytes"
	"encoding/binary"
)

type SeqPointer struct {
	BucketName string
	SequenceName string
	lock *BucketLocks
}

func (s SeqPointer) Set(v uint64) {

	Db.Update(func(tx *bolt.Tx) error {

		bucket, err := tx.CreateBucketIfNotExists([]byte(s.BucketName))

		if err != nil { return err }

		buf := new(bytes.Buffer)
		err = binary.Write(buf, binary.LittleEndian, v)

		if err != nil { return err }

		fmt.Println("SET", s.BucketName, s.SequenceName, v)

		err = bucket.Put([]byte(s.SequenceName), buf.Bytes())

		if err != nil { return err }

		return nil
	})

}

func (s SeqPointer) Get() uint64 {

	var val uint64;

	Db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(s.BucketName))

		if bucket == nil {
			return fmt.Errorf("Bucket %s not found", s.BucketName)
		}

		b := bucket.Get([]byte(s.SequenceName))

		buf := bytes.NewReader(b)
		err := binary.Read(buf, binary.LittleEndian, &val)

		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}

		return nil
	})

	return val
}

func (s SeqPointer) Inc() uint64 {
	var val uint64

	Db.Batch(func(tx *bolt.Tx) error {
		// Get the given sequence from bucket
		bucket := tx.Bucket([]byte(s.BucketName))

		if bucket == nil { return fmt.Errorf("Bucket %s not found", s.BucketName) } // Bucket is not found

		// Convert the binary to Uint64
		bufR := bytes.NewReader(bucket.Get([]byte(s.SequenceName)))
		err := binary.Read(bufR, binary.LittleEndian, &val)

		if err != nil { return err }

		// Increment variable
		val++

		// Put it back
		bufW := new(bytes.Buffer)
		err = binary.Write(bufW, binary.LittleEndian, val)

		if err != nil { return err }

		err = bucket.Put([]byte(s.SequenceName), bufW.Bytes())

		return err
	})

	return val
}

func (s SeqPointer) Lock() {
	s.lock.WaitForLock(s.BucketName)
}

func (s SeqPointer) Unlock() {
	s.lock.RemoveLock(s.BucketName)
}